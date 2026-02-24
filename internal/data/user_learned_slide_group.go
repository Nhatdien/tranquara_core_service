package data

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
)

type UserLearnedSlideGroup struct {
	ID           uuid.UUID `json:"id"`
	UserID       uuid.UUID `json:"user_id"`
	CollectionID uuid.UUID `json:"collection_id"`
	SlideGroupID string    `json:"slide_group_id"`
	CompletedAt  time.Time `json:"completed_at"`
}

type UserLearnedSlideGroupModel struct {
	DB *sql.DB
}

// Insert marks a slide group as completed for a user
// Uses ON CONFLICT to handle duplicate completions gracefully
func (m UserLearnedSlideGroupModel) Insert(learned *UserLearnedSlideGroup) (*UserLearnedSlideGroup, error) {
	query := `
		INSERT INTO user_learned_slide_groups (user_id, collection_id, slide_group_id)
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id, collection_id, slide_group_id) DO NOTHING
		RETURNING id, user_id, collection_id, slide_group_id, completed_at
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query,
		learned.UserID,
		learned.CollectionID,
		learned.SlideGroupID,
	).Scan(
		&learned.ID,
		&learned.UserID,
		&learned.CollectionID,
		&learned.SlideGroupID,
		&learned.CompletedAt,
	)

	if err != nil {
		// ON CONFLICT DO NOTHING returns no rows — treat as already exists
		if errors.Is(err, sql.ErrNoRows) {
			// Fetch the existing record
			return m.GetOne(learned.UserID, learned.CollectionID, learned.SlideGroupID)
		}
		return nil, err
	}

	return learned, nil
}

// GetOne retrieves a specific completion record
func (m UserLearnedSlideGroupModel) GetOne(userID, collectionID uuid.UUID, slideGroupID string) (*UserLearnedSlideGroup, error) {
	query := `
		SELECT id, user_id, collection_id, slide_group_id, completed_at
		FROM user_learned_slide_groups
		WHERE user_id = $1 AND collection_id = $2 AND slide_group_id = $3
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var learned UserLearnedSlideGroup
	err := m.DB.QueryRowContext(ctx, query, userID, collectionID, slideGroupID).Scan(
		&learned.ID,
		&learned.UserID,
		&learned.CollectionID,
		&learned.SlideGroupID,
		&learned.CompletedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}

	return &learned, nil
}

// GetByCollection retrieves all completed slide groups for a user in a specific collection
func (m UserLearnedSlideGroupModel) GetByCollection(userID, collectionID uuid.UUID) ([]*UserLearnedSlideGroup, error) {
	query := `
		SELECT id, user_id, collection_id, slide_group_id, completed_at
		FROM user_learned_slide_groups
		WHERE user_id = $1 AND collection_id = $2
		ORDER BY completed_at ASC
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, userID, collectionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []*UserLearnedSlideGroup
	for rows.Next() {
		var learned UserLearnedSlideGroup
		err = rows.Scan(
			&learned.ID,
			&learned.UserID,
			&learned.CollectionID,
			&learned.SlideGroupID,
			&learned.CompletedAt,
		)
		if err != nil {
			return nil, err
		}
		results = append(results, &learned)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

// GetAllByUser retrieves all completed slide groups for a user across all collections
func (m UserLearnedSlideGroupModel) GetAllByUser(userID uuid.UUID) ([]*UserLearnedSlideGroup, error) {
	query := `
		SELECT id, user_id, collection_id, slide_group_id, completed_at
		FROM user_learned_slide_groups
		WHERE user_id = $1
		ORDER BY completed_at DESC
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []*UserLearnedSlideGroup
	for rows.Next() {
		var learned UserLearnedSlideGroup
		err = rows.Scan(
			&learned.ID,
			&learned.UserID,
			&learned.CollectionID,
			&learned.SlideGroupID,
			&learned.CompletedAt,
		)
		if err != nil {
			return nil, err
		}
		results = append(results, &learned)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

// Delete removes a specific completion record
func (m UserLearnedSlideGroupModel) Delete(id, userID uuid.UUID) error {
	query := `
		DELETE FROM user_learned_slide_groups
		WHERE id = $1 AND user_id = $2
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := m.DB.ExecContext(ctx, query, id, userID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}

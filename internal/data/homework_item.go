package data

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
)

type HomeworkItem struct {
	ID          uuid.UUID  `json:"id"`
	SessionID   uuid.UUID  `json:"session_id"`
	UserID      string     `json:"user_id"`
	Content     string     `json:"content"`
	Completed   bool       `json:"completed"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
}

type HomeworkItemModel struct {
	DB *sql.DB
}

// Insert creates a new homework item
func (m HomeworkItemModel) Insert(item *HomeworkItem) (*HomeworkItem, error) {
	query := `
		INSERT INTO homework_items (session_id, user_id, content)
		VALUES ($1, $2, $3)
		RETURNING id, session_id, user_id, content, completed, completed_at, created_at
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query,
		item.SessionID,
		item.UserID,
		item.Content,
	).Scan(
		&item.ID,
		&item.SessionID,
		&item.UserID,
		&item.Content,
		&item.Completed,
		&item.CompletedAt,
		&item.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return item, nil
}

// GetBySession retrieves all homework items for a session
func (m HomeworkItemModel) GetBySession(sessionID uuid.UUID, userID string) ([]*HomeworkItem, error) {
	query := `
		SELECT id, session_id, user_id, content, completed, completed_at, created_at
		FROM homework_items
		WHERE session_id = $1 AND user_id = $2
		ORDER BY created_at ASC
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, sessionID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*HomeworkItem
	for rows.Next() {
		var item HomeworkItem
		err = rows.Scan(
			&item.ID,
			&item.SessionID,
			&item.UserID,
			&item.Content,
			&item.Completed,
			&item.CompletedAt,
			&item.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		items = append(items, &item)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

// GetAllByUser retrieves all homework items for a user across all sessions
func (m HomeworkItemModel) GetAllByUser(userID string) ([]*HomeworkItem, error) {
	query := `
		SELECT id, session_id, user_id, content, completed, completed_at, created_at
		FROM homework_items
		WHERE user_id = $1
		ORDER BY created_at ASC
	`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*HomeworkItem
	for rows.Next() {
		var item HomeworkItem
		err = rows.Scan(
			&item.ID,
			&item.SessionID,
			&item.UserID,
			&item.Content,
			&item.Completed,
			&item.CompletedAt,
			&item.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		items = append(items, &item)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

// Toggle flips the completed state of a homework item
func (m HomeworkItemModel) Toggle(id uuid.UUID, userID string, completed bool) (*HomeworkItem, error) {
	var completedAt *time.Time
	if completed {
		now := time.Now()
		completedAt = &now
	}

	query := `
		UPDATE homework_items
		SET completed = $1, completed_at = $2
		WHERE id = $3 AND user_id = $4
		RETURNING id, session_id, user_id, content, completed, completed_at, created_at
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var item HomeworkItem
	err := m.DB.QueryRowContext(ctx, query, completed, completedAt, id, userID).Scan(
		&item.ID,
		&item.SessionID,
		&item.UserID,
		&item.Content,
		&item.Completed,
		&item.CompletedAt,
		&item.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}

	return &item, nil
}

// Delete removes a homework item
func (m HomeworkItemModel) Delete(id uuid.UUID, userID string) error {
	query := `
		DELETE FROM homework_items
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

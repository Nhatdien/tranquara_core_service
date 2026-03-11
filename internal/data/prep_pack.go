package data

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
)

type PrepPack struct {
	ID             uuid.UUID       `json:"id"`
	UserID         string          `json:"user_id"`
	DateRangeStart time.Time       `json:"date_range_start"`
	DateRangeEnd   time.Time       `json:"date_range_end"`
	Content        json.RawMessage `json:"content"`
	JournalCount   int             `json:"journal_count"`
	PersonalNotes  *string         `json:"personal_notes,omitempty"`
	CreatedAt      time.Time       `json:"created_at"`
}

type PrepPackModel struct {
	DB *sql.DB
}

// Insert creates a new prep pack
func (m PrepPackModel) Insert(pack *PrepPack) (*PrepPack, error) {
	query := `
		INSERT INTO prep_packs (user_id, date_range_start, date_range_end, content, journal_count, personal_notes)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, user_id, date_range_start, date_range_end, content, journal_count, personal_notes, created_at
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var contentBytes []byte

	err := m.DB.QueryRowContext(ctx, query,
		pack.UserID,
		pack.DateRangeStart,
		pack.DateRangeEnd,
		pack.Content,
		pack.JournalCount,
		pack.PersonalNotes,
	).Scan(
		&pack.ID,
		&pack.UserID,
		&pack.DateRangeStart,
		&pack.DateRangeEnd,
		&contentBytes,
		&pack.JournalCount,
		&pack.PersonalNotes,
		&pack.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	pack.Content = json.RawMessage(contentBytes)
	return pack, nil
}

// Get retrieves a single prep pack by ID and user
func (m PrepPackModel) Get(id uuid.UUID, userID string) (*PrepPack, error) {
	query := `
		SELECT id, user_id, date_range_start, date_range_end, content, journal_count, personal_notes, created_at
		FROM prep_packs
		WHERE id = $1 AND user_id = $2
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var pack PrepPack
	var contentBytes []byte

	err := m.DB.QueryRowContext(ctx, query, id, userID).Scan(
		&pack.ID,
		&pack.UserID,
		&pack.DateRangeStart,
		&pack.DateRangeEnd,
		&contentBytes,
		&pack.JournalCount,
		&pack.PersonalNotes,
		&pack.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}

	pack.Content = json.RawMessage(contentBytes)
	return &pack, nil
}

// GetAllByUser retrieves all prep packs for a user, newest first
func (m PrepPackModel) GetAllByUser(userID string) ([]*PrepPack, error) {
	query := `
		SELECT id, user_id, date_range_start, date_range_end, content, journal_count, personal_notes, created_at
		FROM prep_packs
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var packs []*PrepPack
	for rows.Next() {
		var pack PrepPack
		var contentBytes []byte

		err = rows.Scan(
			&pack.ID,
			&pack.UserID,
			&pack.DateRangeStart,
			&pack.DateRangeEnd,
			&contentBytes,
			&pack.JournalCount,
			&pack.PersonalNotes,
			&pack.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		pack.Content = json.RawMessage(contentBytes)
		packs = append(packs, &pack)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return packs, nil
}

// Delete removes a prep pack
func (m PrepPackModel) Delete(id uuid.UUID, userID string) error {
	query := `
		DELETE FROM prep_packs
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

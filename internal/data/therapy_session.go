package data

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
)

type TherapySession struct {
	ID              uuid.UUID  `json:"id"`
	UserID          string     `json:"user_id"`
	SessionDate     *time.Time `json:"session_date,omitempty"`
	Status          string     `json:"status"`
	MoodBefore      *int       `json:"mood_before,omitempty"`
	TalkingPoints   *string    `json:"talking_points,omitempty"`
	SessionPriority *string    `json:"session_priority,omitempty"`
	PrepPackID      *uuid.UUID `json:"prep_pack_id,omitempty"`
	MoodAfter       *int       `json:"mood_after,omitempty"`
	KeyTakeaways    *string    `json:"key_takeaways,omitempty"`
	SessionRating   *int       `json:"session_rating,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

type TherapySessionModel struct {
	DB *sql.DB
}

// Insert creates a new therapy session
func (m TherapySessionModel) Insert(session *TherapySession) (*TherapySession, error) {
	query := `
		INSERT INTO therapy_sessions (user_id, session_date, status, mood_before, talking_points, session_priority, prep_pack_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, user_id, session_date, status, mood_before, talking_points, session_priority,
		          prep_pack_id, mood_after, key_takeaways, session_rating, created_at, updated_at
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query,
		session.UserID,
		session.SessionDate,
		session.Status,
		session.MoodBefore,
		session.TalkingPoints,
		session.SessionPriority,
		session.PrepPackID,
	).Scan(
		&session.ID,
		&session.UserID,
		&session.SessionDate,
		&session.Status,
		&session.MoodBefore,
		&session.TalkingPoints,
		&session.SessionPriority,
		&session.PrepPackID,
		&session.MoodAfter,
		&session.KeyTakeaways,
		&session.SessionRating,
		&session.CreatedAt,
		&session.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return session, nil
}

// Get retrieves a single therapy session by ID and user
func (m TherapySessionModel) Get(id uuid.UUID, userID string) (*TherapySession, error) {
	query := `
		SELECT id, user_id, session_date, status, mood_before, talking_points, session_priority,
		       prep_pack_id, mood_after, key_takeaways, session_rating, created_at, updated_at
		FROM therapy_sessions
		WHERE id = $1 AND user_id = $2
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var session TherapySession
	err := m.DB.QueryRowContext(ctx, query, id, userID).Scan(
		&session.ID,
		&session.UserID,
		&session.SessionDate,
		&session.Status,
		&session.MoodBefore,
		&session.TalkingPoints,
		&session.SessionPriority,
		&session.PrepPackID,
		&session.MoodAfter,
		&session.KeyTakeaways,
		&session.SessionRating,
		&session.CreatedAt,
		&session.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}

	return &session, nil
}

// GetAllByUser retrieves all therapy sessions for a user, ordered by session_date DESC
func (m TherapySessionModel) GetAllByUser(userID string) ([]*TherapySession, error) {
	query := `
		SELECT id, user_id, session_date, status, mood_before, talking_points, session_priority,
		       prep_pack_id, mood_after, key_takeaways, session_rating, created_at, updated_at
		FROM therapy_sessions
		WHERE user_id = $1
		ORDER BY COALESCE(session_date, created_at) DESC
	`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sessions []*TherapySession
	for rows.Next() {
		var s TherapySession
		err = rows.Scan(
			&s.ID,
			&s.UserID,
			&s.SessionDate,
			&s.Status,
			&s.MoodBefore,
			&s.TalkingPoints,
			&s.SessionPriority,
			&s.PrepPackID,
			&s.MoodAfter,
			&s.KeyTakeaways,
			&s.SessionRating,
			&s.CreatedAt,
			&s.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		sessions = append(sessions, &s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return sessions, nil
}

// Update modifies an existing therapy session (only non-nil fields)
func (m TherapySessionModel) Update(session *TherapySession) error {
	query := `
		UPDATE therapy_sessions
		SET session_date = COALESCE($1, session_date),
		    status = COALESCE($2, status),
		    mood_before = COALESCE($3, mood_before),
		    talking_points = COALESCE($4, talking_points),
		    session_priority = COALESCE($5, session_priority),
		    prep_pack_id = COALESCE($6, prep_pack_id),
		    mood_after = COALESCE($7, mood_after),
		    key_takeaways = COALESCE($8, key_takeaways),
		    session_rating = COALESCE($9, session_rating),
		    updated_at = NOW()
		WHERE id = $10 AND user_id = $11
		RETURNING updated_at
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query,
		session.SessionDate,
		session.Status,
		session.MoodBefore,
		session.TalkingPoints,
		session.SessionPriority,
		session.PrepPackID,
		session.MoodAfter,
		session.KeyTakeaways,
		session.SessionRating,
		session.ID,
		session.UserID,
	).Scan(&session.UpdatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrRecordNotFound
		}
		return err
	}

	return nil
}

// Delete removes a therapy session
func (m TherapySessionModel) Delete(id uuid.UUID, userID string) error {
	query := `
		DELETE FROM therapy_sessions
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

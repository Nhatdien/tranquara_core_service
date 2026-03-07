package data

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

// AIMemory represents a single AI-generated insight about a user.
type AIMemory struct {
	ID               uuid.UUID   `json:"id"`
	UserID           uuid.UUID   `json:"user_id"`
	Content          string      `json:"content"`
	Category         string      `json:"category"`
	SourceJournalIDs []uuid.UUID `json:"source_journal_ids"`
	Confidence       float64     `json:"confidence"`
	CreatedAt        time.Time   `json:"created_at"`
	UpdatedAt        time.Time   `json:"updated_at"`
}

// AIMemoryModel wraps the database connection for ai_memories operations.
type AIMemoryModel struct {
	DB *sql.DB
}

// GetAllByUser returns all memories for a specific user, ordered by creation date (newest first).
// Optionally filters by category.
func (m AIMemoryModel) GetAllByUser(userID uuid.UUID, category string) ([]*AIMemory, error) {
	var query string
	var args []interface{}

	if category != "" {
		query = `
			SELECT id, user_id, content, category, source_journal_ids, confidence, created_at, updated_at
			FROM ai_memories
			WHERE user_id = $1 AND category = $2
			ORDER BY created_at DESC
		`
		args = []interface{}{userID, category}
	} else {
		query = `
			SELECT id, user_id, content, category, source_journal_ids, confidence, created_at, updated_at
			FROM ai_memories
			WHERE user_id = $1
			ORDER BY created_at DESC
		`
		args = []interface{}{userID}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	memories := []*AIMemory{}
	for rows.Next() {
		var mem AIMemory
		err := rows.Scan(
			&mem.ID,
			&mem.UserID,
			&mem.Content,
			&mem.Category,
			pq.Array(&mem.SourceJournalIDs),
			&mem.Confidence,
			&mem.CreatedAt,
			&mem.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		memories = append(memories, &mem)
	}

	return memories, rows.Err()
}

// Delete hard-deletes a single memory by ID and user_id.
func (m AIMemoryModel) Delete(id uuid.UUID, userID uuid.UUID) error {
	query := `DELETE FROM ai_memories WHERE id = $1 AND user_id = $2`

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

// BatchCreate inserts multiple memories for a user. Returns the created records with IDs.
// This is called by the AI service's memory generation scheduler.
func (m AIMemoryModel) BatchCreate(userID uuid.UUID, memories []AIMemory) ([]*AIMemory, error) {
	query := `
		INSERT INTO ai_memories (user_id, content, category, source_journal_ids, confidence)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, user_id, content, category, source_journal_ids, confidence, created_at, updated_at
	`

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	created := []*AIMemory{}
	for _, mem := range memories {
		var result AIMemory
		err := m.DB.QueryRowContext(ctx, query,
			userID,
			mem.Content,
			mem.Category,
			pq.Array(mem.SourceJournalIDs),
			mem.Confidence,
		).Scan(
			&result.ID,
			&result.UserID,
			&result.Content,
			&result.Category,
			pq.Array(&result.SourceJournalIDs),
			&result.Confidence,
			&result.CreatedAt,
			&result.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		created = append(created, &result)
	}

	return created, nil
}

// GetActiveJournalUsersSince returns user IDs that have created/updated journals since the given time.
func (m AIMemoryModel) GetActiveJournalUsersSince(since time.Time) ([]uuid.UUID, error) {
	query := `
		SELECT DISTINCT user_id
		FROM user_journals
		WHERE updated_at >= $1
	`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, since)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var userIDs []uuid.UUID
	for rows.Next() {
		var id uuid.UUID
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		userIDs = append(userIDs, id)
	}

	return userIDs, rows.Err()
}

package data

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

type UserJournal struct {
	ID           uuid.UUID  `json:"id"`
	UserID       uuid.UUID  `json:"user_id"`
	CollectionID *uuid.UUID `json:"collection_id,omitempty"` // Nullable for free-form journals
	Title        string     `json:"title"`
	Content      string     `json:"content"`                // TipTap JSON with embedded emotions + AI
	ContentHTML  *string    `json:"content_html,omitempty"` // Rendered HTML preview
	MoodScore    *int       `json:"mood_score,omitempty"`   // 1-10 scale
	MoodLabel    *string    `json:"mood_label,omitempty"`   // "Storm", "Sunny", etc.
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

// SlideGroup represents a group of slides in a collection
type SlideGroup struct {
	ID          string      `json:"id"`
	Title       string      `json:"title"`
	Description string      `json:"description"`
	Position    int         `json:"position"`
	Slides      []SlideData `json:"slides"`
}

// SlideData represents individual slide configuration
type SlideData struct {
	ID       string                 `json:"id"`
	Type     string                 `json:"type"` // emotion_log, sleep_check, journal_prompt, doc
	Question string                 `json:"question,omitempty"`
	Title    string                 `json:"title,omitempty"`
	Content  string                 `json:"content,omitempty"`
	Config   map[string]interface{} `json:"config,omitempty"`
}

type JournalTemplate struct {
	ID          uuid.UUID       `json:"id"`
	Title       string          `json:"title"`
	Description *string         `json:"description,omitempty"`
	Category    string          `json:"category"`
	SlideGroups json.RawMessage `json:"slide_groups"` // JSONB stored as raw message
	IsActive    bool            `json:"is_active"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}

type UserJournalModel struct {
	DB *sql.DB
}

func (journal UserJournalModel) Get(id uuid.UUID, userID uuid.UUID) (*UserJournal, error) {
	query := `
		SELECT id, user_id, collection_id, title, content, content_html, 
		       mood_score, mood_label, created_at, updated_at 
		FROM user_journals 
		WHERE id = $1 AND user_id = $2  
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var userJournal UserJournal

	err := journal.DB.QueryRowContext(ctx, query, id, userID).Scan(
		&userJournal.ID,
		&userJournal.UserID,
		&userJournal.CollectionID,
		&userJournal.Title,
		&userJournal.Content,
		&userJournal.ContentHTML,
		&userJournal.MoodScore,
		&userJournal.MoodLabel,
		&userJournal.CreatedAt,
		&userJournal.UpdatedAt,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &userJournal, nil
}

func (journal UserJournalModel) GetAllTemplates() ([]*JournalTemplate, error) {
	query := `
		SELECT id, title, description, category, slide_groups, is_active, created_at, updated_at 
		FROM journal_templates
		WHERE is_active = true
		ORDER BY category, title
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	journalTemplates := []*JournalTemplate{}

	rows, err := journal.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var journalTemplate JournalTemplate
		err = rows.Scan(
			&journalTemplate.ID,
			&journalTemplate.Title,
			&journalTemplate.Description,
			&journalTemplate.Category,
			&journalTemplate.SlideGroups,
			&journalTemplate.IsActive,
			&journalTemplate.CreatedAt,
			&journalTemplate.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		journalTemplates = append(journalTemplates, &journalTemplate)
	}

	if err = rows.Err(); err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return journalTemplates, nil
}

func (journal UserJournalModel) GetList(userId uuid.UUID) ([]*UserJournal, error) {
	query := `
		SELECT COUNT(*) OVER(), id, user_id, collection_id, title, content, content_html,
		       mood_score, mood_label, created_at, updated_at 
		FROM user_journals 
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	totalRecords := 0
	userJournals := []*UserJournal{}

	rows, err := journal.DB.QueryContext(ctx, query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var userJournal UserJournal
		err = rows.Scan(
			&totalRecords,
			&userJournal.ID,
			&userJournal.UserID,
			&userJournal.CollectionID,
			&userJournal.Title,
			&userJournal.Content,
			&userJournal.ContentHTML,
			&userJournal.MoodScore,
			&userJournal.MoodLabel,
			&userJournal.CreatedAt,
			&userJournal.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		userJournals = append(userJournals, &userJournal)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return userJournals, nil
}

func (journal UserJournalModel) Insert(userJournal *UserJournal) (*UserJournal, error) {
	query := `
		INSERT INTO user_journals (user_id, collection_id, title, content, content_html, mood_score, mood_label)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, user_id, collection_id, title, content, content_html, mood_score, mood_label, created_at, updated_at
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []any{
		userJournal.UserID,
		userJournal.CollectionID,
		userJournal.Title,
		userJournal.Content,
		userJournal.ContentHTML,
		userJournal.MoodScore,
		userJournal.MoodLabel,
	}

	argsResponse := []any{
		&userJournal.ID,
		&userJournal.UserID,
		&userJournal.CollectionID,
		&userJournal.Title,
		&userJournal.Content,
		&userJournal.ContentHTML,
		&userJournal.MoodScore,
		&userJournal.MoodLabel,
		&userJournal.CreatedAt,
		&userJournal.UpdatedAt,
	}

	err := journal.DB.QueryRowContext(ctx, query, args...).Scan(argsResponse...)

	if err != nil {
		return nil, err
	}

	return userJournal, nil
}

func (journal UserJournalModel) Update(userJournal *UserJournal) (*UserJournal, error) {
	query := `
		UPDATE user_journals
		SET title = $1, content = $2, content_html = $3, mood_score = $4, mood_label = $5
		WHERE id = $6
		RETURNING id, user_id, collection_id, title, content, content_html, mood_score, mood_label, created_at, updated_at
	`

	args := []any{
		userJournal.Title,
		userJournal.Content,
		userJournal.ContentHTML,
		userJournal.MoodScore,
		userJournal.MoodLabel,
		userJournal.ID,
	}

	argsResponse := []any{
		&userJournal.ID,
		&userJournal.UserID,
		&userJournal.CollectionID,
		&userJournal.Title,
		&userJournal.Content,
		&userJournal.ContentHTML,
		&userJournal.MoodScore,
		&userJournal.MoodLabel,
		&userJournal.CreatedAt,
		&userJournal.UpdatedAt,
	}

	err := journal.DB.QueryRow(query, args...).Scan(argsResponse...)

	if err != nil {
		return nil, err
	}
	return userJournal, nil
}

func (journal UserJournalModel) Delete(id uuid.UUID) error {
	query := `
			DELETE FROM user_journals
			WHERE id = $1
			`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()
	result, err := journal.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAfftected, err := result.RowsAffected()

	if err != nil {
		return err
	}
	if rowsAfftected == 0 {
		return ErrRecordNotFound
	}

	return nil
}

// GetListWithFilter retrieves journals with advanced filtering, searching, and pagination.
// Uses the new QueryFilter builder pattern for cleaner query construction.
//
// Supports:
//   - Pagination (page, page_size)
//   - Sorting (any column in safelist)
//   - Full-text search via tsvector (title + content_html)
//   - Time range filtering (created_at, updated_at)
//   - Collection filtering
func (journal UserJournalModel) GetListWithFilter(userID uuid.UUID, filter *QueryFilter, collectionID *uuid.UUID) ([]*UserJournal, Metadata, error) {
	// Build the query dynamically based on filter options
	var queryBuilder strings.Builder
	var args []interface{}
	paramIndex := 1

	// Base SELECT with COUNT for pagination
	queryBuilder.WriteString(`
		SELECT COUNT(*) OVER(), id, user_id, collection_id, title, content, content_html,
		       mood_score, mood_label, created_at, updated_at
		FROM user_journals
		WHERE user_id = $1
	`)
	args = append(args, userID)
	paramIndex++

	// Collection filter (optional)
	if collectionID != nil {
		queryBuilder.WriteString(fmt.Sprintf(" AND collection_id = $%d", paramIndex))
		args = append(args, *collectionID)
		paramIndex++
	}

	// Full-text search condition
	if filter.HasSearch() {
		searchSQL, searchArgs := filter.SearchConditionSQL(paramIndex)
		if searchSQL != "" {
			queryBuilder.WriteString(" AND ")
			queryBuilder.WriteString(searchSQL)
			args = append(args, searchArgs...)
			paramIndex++
		}
	}

	// Time range filter
	if filter.HasTimeRange() {
		timeSQL, timeArgs, nextIdx := filter.TimeRangeConditionSQL(paramIndex)
		if timeSQL != "" {
			queryBuilder.WriteString(" AND ")
			queryBuilder.WriteString(timeSQL)
			args = append(args, timeArgs...)
			paramIndex = nextIdx
		}
	}

	// ORDER BY clause
	// If searching, optionally order by relevance first
	if filter.HasSearch() {
		rankSQL := filter.FullTextRankSQL(2) // search query is at position 2 after userID
		if rankSQL != "" {
			queryBuilder.WriteString(fmt.Sprintf(" ORDER BY %s DESC", rankSQL))
			if filter.SortClause() != "" {
				queryBuilder.WriteString(fmt.Sprintf(", %s", filter.SortClause()))
			}
		} else if filter.SortClause() != "" {
			queryBuilder.WriteString(fmt.Sprintf(" ORDER BY %s", filter.SortClause()))
		}
	} else if filter.SortClause() != "" {
		queryBuilder.WriteString(fmt.Sprintf(" ORDER BY %s", filter.SortClause()))
	} else {
		queryBuilder.WriteString(" ORDER BY created_at DESC")
	}

	// Pagination
	queryBuilder.WriteString(fmt.Sprintf(" LIMIT $%d OFFSET $%d", paramIndex, paramIndex+1))
	args = append(args, filter.Limit(), filter.Offset())

	// Execute query
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := journal.DB.QueryContext(ctx, queryBuilder.String(), args...)
	if err != nil {
		return nil, Metadata{}, err
	}
	defer rows.Close()

	totalRecords := 0
	userJournals := []*UserJournal{}

	for rows.Next() {
		var uj UserJournal
		err = rows.Scan(
			&totalRecords,
			&uj.ID,
			&uj.UserID,
			&uj.CollectionID,
			&uj.Title,
			&uj.Content,
			&uj.ContentHTML,
			&uj.MoodScore,
			&uj.MoodLabel,
			&uj.CreatedAt,
			&uj.UpdatedAt,
		)

		if err != nil {
			return nil, Metadata{}, err
		}

		userJournals = append(userJournals, &uj)
	}

	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	metadata := filter.CalculateMetadata(totalRecords)

	return userJournals, metadata, nil
}

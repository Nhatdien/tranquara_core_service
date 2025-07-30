package data

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
)

type UserJournal struct {
	ID         uuid.UUID `json:"id"`
	UserID     uuid.UUID `json:"user_id"`
	Title      string    `json:"title"`
	Content    string    `json:"content"`
	TemplateID uuid.UUID `json:"template_id"`
	CreatedAt  time.Time `json:"created_at"`
}

type UserJournalModel struct {
	DB *sql.DB
}

func (journal UserJournalModel) Get(id uuid.UUID, userID uuid.UUID) (*UserJournal, error) {
	query := `
				SELECT * FROM user_journals 
				WHERE id = $1 AND user_id = $2  
			`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	var userJournal UserJournal

	err := journal.DB.QueryRowContext(ctx, query, id).Scan(
		&userJournal.ID,
		&userJournal.UserID,
		&userJournal.Title,
		&userJournal.Content,
		&userJournal.TemplateID,
		&userJournal.CreatedAt,
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

func (journal UserJournalModel) GetList(userId uuid.UUID, filter TimeFilter) ([]*UserJournal, TimeFilter, error) {
	query := `
				SELECT COUNT(*) OVER(), id, title, content, template_id, created_at FROM user_journals 
				WHERE user_id = $1 AND create_at BETWEEN $2 AND $3 
				ORDER BY created_at ASC
			`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	totalRecords := 0
	userJournals := []*UserJournal{}

	rows, err := journal.DB.QueryContext(ctx, query, userId, filter.StartTime, filter.EndTime)

	if err != nil {
		return nil, TimeFilter{}, err
	}
	for rows.Next() {
		var userJournal UserJournal
		err = rows.Scan(
			&totalRecords,
			&userJournal.ID,
			&userJournal.Title,
			&userJournal.Content,
			&userJournal.TemplateID,
			&userJournal.CreatedAt,
		)

		if err != nil {
			return nil, TimeFilter{}, err
		}

		userJournals = append(userJournals, &userJournal)
	}

	if err = rows.Err(); err != nil {
		return nil, TimeFilter{}, err
	}

	return userJournals, filter, nil
}

func (journal UserJournalModel) Insert(userJournal *UserJournal) (*UserJournal, error) {
	query := `
			INSERT INTO user_journals (user_id, title, content, template_id)
			VALUES ($1, $2, $3, $4)
			RETURNING *
		`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	args := []any{userJournal.UserID, userJournal.Title, userJournal.Content, userJournal.TemplateID}
	argsResponse := []any{
		&userJournal.ID,
		&userJournal.UserID,
		&userJournal.Title,
		&userJournal.Content,
		&userJournal.TemplateID,
		&userJournal.CreatedAt}

	err := journal.DB.QueryRowContext(ctx, query, args...).Scan(argsResponse...)

	if err != nil {
		return nil, err
	}

	return userJournal, nil
}

func (journal UserJournalModel) Update(userJournal *UserJournal) (*UserJournal, error) {
	query := `
			UPDATE user_journals
			SET title = $1, content = $2, template_id = $3
			WHERE id = $4
			RETURNING *
	`

	args := []any{userJournal.Title, userJournal.Content, userJournal.TemplateID}
	argsResponse := []any{
		&userJournal.ID,
		&userJournal.UserID,
		&userJournal.Title,
		&userJournal.Content,
		&userJournal.TemplateID,
		&userJournal.CreatedAt,
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

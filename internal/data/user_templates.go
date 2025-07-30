package data

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
)

type UserTemplate struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

type UserTemplateModel struct {
	DB *sql.DB
}

func (temp UserTemplateModel) Get(id uuid.UUID, userID uuid.UUID) (*UserTemplate, error) {
	query := `
				SELECT * FROM user_templates 
				WHERE id = $1 AND user_id = $2  
			`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	var userTemplate UserTemplate

	err := temp.DB.QueryRowContext(ctx, query, id).Scan(
		&userTemplate.ID,
		&userTemplate.UserID,
		&userTemplate.Title,
		&userTemplate.Content,
		&userTemplate.CreatedAt,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &userTemplate, nil
}

func (journal UserTemplateModel) GetList(userId uuid.UUID, filter TimeFilter) ([]*UserTemplate, TimeFilter, error) {
	query := `
				SELECT COUNT(*) OVER(), id, title, content, created_at FROM user_templates 
				WHERE user_id = $1 AND create_at BETWEEN $2 AND $3 
				ORDER BY created_at ASC
			`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	totalRecords := 0
	usetTemplates := []*UserTemplate{}

	rows, err := journal.DB.QueryContext(ctx, query, userId, filter.StartTime, filter.EndTime)

	if err != nil {
		return nil, TimeFilter{}, err
	}
	for rows.Next() {
		var userTemplate UserTemplate
		err = rows.Scan(
			&totalRecords,
			&userTemplate.ID,
			&userTemplate.Title,
			&userTemplate.Content,
			&userTemplate.CreatedAt,
		)

		if err != nil {
			return nil, TimeFilter{}, err
		}

		usetTemplates = append(usetTemplates, &userTemplate)
	}

	if err = rows.Err(); err != nil {
		return nil, TimeFilter{}, err
	}

	return usetTemplates, filter, nil
}

func (journal UserTemplateModel) Insert(userTemplate *UserTemplate) (*UserTemplate, error) {
	query := `
			INSERT INTO user_temmplates (user_id, title, content)
			VALUES ($1, $2, $3, $4)
			RETURNING *
		`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	args := []any{userTemplate.UserID, userTemplate.Title, userTemplate.Content}
	argsResponse := []any{
		&userTemplate.ID,
		&userTemplate.UserID,
		&userTemplate.Title,
		&userTemplate.Content,
		&userTemplate.CreatedAt}

	err := journal.DB.QueryRowContext(ctx, query, args...).Scan(argsResponse...)

	if err != nil {
		return nil, err
	}

	return userTemplate, nil
}

func (journal UserTemplateModel) Update(userTemplate *UserTemplate) (*UserTemplate, error) {
	query := `
			UPDATE user_journals
			SET title = $1, content = $2
			WHERE id = $3
			RETURNING *
	`

	args := []any{userTemplate.Title, userTemplate.Content}
	argsResponse := []any{
		&userTemplate.ID,
		&userTemplate.UserID,
		&userTemplate.Title,
		&userTemplate.Content,
		&userTemplate.CreatedAt,
	}

	err := journal.DB.QueryRow(query, args...).Scan(argsResponse...)

	if err != nil {
		return nil, err
	}
	return userTemplate, nil
}

func (journal UserTemplateModel) Delete(id uuid.UUID) error {
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

package data

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type GuiderChatlog struct {
	Id         uuid.UUID `json:"id"`
	UserId     uuid.UUID `json:"user_id"`
	SenderType string    `json:"sender_type"`
	Message    string    `json:"message"`
	CreatedAt  time.Time `json:"created_at"`
}

type GuiderChatlogModel struct {
	DB *sql.DB
}

func (chatlog GuiderChatlogModel) GetList(userUuid uuid.UUID, filter Filter) ([]*GuiderChatlog, Metadata, error) {
	query := `SELECT COUNT(*) OVER(), id, user_id, sender_type, message, created_at FROM ai_guider_chatlog 
			  WHERE user_id = $1
			  ORDER BY created_at ASC
			  LIMIT $2 OFFSET $3`

	totalRecords := 0

	context, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var guiderChatlogs []*GuiderChatlog

	rows, err := chatlog.DB.QueryContext(context, query, userUuid, filter.limit(), filter.offset())

	if err != nil {
		return guiderChatlogs, Metadata{}, err
	}

	if rows == nil {
		return guiderChatlogs, Metadata{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var g GuiderChatlog
		err := rows.Scan(
			&totalRecords,
			&g.Id,
			&g.UserId,
			&g.SenderType,
			&g.Message,
			&g.CreatedAt,
		)
		if err != nil {
			return guiderChatlogs, Metadata{}, err
		}
		// Make a copy for the slice
		guiderChatlogs = append(guiderChatlogs, &g)
	}

	metadata := filter.calculateMetadata(totalRecords, filter.Page, filter.PageSize)
	return guiderChatlogs, metadata, err
}

func (chatlog GuiderChatlogModel) Insert(chatLog *GuiderChatlog) (*GuiderChatlog, error) {
	query := `INSERT INTO ai_guider_chatlog  (user_id, sender_type, message)
			  VALUES ($1, $2, $3)
			  RETURNING id, user_id, sender_type, message, created_at`

	argsResponse := []any{&chatLog.Id, &chatLog.UserId, &chatLog.SenderType, &chatLog.Message, &chatLog.CreatedAt}

	context, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := chatlog.DB.QueryRowContext(context, query, chatLog.UserId, chatLog.SenderType, chatLog.Message).Scan(argsResponse...)

	return chatLog, err
}

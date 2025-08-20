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
	JournalId  uuid.UUID `json:"journal_id"`
	SenderType string    `json:"sender_type"`
	Message    string    `json:"message"`
	CreatedAt  time.Time `json:"created_at"`
}

type GuiderChatlogModel struct {
	DB *sql.DB
}

func (chatlog GuiderChatlogModel) GetList(userUuid uuid.UUID, journalId uuid.UUID) ([]*GuiderChatlog, error) {
	query := `SELECT COUNT(*) OVER(), id, user_id, journal_id, sender_type, message, created_at FROM ai_guider_chatlog 
			  WHERE user_id = $1 AND journal_id = $2
			  ORDER BY created_at ASC
			  `

	totalRecords := 0

	context, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var guiderChatlogs []*GuiderChatlog

	rows, err := chatlog.DB.QueryContext(context, query, userUuid, journalId)

	if err != nil {
		return guiderChatlogs, err
	}

	if rows == nil {
		return guiderChatlogs, err
	}
	defer rows.Close()

	for rows.Next() {
		var g GuiderChatlog
		err := rows.Scan(
			&totalRecords,
			&g.Id,
			&g.UserId,
			&g.JournalId,
			&g.SenderType,
			&g.Message,
			&g.CreatedAt,
		)
		if err != nil {
			return guiderChatlogs, err
		}
		// Make a copy for the slice
		guiderChatlogs = append(guiderChatlogs, &g)
	}

	return guiderChatlogs, err
}

func (chatlog GuiderChatlogModel) Insert(chatLog *GuiderChatlog) (*GuiderChatlog, error) {
	query := `INSERT INTO ai_guider_chatlog  (user_id, sender_type, message, journal_id)
			  VALUES ($1, $2, $3, $4)
			  RETURNING id, user_id, sender_type, message, journal_id, created_at`

	argsResponse := []any{&chatLog.Id, &chatLog.UserId, &chatLog.SenderType, &chatLog.Message, &chatLog.JournalId, &chatLog.CreatedAt}

	context, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := chatlog.DB.QueryRowContext(context, query, chatLog.UserId, chatLog.SenderType, chatLog.Message, chatLog.JournalId).Scan(argsResponse...)

	return chatLog, err
}

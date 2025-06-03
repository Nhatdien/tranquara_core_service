package data

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type EmotionLog struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	Emotion   string    `json:"emotion"`
	Source    string    `json:"source"`
	Context   string    `json:"context"`
	CreatedAt time.Time `json:"created_at"`
}

type EmotionLogModel struct {
	DB *sql.DB
}

func (emo EmotionLogModel) GetList(userId uuid.UUID, filter TimeFilter) ([]*EmotionLog, TimeFilter, error) {
	query := fmt.Sprintf(`
				SELECT COUNT(*) OVER(), id, emotion, source, context, created_at FROM emotion_logs 
				WHERE user_id = $1
				AND created_at BETWEEN $2 and $3
				ORDER BY created_at
			`)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	totalRecords := 0
	emotionLogs := []*EmotionLog{}

	rows, err := emo.DB.QueryContext(ctx, query, filter.StartTime, filter.EndTime)

	if err != nil {
		return nil, TimeFilter{}, err
	}
	for rows.Next() {
		var emotionLog EmotionLog
		err = rows.Scan(
			&totalRecords,
			&emotionLog.ID,
			&emotionLog.Emotion,
			&emotionLog.Source,
			&emotionLog.Context,
			&emotionLog.CreatedAt,
		)

		if err != nil {
			return nil, TimeFilter{}, err
		}

		emotionLogs = append(emotionLogs, &emotionLog)
	}

	if err = rows.Err(); err != nil {
		return nil, TimeFilter{}, err
	}

	return emotionLogs, filter, nil
}

func (emo EmotionLogModel) Insert(emotionLog *EmotionLog) (*EmotionLog, error) {
	query := `
			INSERT INTO emotion_logs (title, description, media_link, exercise_type)
			VALUES ($1, $2, $3, $4)
		`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	args := []any{emotionLog.UserID, emotionLog.Emotion, emotionLog.Source, emotionLog.Context}
	argsResponse := []any{
		&emotionLog.ID,
		&emotionLog.Emotion,
		&emotionLog.Source,
		&emotionLog.Context,
		&emotionLog.CreatedAt}

	err := emo.DB.QueryRowContext(ctx, query, args...).Scan(argsResponse...)

	if err != nil {
		return nil, err
	}
	return emotionLog, nil
}

package data

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"time"
)

type UserInformation struct {
	UserID              int64                  `json:"user_id"`
	Age                 int16                  `json:"age"`
	KYCAnswers          map[string]interface{} `json:"kyc_answers"` // handles JSONB
	ProgramMode         string                 `json:"program_mode"`
	DailyReminderTime   string                 `json:"daily_reminder_time"` // or time.Time if parsed
	NotificationEnabled bool                   `json:"notification_enabled"`
}

type UserInformationModel struct {
	DB *sql.DB
}

func (m UserInformationModel) Get(userID int64) (*UserInformation, error) {
	query := `
		SELECT user_id, age, kyc_answers, program_mode, daily_reminder_time, notification_enabled
		FROM user_information
		WHERE user_id = $1
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var info UserInformation
	var kycRaw []byte

	err := m.DB.QueryRowContext(ctx, query, userID).Scan(
		&info.UserID,
		&info.Age,
		&kycRaw,
		&info.ProgramMode,
		&info.DailyReminderTime,
		&info.NotificationEnabled,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}

	err = json.Unmarshal(kycRaw, &info.KYCAnswers)
	if err != nil {
		return nil, err
	}

	return &info, nil
}

func (m UserInformationModel) Insert(info *UserInformation) error {
	query := `
		INSERT INTO user_information (user_id, age, kyc_answers, program_mode, daily_reminder_time, notification_enabled)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING *
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	kycJSON, err := json.Marshal(info.KYCAnswers)
	if err != nil {
		return err
	}

	return m.DB.QueryRowContext(ctx, query,
		info.UserID,
		info.Age,
		kycJSON,
		info.ProgramMode,
		info.DailyReminderTime,
		info.NotificationEnabled,
	).Scan(
		&info.UserID,
		&info.Age,
		&kycJSON,
		&info.ProgramMode,
		&info.DailyReminderTime,
		&info.NotificationEnabled,
	)
}

func (m UserInformationModel) Update(info *UserInformation) error {
	query := `
		UPDATE user_information
		SET age = $1,
			kyc_answers = $2,
			program_mode = $3,
			daily_reminder_time = $4,
			notification_enabled = $5
		WHERE user_id = $6
		RETURNING *
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	kycJSON, err := json.Marshal(info.KYCAnswers)
	if err != nil {
		return err
	}

	return m.DB.QueryRowContext(ctx, query,
		info.Age,
		kycJSON,
		info.ProgramMode,
		info.DailyReminderTime,
		info.NotificationEnabled,
		info.UserID,
	).Scan(
		&info.UserID,
		&info.Age,
		&kycJSON,
		&info.ProgramMode,
		&info.DailyReminderTime,
		&info.NotificationEnabled,
	)
}

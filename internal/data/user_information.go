package data

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
)

type UserInformation struct {
	UserID     uuid.UUID      `json:"user_id"`
	Name       string         `json:"name"`
	Age        int16          `json:"age"`
	Gender     string         `json:"gender"`
	KYCAnswers map[string]any `json:"kyc_answers"` // handles JSONB
	Settings   map[string]any `json:"settings"`
	CreatedAt  time.Time      `json:"created_at"`
}

type UserInformationModel struct {
	DB *sql.DB
}

func (m UserInformationModel) Get(userID uuid.UUID) (*UserInformation, error) {
	query := `
		SELECT user_id, name, age, gender, kyc_answers, settings, created_at
		FROM user_information
		WHERE user_id = $1
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var info UserInformation
	var kycRaw []byte
	var settingRaw []byte

	err := m.DB.QueryRowContext(ctx, query, userID).Scan(
		&info.UserID,
		&info.Name,
		&info.Age,
		&info.Gender,
		&kycRaw,
		&settingRaw,
		&info.CreatedAt,
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

	err = json.Unmarshal(settingRaw, &info.Settings)
	if err != nil {
		return nil, err
	}

	return &info, nil
}

func (m UserInformationModel) Insert(info *UserInformation) error {
	query := `
		INSERT INTO user_information (user_id, name, age, gender, kyc_answers, settings)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING *
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	kycJSON, err := json.Marshal(info.KYCAnswers)
	if err != nil {
		return err
	}

	settingJSON, err := json.Marshal(info.Settings)
	if err != nil {
		return err
	}

	return m.DB.QueryRowContext(ctx, query,
		info.UserID,
		info.Name,
		info.Age,
		info.Gender,
		kycJSON,
		settingJSON,
		info.CreatedAt,
	).Scan(
		&info.UserID,
		&info.Name,
		&info.Age,
		&info.Gender,
		&info.KYCAnswers,
		&info.Settings,
		&info.CreatedAt,
	)
}

func (m UserInformationModel) Update(info *UserInformation) error {
	query := `
		UPDATE user_information
		SET 
			name = $1
			age = $2,
			gender = $3,
			kyc_answers = $4,
			settings = $5
		WHERE user_id = $6
		RETURNING *
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	kycJSON, err := json.Marshal(info.KYCAnswers)
	if err != nil {
		return err
	}

	settingJSON, err := json.Marshal(info.Settings)
	if err != nil {
		return err
	}

	return m.DB.QueryRowContext(ctx, query,
		info.Name,
		info.Age,
		info.Gender,
		kycJSON,
		settingJSON,
		info.UserID,
	).Scan(
		&info.UserID,
		&info.Name,
		&info.Age,
		&info.Gender,
		&info.KYCAnswers,
		&info.Settings,
		&info.CreatedAt,
	)
}

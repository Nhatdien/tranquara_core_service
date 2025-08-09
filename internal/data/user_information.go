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
	AgeRange   string         `json:"age_range"`
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
		SELECT user_id, name, age_range, gender, kyc_answers, settings, created_at
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
		&info.AgeRange,
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
		INSERT INTO user_information (user_id, name, age_range, gender, kyc_answers, settings)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING user_id, name, age_range, gender, kyc_answers, settings, created_at
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

	var kycAnswersByte, settingBytes []byte
	err = m.DB.QueryRowContext(ctx, query,
		info.UserID,
		info.Name,
		info.AgeRange,
		info.Gender,
		kycJSON,
		settingJSON,
	).Scan(
		&info.UserID,
		&info.Name,
		&info.AgeRange,
		&info.Gender,
		&kycAnswersByte,
		&settingBytes,
		&info.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil // or return ErrNotFound
	}
	if err != nil {
		return err
	}
	// Unmarshal raw JSON bytes back into maps
	if err := json.Unmarshal(kycAnswersByte, &info.KYCAnswers); err != nil {
		return err
	}
	if err := json.Unmarshal(settingBytes, &info.Settings); err != nil {
		return err
	}
	return nil
}

func (m UserInformationModel) Update(info *UserInformation) error {
	query := `
		UPDATE user_information
		SET 
			name = $1
			age_range = $2,
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
		info.AgeRange,
		info.Gender,
		kycJSON,
		settingJSON,
		info.UserID,
	).Scan(
		&info.UserID,
		&info.Name,
		&info.AgeRange,
		&info.Gender,
		&info.KYCAnswers,
		&info.Settings,
		&info.CreatedAt,
	)
}

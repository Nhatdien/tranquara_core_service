package data

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type UserStreak struct {
	UserId        uuid.UUID `json:"user_id"`
	CurrentStreak int       `json:"current_streak"`
	LongestStreak int       `json:"longest_streak"`
	LastActive    time.Time `json:"last_active"`
	TotalEntries  int       `json:"total_entries"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type UserStreakModel struct {
	DB *sql.DB
}

func (streak UserStreakModel) Get(userUuid uuid.UUID) (UserStreak, error) {
	query := `SELECT user_id, current_streak, longest_streak, last_active, total_entries, updated_at
			  FROM user_streaks
			  WHERE user_id = $1`

	var userStreak UserStreak

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := streak.DB.QueryRowContext(ctx, query, userUuid).Scan(
		&userStreak.UserId,
		&userStreak.CurrentStreak,
		&userStreak.LongestStreak,
		&userStreak.LastActive,
		&userStreak.TotalEntries,
		&userStreak.UpdatedAt,
	)

	return userStreak, err
}

func (streak UserStreakModel) Insert(userUUID uuid.UUID) error {
	query := `INSERT INTO user_streaks (user_id, last_active)
			  VALUES ($1, CURRENT_TIMESTAMP)
			  ON CONFLICT (user_id) DO NOTHING`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := streak.DB.ExecContext(ctx, query, userUUID)

	return err
}

func (streak UserStreakModel) UpdateOrReset(userUuid uuid.UUID) error {
	query := `WITH updated_values AS (
				SELECT 
					user_id,
					CASE 
						WHEN CURRENT_DATE - last_active::date > 1 THEN 1
						WHEN CURRENT_DATE - last_active::date = 1 THEN current_streak + 1
						ELSE current_streak
					END AS new_current_streak
				FROM user_streaks
				WHERE user_id = $1
			)
			UPDATE user_streaks u
			SET 
				current_streak = uv.new_current_streak,
				longest_streak = GREATEST(uv.new_current_streak, u.longest_streak),
				last_active = CURRENT_TIMESTAMP,
				total_entries = u.total_entries + 1,
				updated_at = CURRENT_TIMESTAMP
			FROM updated_values uv
			WHERE u.user_id = uv.user_id`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := streak.DB.ExecContext(ctx, query, userUuid)

	return err
}

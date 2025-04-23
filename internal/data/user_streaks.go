package data

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type UserStreak struct {
	UserId        uuid.UUID `json:"user_id"`
	CurrentStreak int8      `json:"current_streak"`
	LongestStreak int8      `json:"longest_streak"`
	LastActive    time.Time `json:"last_active"`
}

type UserStreakModel struct {
	DB *sql.DB
}

func (streak UserStreakModel) Get(userUuid uuid.UUID) (UserStreak, error) {
	query := `SELECT * FROM user_streaks
			  WHERE user_id = $1`

	var userStreak UserStreak
	argsResponse := []any{&userStreak.UserId, &userStreak.CurrentStreak, &userStreak.LongestStreak, &userStreak.LastActive}

	context, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := streak.DB.QueryRowContext(context, query, userUuid).Scan(argsResponse...)

	return userStreak, err
}

func (streak UserStreakModel) Insert(userUUID uuid.UUID) (UserStreak, error) {
	query := `INSERT INTO user_streaks (user_id)
			  VALUES ($1)`

	var userStreak UserStreak
	argsResponse := []any{&userStreak.UserId, &userStreak.CurrentStreak, &userStreak.LongestStreak, &userStreak.LastActive}

	context, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := streak.DB.QueryRowContext(context, query, userUUID).Scan(argsResponse...)

	return userStreak, err
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
				last_active = CURRENT_TIMESTAMP
			FROM updated_values uv
			WHERE u.user_id = uv.user_id
`

	context, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := streak.DB.QueryContext(context, query, userUuid)

	return err
}

package data

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type UserStreak struct {
	UserId        uuid.UUID
	CurrentStreak int8
	LongestStreak int8
	LastActive    time.Time
}

type UserStreakModel struct {
	db *sql.DB
}

func (streak UserStreakModel) Get(userUuid uuid.UUID) (UserStreak, error) {
	query := `SELECT * FROM user_streaks
			  WHERE user_id = $1`

	var userStreak UserStreak
	args := []any{userStreak.UserId}
	argsResponse := []any{&userStreak.UserId, &userStreak.CurrentStreak, &userStreak.LongestStreak, &userStreak.LastActive}

	context, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := streak.db.QueryRowContext(context, query, args...).Scan(argsResponse...)

	return userStreak, err
}

func (streak UserStreakModel) Insert() (UserStreak, error) {
	query := `INSERT INTO user_streaks (user_id)
			  VALUES ($1)`

	var userStreak UserStreak
	args := []any{userStreak.UserId}
	argsResponse := []any{&userStreak.UserId, &userStreak.CurrentStreak, &userStreak.LongestStreak, &userStreak.LastActive}

	context, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := streak.db.QueryRowContext(context, query, args...).Scan(argsResponse...)

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

	return streak.db.QueryRowContext(context, query, userUuid).Scan()
}

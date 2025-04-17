package data

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"tranquara.net/internal/validator"
)

type UserCompletedExercise struct {
	Id          int64     `json:"id"`
	UserId      uuid.UUID `json:"user_id"`
	WeekNumber  int       `json:"week_number"`
	DayNumber   int       `json:"day_number"`
	ExerciseId  int       `json:"exercise_id"`
	CompletedAt time.Time `json:"completed_at"`
	Notes       string    `json:"notes"`
}

type UserCompletedExerciseModel struct {
	DB *sql.DB
}

func UserCompleteExercise(v *validator.Validator, uce *UserCompletedExercise) {
	v.Check((uce.WeekNumber > 1) && (uce.WeekNumber < 8), "week_number", "week_number must be between 1 and 7")
	v.Check((uce.DayNumber > 1) && (uce.DayNumber < 8), "day_number", "day_number must be between 1 and 7")
}

func (uce *UserCompletedExerciseModel) Insert(completeExercise *UserCompletedExercise) error {
	query := `INSERT INTO user_completed_exercises (user_id, week_number, day_number, exercise_id, notes)
			 VALUES ($1, $2, $3, $4, $5)
			 RETURNING *`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	args := []any{completeExercise.UserId, completeExercise.WeekNumber, completeExercise.DayNumber, completeExercise.ExerciseId, completeExercise.Notes}
	argsResponse := []any{&completeExercise.Id, &completeExercise.UserId, &completeExercise.WeekNumber,
		&completeExercise.DayNumber, &completeExercise.ExerciseId,
		&completeExercise.CompletedAt, &completeExercise.Notes}

	return uce.DB.QueryRowContext(ctx, query, args...).Scan(argsResponse...)
}

func (e UserCompletedExerciseModel) GetList(fromTime, toTime time.Time, userID uuid.UUID, filter Filter) ([]*UserCompletedExercise, Metadata, error) {
	query := fmt.Sprintf(`
					SELECT COUNT(*) OVER(), user_id, week_number , day_number , exercise_id , completed_at, notes FROM user_completed_exercises 
					WHERE completed_at BETWEEN $1 AND $2
					AND user_id = $3
					ORDER BY %s %s, id DESC
					LIMIT $4 OFFSET $5
				`, filter.sortColumn(), filter.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	totalRecords := 0
	completedExercises := []*UserCompletedExercise{}

	rows, err := e.DB.QueryContext(ctx, query, fromTime, toTime, userID, filter.limit(), filter.offset())

	if err != nil {
		return nil, Metadata{}, err
	}
	for rows.Next() {
		var completedExercise UserCompletedExercise
		err = rows.Scan(
			&totalRecords,
			&completedExercise.UserId, &completedExercise.WeekNumber,
			&completedExercise.DayNumber, &completedExercise.ExerciseId,
			&completedExercise.CompletedAt, &completedExercise.Notes,
		)

		if err != nil {
			return nil, Metadata{}, err
		}

		completedExercises = append(completedExercises, &completedExercise)
	}

	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculateMetadata(totalRecords, filter.Page, filter.PageSize)

	return completedExercises, metadata, nil
}

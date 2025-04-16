package data

import (
	"database/sql"
	"time"

	"tranquara.net/internal/validator"
)

type UserCompletedExercise struct {
	Id          int64     `json:"id"`
	UserId      int64     `json:"user_id"`
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

func (uce *UserCompletedExerciseModel) Insert() (*UserCompletedExercise, error) {
	return nil, nil
}

func (e UserCompletedExerciseModel) GetList(weekNumber int) ([]*UserCompletedExercise, error) {
	return nil, nil
}

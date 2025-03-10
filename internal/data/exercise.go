package data

import (
	"database/sql"
	"errors"
)

type Exercise struct {
	ExerciseID      int64  `json:"exercise_id"`
	Title           string `json:"title"`
	Description     string `json:"description"`
	DurationMinutes uint   `json:"duration_minutes"`
	ExerciseType    string `json:"exercise_type"`
}

type ExerciseModel struct {
	DB *sql.DB
}

func (e ExerciseModel) Insert(exercise *Exercise) error {
	query := `
		INSERT INTO exercises (title, description, duration_minutes, exercise_type)
		VALUES ($1, $2, $3, $4)
		RETURNING *
	`

	args := []any{exercise.Title, exercise.Description, exercise.DurationMinutes, exercise.ExerciseType}
	argsResponse := []any{&exercise.ExerciseID, &exercise.Title, &exercise.Description, &exercise.DurationMinutes, &exercise.ExerciseType}

	return e.DB.QueryRow(query, args...).Scan(argsResponse...)
}

func (e ExerciseModel) Get(id int64) (*Exercise, error) {
	query := `
				SELECT * FROM exercises 
				WHERE exercise_id = $1
			`

	var exercise Exercise

	err := e.DB.QueryRow(query, id).Scan(
		&exercise.ExerciseID,
		&exercise.Title,
		&exercise.Description,
		&exercise.DurationMinutes,
		&exercise.ExerciseType,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &exercise, nil
}
func (e ExerciseModel) Update(exercise *Exercise) error {
	query := `
			UPDATE exercises
			SET title = $1, description = $2, duration_minutes = $3, exercise_type = $4
			WHERE exercise_id = $5
			RETURNING *
	`

	args := []any{exercise.Title, exercise.Description, exercise.DurationMinutes, exercise.ExerciseType, exercise.ExerciseID}
	argsResponse := []any{&exercise.ExerciseID, &exercise.Title, &exercise.Description, &exercise.DurationMinutes, &exercise.ExerciseType}

	return e.DB.QueryRow(query, args...).Scan(argsResponse...)
}
func (e ExerciseModel) Delete(id int64) error {
	return nil
}

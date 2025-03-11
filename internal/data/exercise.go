package data

import (
	"context"
	"database/sql"
	"errors"
	"time"
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

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	args := []any{exercise.Title, exercise.Description, exercise.DurationMinutes, exercise.ExerciseType}
	argsResponse := []any{&exercise.ExerciseID, &exercise.Title, &exercise.Description, &exercise.DurationMinutes, &exercise.ExerciseType}

	return e.DB.QueryRowContext(ctx, query, args...).Scan(argsResponse...)
}

func (e ExerciseModel) Get(id int64) (*Exercise, error) {
	query := `
				SELECT  exercise_id, title, description, duration_minutes, exercise_type  FROM exercises 
				WHERE exercise_id = $1
			`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	var exercise Exercise

	err := e.DB.QueryRowContext(ctx, query, id).Scan(
		&[]byte{},
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
	query := `
			DELETE FROM exercises
			WHERE exercise_id = $1
			`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()
	result, err := e.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAfftected, err := result.RowsAffected()

	if err != nil {
		return err
	}
	if rowsAfftected == 0 {
		return ErrRecordNotFound
	}

	return nil
}

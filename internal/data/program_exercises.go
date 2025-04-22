package data

import (
	"context"
	"database/sql"
	"time"
)

type ProgramExercise struct {
	ID         int64 `json:"id"`
	WeekNumber int   `json:"week_number"`
	DayNumber  int   `json:"day_number"`
	ExerciseID int64 `json:"exercise_id"`
}

type ProgramExerciseModel struct {
	DB *sql.DB
}

func (m ProgramExerciseModel) Insert(pe *ProgramExercise) error {
	query := `
		INSERT INTO program_exercises (week_number, day_number, exercise_id)
		VALUES ($1, $2, $3)
		RETURNING id, week_number, day_number, exercise_id
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query,
		pe.WeekNumber,
		pe.DayNumber,
		pe.ExerciseID,
	).Scan(
		&pe.ID,
		&pe.WeekNumber,
		&pe.DayNumber,
		&pe.ExerciseID,
	)
}

func (m ProgramExerciseModel) Get() ([]ProgramExercise, error) {
	query := `
		SELECT id, week_number, day_number, exercise_id
		FROM program_exercises
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	peArr := []ProgramExercise{}

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var pe ProgramExercise
		err = rows.Scan(
			&pe.ID,
			&pe.WeekNumber,
			&pe.DayNumber,
			&pe.ExerciseID,
		)

		if err != nil {
			return nil, err
		}

		peArr = append(peArr, pe)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return peArr, nil
}

func (m ProgramExerciseModel) Update(pe *ProgramExercise) error {
	query := `
		UPDATE program_exercises
		SET week_number = $1,
			day_number = $2,
			exercise_id = $3
		WHERE id = $4
		RETURNING id, week_number, day_number, exercise_id
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query,
		pe.WeekNumber,
		pe.DayNumber,
		pe.ExerciseID,
		pe.ID,
	).Scan(
		&pe.ID,
		&pe.WeekNumber,
		&pe.DayNumber,
		&pe.ExerciseID,
	)
}

package data

import (
	"database/sql"
)

type ExerciseModel struct {
	DB *sql.DB
}

func (e ExerciseModel) Insert(exercise ExerciseModel) error {
	return nil
}

func (e ExerciseModel) Get(id int64) (*ExerciseModel, error) {
	return nil, nil
}
func (e ExerciseModel) Update(exercise ExerciseModel) error {
	return nil
}
func (e ExerciseModel) Delete(id int64) error {
	return nil
}

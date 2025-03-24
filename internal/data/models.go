package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
)

type Models struct {
	Exercise ExerciseModel
	User     UserModel
	Token    TokenModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Exercise: ExerciseModel{DB: db},
		User:     UserModel{DB: db},
		Token:    TokenModel{DB: db},
	}

}

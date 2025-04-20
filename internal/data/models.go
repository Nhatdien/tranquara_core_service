package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
)

type Models struct {
	Exercise                       ExerciseModel
	User                           UserModel
	UserCompletedExercise          UserCompletedExerciseModel
	UserCompletedSelfGuideActivity UserCompletedSelfGuideActivityModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Exercise:                       ExerciseModel{DB: db},
		User:                           UserModel{DB: db},
		UserCompletedExercise:          UserCompletedExerciseModel{DB: db},
		UserCompletedSelfGuideActivity: UserCompletedSelfGuideActivityModel{DB: db},
	}

}

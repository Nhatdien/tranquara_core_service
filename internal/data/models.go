package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
)

type Models struct {
	Exercise              ExerciseModel
	User                  UserModel
	UserCompletedExercise UserCompletedExerciseModel
	UserInformation       UserInformationModel
	GuiderChatlog         GuiderChatlogModel
	UserStreak            UserStreakModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Exercise:              ExerciseModel{DB: db},
		User:                  UserModel{DB: db},
		UserCompletedExercise: UserCompletedExerciseModel{DB: db},
		UserInformation:       UserInformationModel{DB: db},
		GuiderChatlog:         GuiderChatlogModel{DB: db},
		UserStreak:            UserStreakModel{DB: db},
	}

}

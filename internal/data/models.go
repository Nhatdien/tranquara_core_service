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
	EmotionLog            EmotionLogModel
	UserJournal           UserJournalModel
	UserLearnedSlideGroup UserLearnedSlideGroupModel
	AIMemory              AIMemoryModel
	TherapySession        TherapySessionModel
	HomeworkItem          HomeworkItemModel
	PrepPack              PrepPackModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Exercise:              ExerciseModel{DB: db},
		User:                  UserModel{DB: db},
		UserCompletedExercise: UserCompletedExerciseModel{DB: db},
		UserInformation:       UserInformationModel{DB: db},
		GuiderChatlog:         GuiderChatlogModel{DB: db},
		UserStreak:            UserStreakModel{DB: db},
		EmotionLog:            EmotionLogModel{DB: db},
		UserJournal:           UserJournalModel{DB: db},
		UserLearnedSlideGroup: UserLearnedSlideGroupModel{DB: db},
		AIMemory:              AIMemoryModel{DB: db},
		TherapySession:        TherapySessionModel{DB: db},
		HomeworkItem:          HomeworkItemModel{DB: db},
		PrepPack:              PrepPackModel{DB: db},
	}

}

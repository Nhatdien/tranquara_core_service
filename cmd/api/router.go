package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {

	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundRespond)
	router.MethodNotAllowed = http.HandlerFunc(app.notFoundRespond)

	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)

	// Auth routes (public - no auth middleware)
	router.HandlerFunc(http.MethodPost, "/v1/auth/register", app.registerUserHandler)

	// User sync route (requires auth)
	router.HandlerFunc(http.MethodPost, "/v1/users/sync", app.authMiddleWare(app.syncUserHandler))

	//Exercises handlers
	router.HandlerFunc(http.MethodGet, "/v1/exercise/:id", app.authMiddleWare(app.showExerciseHanlder))
	router.HandlerFunc(http.MethodGet, "/v1/exercise", app.authMiddleWare(app.listExerciseHandler))
	router.HandlerFunc(http.MethodPost, "/v1/exercise", app.authMiddleWare(app.createExerciseHandler))
	router.HandlerFunc(http.MethodPatch, "/v1/exercise/:id", app.authMiddleWare(app.updateExerciseHandler))
	router.HandlerFunc(http.MethodDelete, "/v1/exercise/:id", app.authMiddleWare(app.deleteExerciseHandler))

	//User completed exercise
	router.HandlerFunc(http.MethodPost, "/v1/user_completed_exercise", app.authMiddleWare(app.createUserCompletedExerciseHandler))
	router.HandlerFunc(http.MethodGet, "/v1/user_completed_exercise", app.authMiddleWare(app.listCompletedExerciseHandler))

	// router.HandlerFunc(http.MethodPut, "/v1/program_exercises/:id", app.authMiddleWare(app.updateProgramExerciseHandler))
	//emotion logs routes
	router.HandlerFunc(http.MethodGet, "/v1/emotion_log", app.authMiddleWare(app.GetEmotionLogs))
	router.HandlerFunc(http.MethodPost, "/v1/emotion_log", app.authMiddleWare(app.CreateEmotionLog))

	//User journal routes
	router.HandlerFunc(http.MethodGet, "/v1/journal", app.authMiddleWare(app.GetUserJournal))
	router.HandlerFunc(http.MethodGet, "/v1/journals", app.authMiddleWare(app.GetUserJournals))
	router.HandlerFunc(http.MethodPost, "/v1/journal", app.authMiddleWare(app.CreateUserJournal))
	router.HandlerFunc(http.MethodPut, "/v1/journal", app.authMiddleWare(app.UpdateUserJournal))
	router.HandlerFunc(http.MethodDelete, "/v1/journal", app.authMiddleWare(app.DeleteUserJournal))

	//chat log routes
	router.HandlerFunc(http.MethodGet, "/v1/guider_chatlogs", app.authMiddleWare(app.getChatLogHandler))

	//User info routes
	router.HandlerFunc(http.MethodGet, "/v1/user_information", app.authMiddleWare(app.getUserInformationHandler))
	router.HandlerFunc(http.MethodPost, "/v1/user_information", app.authMiddleWare(app.createUserInformationHandler))
	router.HandlerFunc(http.MethodPut, "/v1/user_information", app.authMiddleWare(app.updateUserInformationHandler))

	//UserTemplate
	router.HandlerFunc(http.MethodGet, "/v1/user-template", app.authMiddleWare(app.GetUserJournals))
	router.HandlerFunc(http.MethodPost, "/v1/user-template", app.authMiddleWare(app.CreateUserJournal))

	router.HandlerFunc(http.MethodGet, "/v1/tempalte-gallary", app.authMiddleWare(app.GetAllTemplates))
	router.HandlerFunc(http.MethodGet, "/v1/user-template/:id", app.authMiddleWare(app.GetUserJournal))
	router.HandlerFunc(http.MethodPut, "/v1/user-template/:id", app.authMiddleWare(app.UpdateUserJournal))
	router.HandlerFunc(http.MethodDelete, "/v1/user-template/:id", app.authMiddleWare(app.DeleteUserJournal))

	// UserStreak routes
	router.HandlerFunc(http.MethodGet, "/v1/user_streaks", app.authMiddleWare(app.getUserStreakHandler))
	router.HandlerFunc(http.MethodPut, "/v1/user_streaks", app.authMiddleWare(app.updateUserStreakHandler))

	// Learned progress routes
	router.HandlerFunc(http.MethodPost, "/v1/learned", app.authMiddleWare(app.CreateLearnedSlideGroup))
	router.HandlerFunc(http.MethodGet, "/v1/learned", app.authMiddleWare(app.GetAllLearned))
	router.HandlerFunc(http.MethodGet, "/v1/learned/:collection_id", app.authMiddleWare(app.GetLearnedByCollection))
	router.HandlerFunc(http.MethodDelete, "/v1/learned/:id", app.authMiddleWare(app.DeleteLearnedSlideGroup))

	// AI Memory routes (public — requires user auth)
	router.HandlerFunc(http.MethodGet, "/v1/ai-memories", app.authMiddleWare(app.listAIMemoriesHandler))
	router.HandlerFunc(http.MethodDelete, "/v1/ai-memories/:id", app.authMiddleWare(app.deleteAIMemoryHandler))

	// AI Memory routes (internal — called by AI service with API key)
	router.HandlerFunc(http.MethodGet, "/v1/internal/active-journal-users", app.internalAuthMiddleware(app.internalGetActiveJournalUsersHandler))
	router.HandlerFunc(http.MethodGet, "/v1/internal/user-journals", app.internalAuthMiddleware(app.internalGetUserJournalsHandler))
	router.HandlerFunc(http.MethodGet, "/v1/internal/ai-memories/:user_id", app.internalAuthMiddleware(app.internalGetUserMemoriesHandler))
	router.HandlerFunc(http.MethodPost, "/v1/internal/ai-memories/batch", app.internalAuthMiddleware(app.internalBatchCreateMemoriesHandler))

	// Therapy Toolkit — Session routes
	router.HandlerFunc(http.MethodPost, "/v1/therapy-sessions", app.authMiddleWare(app.createSessionHandler))
	router.HandlerFunc(http.MethodGet, "/v1/therapy-sessions", app.authMiddleWare(app.listSessionsHandler))
	router.HandlerFunc(http.MethodPut, "/v1/therapy-sessions", app.authMiddleWare(app.updateSessionHandler))
	router.HandlerFunc(http.MethodDelete, "/v1/therapy-sessions", app.authMiddleWare(app.deleteSessionHandler))

	// Therapy Toolkit — Homework routes
	router.HandlerFunc(http.MethodPost, "/v1/homework", app.authMiddleWare(app.createHomeworkHandler))
	router.HandlerFunc(http.MethodGet, "/v1/homework", app.authMiddleWare(app.listHomeworkHandler))
	router.HandlerFunc(http.MethodPatch, "/v1/homework", app.authMiddleWare(app.toggleHomeworkHandler))
	router.HandlerFunc(http.MethodDelete, "/v1/homework", app.authMiddleWare(app.deleteHomeworkHandler))

	return app.recoverPanic(app.rateLimit(router))
}

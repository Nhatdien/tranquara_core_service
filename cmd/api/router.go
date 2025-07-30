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

	//chat log routes
	router.HandlerFunc(http.MethodGet, "/v1/guider_chatlogs", app.authMiddleWare(app.getChatLogHandler))

	//User info routes
	router.HandlerFunc(http.MethodGet, "/v1/user_information", app.authMiddleWare(app.getUserInformationHandler))
	router.HandlerFunc(http.MethodPost, "/v1/user_information", app.authMiddleWare(app.createUserInformationHandler))
	router.HandlerFunc(http.MethodPut, "/v1/user_information", app.authMiddleWare(app.updateUserInformationHandler))

	//user-journals routes
	router.HandlerFunc(http.MethodGet, "/v1/user-journal", app.authMiddleWare(app.GetUserTemplates))
	router.HandlerFunc(http.MethodPost, "/v1/user-journal", app.authMiddleWare(app.CreateUserTemplate))
	router.HandlerFunc(http.MethodGet, "/v1/user-journal/:id", app.authMiddleWare(app.GetUserTemplate))
	router.HandlerFunc(http.MethodPut, "/v1/user-journal/:id", app.authMiddleWare(app.UpdateUserTemplate))
	router.HandlerFunc(http.MethodDelete, "/v1/user-journal/:id", app.authMiddleWare(app.DeleteUserTemplate))

	//UserTemplate
	router.HandlerFunc(http.MethodGet, "/v1/user-template", app.authMiddleWare(app.GetUserJournals))
	router.HandlerFunc(http.MethodPost, "/v1/user-template", app.authMiddleWare(app.CreateUserJournal))
	router.HandlerFunc(http.MethodGet, "/v1/user-template/:id", app.authMiddleWare(app.GetUserJournal))
	router.HandlerFunc(http.MethodPut, "/v1/user-template/:id", app.authMiddleWare(app.UpdateUserJournal))
	router.HandlerFunc(http.MethodDelete, "/v1/user-template/:id", app.authMiddleWare(app.DeleteUserJournal))

	// UserStreak routes
	router.HandlerFunc(http.MethodGet, "/v1/user_streaks", app.authMiddleWare(app.getUserStreakHandler))
	router.HandlerFunc(http.MethodPut, "/v1/user_streaks", app.authMiddleWare(app.updateUserStreakHandler))

	return app.recoverPanic(app.rateLimit(router))
}

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

	//User completed self guide activity
	router.HandlerFunc(http.MethodPost, "/v1/user_self_guided_activitiy", app.authMiddleWare(app.createUserCompletedSelfGuideActivityHandler))
	router.HandlerFunc(http.MethodGet, "/v1/user_self_guided_activitiy", app.authMiddleWare(app.listCompletedSelfGuideActivityHandler))

	// ProgramExercise routes
	router.HandlerFunc(http.MethodGet, "/v1/program_exercises", app.authMiddleWare(app.getProgramExerciseHandler))
	router.HandlerFunc(http.MethodPost, "/v1/program_exercises", app.authMiddleWare(app.createProgramExerciseHandler))
	// router.HandlerFunc(http.MethodPut, "/v1/program_exercises/:id", app.authMiddleWare(app.updateProgramExerciseHandler))

	//User info routes
	router.HandlerFunc(http.MethodGet, "/v1/user_information", app.authMiddleWare(app.getUserInformationHandler))
	router.HandlerFunc(http.MethodPost, "/v1/user_information", app.authMiddleWare(app.createUserInformationHandler))
	router.HandlerFunc(http.MethodPut, "/v1/user_information", app.authMiddleWare(app.updateUserInformationHandler))

	// UserStreak routes
	router.HandlerFunc(http.MethodGet, "/v1/user_streaks", app.authMiddleWare(app.getUserStreakHandler))
	router.HandlerFunc(http.MethodPost, "/v1/user_streaks", app.authMiddleWare(app.createUserStreakHandler))
	router.HandlerFunc(http.MethodPut, "/v1/user_streaks", app.authMiddleWare(app.updateUserStreakHandler))
	//AI guidence handler
	router.HandlerFunc(http.MethodPost, "/v1/provide_guidence", app.ProvideGuidenceHandler)

	return app.recoverPanic(app.rateLimit(app.testPostMiddleWare(router)))
}

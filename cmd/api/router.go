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
	router.HandlerFunc(http.MethodGet, "/v1/exercise/:id", app.showExerciseHanlder)
	router.HandlerFunc(http.MethodGet, "/v1/exercise", app.listExerciseHandler)
	router.HandlerFunc(http.MethodPost, "/v1/exercise", app.createExerciseHandler)
	router.HandlerFunc(http.MethodPatch, "/v1/exercise/:id", app.updateExerciseHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/exercise/:id", app.deleteExerciseHandler)

	//AI guidence handler
	router.HandlerFunc(http.MethodPost, "/v1/provide_guidence", app.ProvideGuidenceHandler)

	return app.recoverPanic(app.rateLimit(router))
}

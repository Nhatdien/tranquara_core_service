package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() *httprouter.Router {

	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundRespond)
	router.MethodNotAllowed = http.HandlerFunc(app.notFoundRespond)

	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)
	router.HandlerFunc(http.MethodPost, "/v1/provide_guidence", app.ProvideGuidenceHandler)

	return router
}

package main

import (
	"fmt"
	"net/http"
)

func (app *application) logError(r *http.Request, err error) {
	app.logger.PrintError(err, map[string]string{
		"request_method": r.Method,
		"request_url":    r.RequestURI,
	})

}

func (app *application) errorResponse(w http.ResponseWriter, r *http.Request, status int, message interface{}) {
	env := map[string]any{"error": message}
	err := app.writeJson(w, status, env, nil)
	if err != nil {
		app.logError(r, err)
		w.WriteHeader(500)
	}
}

func (app *application) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logError(r, err)
	message := "the server encountered a problem and could not process your request:  " + err.Error()
	app.errorResponse(w, r, http.StatusInternalServerError, message)
}

func (app *application) notFoundRespond(w http.ResponseWriter, r *http.Request) {
	message := "Cannot found resource"
	app.errorResponse(w, r, http.StatusNotFound, message)
}

func (app *application) methodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	message := fmt.Sprintf("the %s method is not supported for this resource", r.Method)
	app.errorResponse(w, r, http.StatusMethodNotAllowed, message)
}

func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.errorResponse(w, r, http.StatusBadRequest, err.Error())
}

func (app *application) failedValidationResponse(w http.ResponseWriter, r *http.Request, errors map[string]string) {
	app.errorResponse(w, r, http.StatusUnprocessableEntity, errors)
}

func (app *application) rateLimitExceedResponse(w http.ResponseWriter, r *http.Request) {
	app.errorResponse(w, r, http.StatusTooManyRequests, "The request limit exceeded")
}

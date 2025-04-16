package main

import (
	"net/http"
)

func (app *application) createUserCompletedExerciseHandler(w http.ResponseWriter, r *http.Request) {
	app.writeJson(w, http.StatusOK, app.GetUserFromContext(r.Context()), nil)
}

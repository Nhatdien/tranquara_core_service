package main

import (
	"net/http"
)

func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"status":      "available",
		"environment": app.config.env,
		"version":     version,
	}
	err := app.writeJson(w, 200, data, nil)
	if err != nil {
		app.logError(r, err)
		app.errorResponse(w, r, http.StatusInternalServerError, "The server encountered a problem and could not process your request")
	}
}

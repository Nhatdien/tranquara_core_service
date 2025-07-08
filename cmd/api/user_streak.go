package main

import (
	"errors"
	"net/http"

	"tranquara.net/internal/data"
)

func (app *application) getUserStreakHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.GetUserUUIDFromContext(r.Context())
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	streak, err := app.models.UserStreak.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundRespond(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJson(w, http.StatusOK, envolope{"user_streak": streak}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) updateUserStreakHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.GetUserUUIDFromContext(r.Context())
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	err = app.models.UserStreak.UpdateOrReset(id)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.writeJson(w, http.StatusOK, envolope{"message": "User streak updated."}, nil)
}

package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"tranquara.net/internal/data"
	"tranquara.net/internal/validator"
)

func (app *application) createUserCompletedExerciseHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		UserId     uuid.UUID `json:"user_id"`
		WeekNumber int       `json:"week_number"`
		Duration   int8      `json:"duration"`
		ExerciseId int       `json:"exercise_id"`
		Notes      string    `json:"notes"`
	}

	err := app.readJson(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	userUUID, err := app.GetUserUUIDFromContext(r.Context())
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	completedExercise := &data.UserCompletedExercise{
		UserId:     userUUID,
		Duration:   input.Duration,
		ExerciseId: input.ExerciseId,
	}

	err = app.models.UserCompletedExercise.Insert(completedExercise)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	header := make(http.Header)
	header.Set("Content-Location", fmt.Sprintf("v1/user_completed_exercise/%d", completedExercise.Id))

	err = app.writeJson(w, http.StatusCreated, completedExercise, header)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) listCompletedExerciseHandler(w http.ResponseWriter, r *http.Request) {

	v := validator.New()
	qs := r.URL.Query()

	// Parse time range
	var fromTime, toTime time.Time
	var err error

	fromStr := app.readString(qs, "from_time", "")
	if fromStr == "" {
		fromTime = time.Time{} // zero value
	} else {
		fromTime, err = time.Parse(time.RFC3339, fromStr)
		if err != nil {
			v.AddError("from_time", "must be a valid RFC3339 timestamp")
		}
	}

	toStr := app.readString(qs, "to_time", "")
	if toStr == "" {
		toTime = time.Now()
	} else {
		toTime, err = time.Parse(time.RFC3339, toStr)
		if err != nil {
			v.AddError("to_time", "must be a valid RFC3339 timestamp")
		}
	}

	filter := app.readQueryFilter(qs, v, DefaultFilterOptions(
		"id",
		[]string{"id", "-id"},
	))

	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	userUUID, err := app.GetUserUUIDFromContext(r.Context())
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	completedExercises, metadata, err := app.models.UserCompletedExercise.GetList(fromTime, toTime, userUUID, filter)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJson(w, http.StatusOK, envolope{"metadata": metadata, "user_completed_exercises": completedExercises}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

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
		DayNumber  int       `json:"day_number"`
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
		WeekNumber: input.WeekNumber,
		DayNumber:  input.DayNumber,
		ExerciseId: input.ExerciseId,
		Notes:      input.Notes,
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

	var input struct {
		FromTime time.Time `json:"from_time"`
		ToTime   time.Time `json:"to_time"`
		data.Filter
	}

	fromStr := app.readString(qs, "from_time", "") // fallback to empty string
	var err error
	if fromStr == "" {
		input.FromTime = time.Time{} // zero value
	} else {
		input.FromTime, err = time.Parse(time.RFC3339, fromStr)
		if err != nil {
			// handle invalid time format
		}
	}

	toStr := app.readString(qs, "to_time", "") // fallback to empty string
	if toStr == "" {
		input.ToTime = time.Now() // zero value
	} else {
		input.ToTime, err = time.Parse(time.RFC3339, toStr)
		if err != nil {
			// handle invalid time format
		}
	}

	input.Page = app.readInt(qs, "page", 1, v)
	input.PageSize = app.readInt(qs, "page_size", 20, v)
	input.Sort = app.readString(qs, "sort", "id")

	input.SortSafelist = []string{"id", "-id"}

	if data.ValidateFilter(v, input.Filter); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	userUUID, err := app.GetUserUUIDFromContext(r.Context())
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	completedExercises, metadata, err := app.models.UserCompletedExercise.GetList(input.FromTime, input.ToTime, userUUID, input.Filter)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJson(w, http.StatusOK, envolope{"metadata": metadata, "completed_exercise": completedExercises}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	// Dump the contents of the input struct in a HTTP response.

}

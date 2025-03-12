package main

import (
	"errors"
	"fmt"
	"net/http"

	"tranquara.net/internal/data"
	"tranquara.net/internal/validator"
)

func (app *application) showExerciseHanlder(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundRespond(w, r)
	}

	exercise, err := app.models.Exercise.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundRespond(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJson(w, http.StatusOK, envolope{"exercise": exercise}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

func (app *application) createExerciseHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title           string `json:"title"`
		Description     string `json:"description"`
		DurationMinutes uint   `json:"duration_minutes"`
		ExerciseType    string `json:"exercise_type"`
	}

	err := app.readJson(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	exercise := &data.Exercise{
		Title:           input.Title,
		Description:     input.Description,
		DurationMinutes: input.DurationMinutes,
		ExerciseType:    input.ExerciseType,
	}

	err = app.models.Exercise.Insert(exercise)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	header := make(http.Header)
	header.Set("Content-Location", fmt.Sprintf("v1/exercises/%d", exercise.ExerciseID))

	err = app.writeJson(w, http.StatusCreated, exercise, header)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) updateExerciseHandler(w http.ResponseWriter, r *http.Request) {
	exerciseID, err := app.readIDParam(r)
	if err != nil || exerciseID < 0 {
		app.notFoundRespond(w, r)
	}

	var input struct {
		Title           *string `json:"title"`
		Description     *string `json:"description"`
		DurationMinutes *uint   `json:"duration_minutes"`
		ExerciseType    *string `json:"exercise_type"`
	}

	err = app.readJson(w, r, &input)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

	exercise, err := app.models.Exercise.Get(exerciseID)

	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundRespond(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
	}

	if input.Title != nil {

		exercise.Title = *input.Title
	}

	if input.Description != nil {
		exercise.Description = *input.Description
	}

	if input.DurationMinutes != nil {
		exercise.DurationMinutes = *input.DurationMinutes
	}

	if input.ExerciseType != nil {
		exercise.ExerciseType = *input.ExerciseType
	}

	app.logger.Print(input)
	app.logger.Print(exercise)

	err = app.models.Exercise.Update(exercise)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

	err = app.writeJson(w, http.StatusOK, exercise, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deleteExerciseHandler(w http.ResponseWriter, r *http.Request) {

	exerciseID, err := app.readIDParam(r)
	if err != nil {
		app.notFoundRespond(w, r)
		return
	}

	err = app.models.Exercise.Delete(exerciseID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	fmt.Fprintf(w, "Exercise id: %d deleted", exerciseID)
}

func (app *application) listExerciseHandler(w http.ResponseWriter, r *http.Request) {

	v := validator.New()

	qs := r.URL.Query()

	var input struct {
		Title        string
		ExerciseType string
		data.Filter
	}

	input.Title = app.readString(qs, "title", "")
	input.ExerciseType = app.readString(qs, "exercise_type", "")

	input.Page = app.readInt(qs, "page", 1, v)
	input.PageSize = app.readInt(qs, "page_size", 20, v)
	input.Sort = app.readString(qs, "sort", "exercise_id")

	input.SortSafelist = []string{"exercise_id", "title", "exercise_type", "-exercise_id", "-title", "-exercise_type"}

	if data.ValidateFilter(v, input.Filter); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	exercises, metadata, err := app.models.Exercise.GetList(input.Title, input.ExerciseType, input.Filter)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJson(w, http.StatusOK, envolope{"metadata": metadata, "exercises": exercises}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	// Dump the contents of the input struct in a HTTP response.

}

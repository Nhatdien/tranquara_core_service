package main

import (
	"errors"
	"fmt"
	"net/http"

	"tranquara.net/internal/data"
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
	}

	err = app.writeJson(w, http.StatusOK, exercise, nil)
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
		Title           string `json:"title"`
		Description     string `json:"description"`
		DurationMinutes uint   `json:"duration_minutes"`
		ExerciseType    string `json:"exercise_type"`
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

	exercise.Title = input.Title
	exercise.Description = input.Description
	exercise.DurationMinutes = input.DurationMinutes
	exercise.ExerciseType = input.ExerciseType

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

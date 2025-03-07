package main

import (
	"fmt"
	"net/http"

	"tranquara.net/internal/data"
)

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

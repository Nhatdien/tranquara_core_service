package main

import (
	"errors"
	"fmt"
	"net/http"

	"tranquara.net/internal/data"
)

func (app *application) createProgramExerciseHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		WeekNumber int   `json:"week_number"`
		DayNumber  int   `json:"day_number"`
		ExerciseID int64 `json:"exercise_id"`
	}

	err := app.readJson(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	programExercise := &data.ProgramExercise{
		WeekNumber: input.WeekNumber,
		DayNumber:  input.DayNumber,
		ExerciseID: input.ExerciseID,
	}

	err = app.models.ProgramExercise.Insert(programExercise)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	header := make(http.Header)
	header.Set("Content-Location", fmt.Sprintf("v1/program-exercises/%d", programExercise.ID))

	err = app.writeJson(w, http.StatusCreated, programExercise, header)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) getProgramExerciseHandler(w http.ResponseWriter, r *http.Request) {
	pe, err := app.models.ProgramExercise.Get()
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundRespond(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJson(w, http.StatusOK, envolope{"program_exercise": pe}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// func (app *application) updateProgramExerciseHandler(w http.ResponseWriter, r *http.Request) {
// 	pe, err := app.models.ProgramExercise.Get()
// 	if err != nil {
// 		switch {
// 		case errors.Is(err, data.ErrRecordNotFound):
// 			app.notFoundRespond(w, r)
// 		default:
// 			app.serverErrorResponse(w, r, err)
// 		}
// 		return
// 	}

// 	var input struct {
// 		WeekNumber *int   `json:"week_number"`
// 		DayNumber  *int   `json:"day_number"`
// 		ExerciseID *int64 `json:"exercise_id"`
// 	}

// 	err = app.readJson(w, r, &input)
// 	if err != nil {
// 		app.serverErrorResponse(w, r, err)
// 		return
// 	}

// 	if input.WeekNumber != nil {
// 		pe.WeekNumber = *input.WeekNumber
// 	}
// 	if input.DayNumber != nil {
// 		pe.DayNumber = *input.DayNumber
// 	}
// 	if input.ExerciseID != nil {
// 		pe.ExerciseID = *input.ExerciseID
// 	}

// 	err = app.models.ProgramExercise.Update(pe)
// 	if err != nil {
// 		app.serverErrorResponse(w, r, err)
// 		return
// 	}

// 	err = app.writeJson(w, http.StatusOK, pe, nil)
// 	if err != nil {
// 		app.serverErrorResponse(w, r, err)
// 	}
// }

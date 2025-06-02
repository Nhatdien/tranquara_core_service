package main

import (
	"fmt"
	"net/http"

	"tranquara.net/internal/data"
)

func (app *application) getUserInformationHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.GetUserUUIDFromContext(r.Context())
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	info, err := app.models.UserInformation.Get(id)
	info.UserID = id
	if err != nil {
		if err == data.ErrRecordNotFound {
			app.notFoundRespond(w, r)
		} else {
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJson(w, http.StatusOK, envolope{"user_info": info}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) createUserInformationHandler(w http.ResponseWriter, r *http.Request) {
	var input data.UserInformation

	id, err := app.GetUserUUIDFromContext(r.Context())
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.readJson(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	input.UserID = id
	err = app.models.UserInformation.Insert(&input)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJson(w, http.StatusCreated, envolope{"user_info": input}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) updateUserInformationHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.GetUserUUIDFromContext(r.Context())
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.logger.PrintInfo(fmt.Sprintf("user id: %s called the api", id), nil)
	info, err := app.models.UserInformation.Get(id)
	if err != nil {
		if err == data.ErrRecordNotFound {
			app.notFoundRespond(w, r)
		} else {
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.readJson(w, r, info)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	info.UserID = id

	err = app.models.UserInformation.Update(info)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJson(w, http.StatusOK, envolope{"user_info": info}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

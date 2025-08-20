package main

import (
	"errors"
	"net/http"

	"github.com/google/uuid"
	"tranquara.net/internal/data"
)

func (app *application) GetUserJournal(w http.ResponseWriter, r *http.Request) {
	idParam := r.URL.Query().Get("id")
	if idParam == "" {
		http.Error(w, "Missing 'id' query parameter", http.StatusBadRequest)
		return
	}

	journalID, err := uuid.FromBytes([]byte(idParam))
	if err != nil {
		http.Error(w, "Invalid 'id' format", http.StatusBadRequest)
		return
	}

	id, err := app.GetUserUUIDFromContext(r.Context())
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	userJournal, err := app.models.UserJournal.Get(journalID, id)
	if err != nil {
		if errors.Is(err, data.ErrRecordNotFound) {
			http.NotFound(w, r)
			return
		}
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJson(w, http.StatusOK, userJournal, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) GetUserJournals(w http.ResponseWriter, r *http.Request) {
	userID, err := app.GetUserUUIDFromContext(r.Context())
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	journals, err := app.models.UserJournal.GetList(userID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJson(w, http.StatusOK, envolope{
		"user_journals": journals,
	}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) GetAllTemplates(w http.ResponseWriter, r *http.Request) {
	templates, err := app.models.UserJournal.GetAllTemplates()
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJson(w, http.StatusOK, envolope{
		"templates": templates,
	}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) CreateUserJournal(w http.ResponseWriter, r *http.Request) {
	var input data.UserJournal

	userID, err := app.GetUserUUIDFromContext(r.Context())
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.readJson(w, r, &input)
	if err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	input.UserID = userID

	newJournal, err := app.models.UserJournal.Insert(&input)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJson(w, http.StatusCreated, newJournal, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) UpdateUserJournal(w http.ResponseWriter, r *http.Request) {
	var input data.UserJournal

	err := app.readJson(w, r, &input)
	if err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	if input.ID == uuid.Nil {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	updatedJournal, err := app.models.UserJournal.Update(&input)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJson(w, http.StatusOK, updatedJournal, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) DeleteUserJournal(w http.ResponseWriter, r *http.Request) {
	idParam := r.URL.Query().Get("id")
	if idParam == "" {
		http.Error(w, "Missing 'id' query parameter", http.StatusBadRequest)
		return
	}

	journalID, err := uuid.FromBytes([]byte(idParam))
	if err != nil {
		http.Error(w, "Invalid 'id' format", http.StatusBadRequest)
		return
	}

	err = app.models.UserJournal.Delete(journalID)
	if err != nil {
		if errors.Is(err, data.ErrRecordNotFound) {
			http.NotFound(w, r)
			return
		}
		app.serverErrorResponse(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

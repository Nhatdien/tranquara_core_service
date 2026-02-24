package main

import (
	"errors"
	"net/http"

	"github.com/google/uuid"
	"tranquara.net/internal/data"
	"tranquara.net/internal/validator"
)

func (app *application) GetUserJournal(w http.ResponseWriter, r *http.Request) {
	idParam := r.URL.Query().Get("id")
	if idParam == "" {
		http.Error(w, "Missing 'id' query parameter", http.StatusBadRequest)
		return
	}

	journalID, err := uuid.Parse(idParam)
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

	v := validator.New()
	qs := r.URL.Query()

	// Parse filter using the new builder pattern
	filter := app.readQueryFilter(qs, v, FilterOptions{
		DefaultPage:     1,
		DefaultPageSize: 20,
		DefaultSort:     "-created_at",
		SortSafelist: []string{
			"created_at", "-created_at",
			"updated_at", "-updated_at",
			"title", "-title",
			"mood_score", "-mood_score",
		},
		TsVectorColumn: "search_vector", // Full-text search column
		UseFullText:    true,
		TimeField:      "created_at", // Enable time range filtering
	})

	// Optional collection filter
	var collectionID *uuid.UUID
	if collectionIDStr := app.readString(qs, "collection_id", ""); collectionIDStr != "" {
		parsed, err := uuid.Parse(collectionIDStr)
		if err != nil {
			v.AddError("collection_id", "must be a valid UUID")
		} else {
			collectionID = &parsed
		}
	}

	// Validate filter
	filter.Validate(v)
	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	// Fetch journals with filter
	journals, metadata, err := app.models.UserJournal.GetListWithFilter(userID, filter, collectionID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJson(w, http.StatusOK, envolope{
		"user_journals": journals,
		"metadata":      metadata,
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

	// Auto-update streak when a journal is created
	if streakErr := app.models.UserStreak.UpdateOrReset(userID); streakErr != nil {
		app.logger.PrintError(streakErr, map[string]string{"action": "update_streak_on_journal_create"})
		// Don't fail journal creation if streak update fails
	}

	// Publish journal to AI service for Qdrant indexing (non-blocking)
	app.publishJournalToAI(newJournal)

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

	// Publish updated journal to AI service for Qdrant re-indexing (non-blocking)
	app.publishJournalToAI(updatedJournal)

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

	journalID, err := uuid.Parse(idParam)
	if err != nil {
		http.Error(w, "Invalid 'id' format", http.StatusBadRequest)
		return
	}

	userID, err := app.GetUserUUIDFromContext(r.Context())
	if err != nil {
		app.serverErrorResponse(w, r, err)
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

	// Remove journal vector from Qdrant (non-blocking)
	app.publishJournalDeleteToAI(journalID, userID)

	w.WriteHeader(http.StatusNoContent)
}

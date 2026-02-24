package main

import (
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
	"tranquara.net/internal/data"
)

// CreateLearnedSlideGroup marks a slide group as completed for the authenticated user
// POST /v1/learned
func (app *application) CreateLearnedSlideGroup(w http.ResponseWriter, r *http.Request) {
	userID, err := app.GetUserUUIDFromContext(r.Context())
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	var input struct {
		CollectionID string `json:"collection_id"`
		SlideGroupID string `json:"slide_group_id"`
	}

	err = app.readJson(w, r, &input)
	if err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	if input.CollectionID == "" || input.SlideGroupID == "" {
		http.Error(w, "collection_id and slide_group_id are required", http.StatusBadRequest)
		return
	}

	collectionID, err := uuid.Parse(input.CollectionID)
	if err != nil {
		http.Error(w, "Invalid collection_id format", http.StatusBadRequest)
		return
	}

	learned := &data.UserLearnedSlideGroup{
		UserID:       userID,
		CollectionID: collectionID,
		SlideGroupID: input.SlideGroupID,
	}

	result, err := app.models.UserLearnedSlideGroup.Insert(learned)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJson(w, http.StatusCreated, envolope{
		"learned": result,
	}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// GetLearnedByCollection retrieves all completed slide groups for a collection
// GET /v1/learned/:collection_id
func (app *application) GetLearnedByCollection(w http.ResponseWriter, r *http.Request) {
	userID, err := app.GetUserUUIDFromContext(r.Context())
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	params := httprouter.ParamsFromContext(r.Context())
	collectionIDStr := params.ByName("collection_id")

	collectionID, err := uuid.Parse(collectionIDStr)
	if err != nil {
		http.Error(w, "Invalid collection_id format", http.StatusBadRequest)
		return
	}

	learned, err := app.models.UserLearnedSlideGroup.GetByCollection(userID, collectionID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJson(w, http.StatusOK, envolope{
		"learned": learned,
	}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// GetAllLearned retrieves all completed slide groups for the authenticated user
// GET /v1/learned
func (app *application) GetAllLearned(w http.ResponseWriter, r *http.Request) {
	userID, err := app.GetUserUUIDFromContext(r.Context())
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	learned, err := app.models.UserLearnedSlideGroup.GetAllByUser(userID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJson(w, http.StatusOK, envolope{
		"learned": learned,
	}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// DeleteLearnedSlideGroup removes a completion record
// DELETE /v1/learned/:id
func (app *application) DeleteLearnedSlideGroup(w http.ResponseWriter, r *http.Request) {
	userID, err := app.GetUserUUIDFromContext(r.Context())
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	params := httprouter.ParamsFromContext(r.Context())
	idStr := params.ByName("id")

	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid id format", http.StatusBadRequest)
		return
	}

	err = app.models.UserLearnedSlideGroup.Delete(id, userID)
	if err != nil {
		if errors.Is(err, data.ErrRecordNotFound) {
			app.notFoundRespond(w, r)
			return
		}
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJson(w, http.StatusOK, envolope{
		"message": "learned record deleted",
	}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/google/uuid"
	"tranquara.net/internal/data"
)

// ─── Prep Pack Handlers ────────────────────────────────────

// createPrepPackHandler creates a new prep pack
// POST /v1/prep-packs
func (app *application) createPrepPackHandler(w http.ResponseWriter, r *http.Request) {
	userUUID, err := app.GetUserUUIDFromContext(r.Context())
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	var input struct {
		DateRangeStart string          `json:"date_range_start"`
		DateRangeEnd   string          `json:"date_range_end"`
		Content        json.RawMessage `json:"content"`
		JournalCount   int             `json:"journal_count"`
		PersonalNotes  *string         `json:"personal_notes"`
	}

	err = app.readJson(w, r, &input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if input.DateRangeStart == "" || input.DateRangeEnd == "" {
		http.Error(w, "date_range_start and date_range_end are required", http.StatusBadRequest)
		return
	}

	if input.Content == nil || len(input.Content) == 0 {
		http.Error(w, "content is required", http.StatusBadRequest)
		return
	}

	startDate, err := time.Parse("2006-01-02", input.DateRangeStart)
	if err != nil {
		http.Error(w, "Invalid date_range_start format (expected YYYY-MM-DD)", http.StatusBadRequest)
		return
	}

	endDate, err := time.Parse("2006-01-02", input.DateRangeEnd)
	if err != nil {
		http.Error(w, "Invalid date_range_end format (expected YYYY-MM-DD)", http.StatusBadRequest)
		return
	}

	pack := &data.PrepPack{
		UserID:         userUUID.String(),
		DateRangeStart: startDate,
		DateRangeEnd:   endDate,
		Content:        input.Content,
		JournalCount:   input.JournalCount,
		PersonalNotes:  input.PersonalNotes,
	}

	result, err := app.models.PrepPack.Insert(pack)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJson(w, http.StatusCreated, envolope{"prep_pack": result}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// listPrepPacksHandler retrieves all prep packs for the authenticated user
// GET /v1/prep-packs
func (app *application) listPrepPacksHandler(w http.ResponseWriter, r *http.Request) {
	userUUID, err := app.GetUserUUIDFromContext(r.Context())
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	packs, err := app.models.PrepPack.GetAllByUser(userUUID.String())
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJson(w, http.StatusOK, envolope{"prep_packs": packs}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// getPrepPackHandler retrieves a single prep pack by ID
// GET /v1/prep-packs/:id
func (app *application) getPrepPackHandler(w http.ResponseWriter, r *http.Request) {
	userUUID, err := app.GetUserUUIDFromContext(r.Context())
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "id query parameter is required", http.StatusBadRequest)
		return
	}

	packID, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid id format", http.StatusBadRequest)
		return
	}

	pack, err := app.models.PrepPack.Get(packID, userUUID.String())
	if err != nil {
		if errors.Is(err, data.ErrRecordNotFound) {
			http.Error(w, "Prep pack not found", http.StatusNotFound)
			return
		}
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJson(w, http.StatusOK, envolope{"prep_pack": pack}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// deletePrepPackHandler deletes a prep pack
// DELETE /v1/prep-packs?id=<uuid>
func (app *application) deletePrepPackHandler(w http.ResponseWriter, r *http.Request) {
	userUUID, err := app.GetUserUUIDFromContext(r.Context())
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "id query parameter is required", http.StatusBadRequest)
		return
	}

	packID, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid id format", http.StatusBadRequest)
		return
	}

	err = app.models.PrepPack.Delete(packID, userUUID.String())
	if err != nil {
		if errors.Is(err, data.ErrRecordNotFound) {
			http.Error(w, "Prep pack not found", http.StatusNotFound)
			return
		}
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJson(w, http.StatusOK, envolope{"message": "prep pack deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

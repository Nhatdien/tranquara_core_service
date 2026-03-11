package main

import (
	"errors"
	"net/http"
	"time"

	"github.com/google/uuid"
	"tranquara.net/internal/data"
)

// ─── Therapy Session Handlers ──────────────────────────────

// createSessionHandler creates a new therapy session
// POST /v1/therapy-sessions
func (app *application) createSessionHandler(w http.ResponseWriter, r *http.Request) {
	userUUID, err := app.GetUserUUIDFromContext(r.Context())
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	var input struct {
		SessionDate     *time.Time `json:"session_date"`
		Status          *string    `json:"status"`
		MoodBefore      *int       `json:"mood_before"`
		TalkingPoints   *string    `json:"talking_points"`
		SessionPriority *string    `json:"session_priority"`
		PrepPackID      *string    `json:"prep_pack_id"`
	}

	err = app.readJson(w, r, &input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	status := "scheduled"
	if input.Status != nil && *input.Status != "" {
		status = *input.Status
	}

	session := &data.TherapySession{
		UserID:          userUUID.String(),
		SessionDate:     input.SessionDate,
		Status:          status,
		MoodBefore:      input.MoodBefore,
		TalkingPoints:   input.TalkingPoints,
		SessionPriority: input.SessionPriority,
	}

	if input.PrepPackID != nil && *input.PrepPackID != "" {
		ppID, err := uuid.Parse(*input.PrepPackID)
		if err != nil {
			http.Error(w, "Invalid prep_pack_id format", http.StatusBadRequest)
			return
		}
		session.PrepPackID = &ppID
	}

	result, err := app.models.TherapySession.Insert(session)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJson(w, http.StatusCreated, envolope{"session": result}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// listSessionsHandler retrieves all therapy sessions for the authenticated user
// GET /v1/therapy-sessions
func (app *application) listSessionsHandler(w http.ResponseWriter, r *http.Request) {
	userUUID, err := app.GetUserUUIDFromContext(r.Context())
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	sessions, err := app.models.TherapySession.GetAllByUser(userUUID.String())
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJson(w, http.StatusOK, envolope{"sessions": sessions}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// updateSessionHandler updates an existing therapy session
// PUT /v1/therapy-sessions?id=<uuid>
func (app *application) updateSessionHandler(w http.ResponseWriter, r *http.Request) {
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

	sessionID, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid id format", http.StatusBadRequest)
		return
	}

	var input struct {
		SessionDate     *time.Time `json:"session_date"`
		Status          *string    `json:"status"`
		MoodBefore      *int       `json:"mood_before"`
		TalkingPoints   *string    `json:"talking_points"`
		SessionPriority *string    `json:"session_priority"`
		PrepPackID      *string    `json:"prep_pack_id"`
		MoodAfter       *int       `json:"mood_after"`
		KeyTakeaways    *string    `json:"key_takeaways"`
		SessionRating   *int       `json:"session_rating"`
	}

	err = app.readJson(w, r, &input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	session := &data.TherapySession{
		ID:              sessionID,
		UserID:          userUUID.String(),
		SessionDate:     input.SessionDate,
		Status:          "",
		MoodBefore:      input.MoodBefore,
		TalkingPoints:   input.TalkingPoints,
		SessionPriority: input.SessionPriority,
		MoodAfter:       input.MoodAfter,
		KeyTakeaways:    input.KeyTakeaways,
		SessionRating:   input.SessionRating,
	}

	// Only set status if provided (COALESCE in SQL handles nil)
	if input.Status != nil {
		session.Status = *input.Status
	}

	if input.PrepPackID != nil && *input.PrepPackID != "" {
		ppID, err := uuid.Parse(*input.PrepPackID)
		if err != nil {
			http.Error(w, "Invalid prep_pack_id format", http.StatusBadRequest)
			return
		}
		session.PrepPackID = &ppID
	}

	err = app.models.TherapySession.Update(session)
	if err != nil {
		if errors.Is(err, data.ErrRecordNotFound) {
			http.Error(w, "Session not found", http.StatusNotFound)
			return
		}
		app.serverErrorResponse(w, r, err)
		return
	}

	// Fetch the updated session to return full data
	updated, err := app.models.TherapySession.Get(sessionID, userUUID.String())
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJson(w, http.StatusOK, envolope{"session": updated}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// deleteSessionHandler deletes a therapy session (cascades to homework)
// DELETE /v1/therapy-sessions?id=<uuid>
func (app *application) deleteSessionHandler(w http.ResponseWriter, r *http.Request) {
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

	sessionID, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid id format", http.StatusBadRequest)
		return
	}

	err = app.models.TherapySession.Delete(sessionID, userUUID.String())
	if err != nil {
		if errors.Is(err, data.ErrRecordNotFound) {
			http.Error(w, "Session not found", http.StatusNotFound)
			return
		}
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJson(w, http.StatusOK, envolope{"message": "session deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// ─── Homework Handlers ────────────────────────────────────

// createHomeworkHandler creates a new homework item
// POST /v1/homework
func (app *application) createHomeworkHandler(w http.ResponseWriter, r *http.Request) {
	userUUID, err := app.GetUserUUIDFromContext(r.Context())
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	var input struct {
		SessionID string `json:"session_id"`
		Content   string `json:"content"`
	}

	err = app.readJson(w, r, &input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if input.SessionID == "" || input.Content == "" {
		http.Error(w, "session_id and content are required", http.StatusBadRequest)
		return
	}

	sessionID, err := uuid.Parse(input.SessionID)
	if err != nil {
		http.Error(w, "Invalid session_id format", http.StatusBadRequest)
		return
	}

	// Verify the session belongs to this user
	_, err = app.models.TherapySession.Get(sessionID, userUUID.String())
	if err != nil {
		if errors.Is(err, data.ErrRecordNotFound) {
			http.Error(w, "Session not found", http.StatusNotFound)
			return
		}
		app.serverErrorResponse(w, r, err)
		return
	}

	item := &data.HomeworkItem{
		SessionID: sessionID,
		UserID:    userUUID.String(),
		Content:   input.Content,
	}

	result, err := app.models.HomeworkItem.Insert(item)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJson(w, http.StatusCreated, envolope{"homework": result}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// listHomeworkHandler retrieves homework items
// GET /v1/homework                   → all homework for user
// GET /v1/homework?session_id=<uuid> → homework for specific session
func (app *application) listHomeworkHandler(w http.ResponseWriter, r *http.Request) {
	userUUID, err := app.GetUserUUIDFromContext(r.Context())
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	sessionIDStr := r.URL.Query().Get("session_id")

	var items []*data.HomeworkItem

	if sessionIDStr != "" {
		sessionID, err := uuid.Parse(sessionIDStr)
		if err != nil {
			http.Error(w, "Invalid session_id format", http.StatusBadRequest)
			return
		}
		items, err = app.models.HomeworkItem.GetBySession(sessionID, userUUID.String())
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}
	} else {
		items, err = app.models.HomeworkItem.GetAllByUser(userUUID.String())
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}
	}

	err = app.writeJson(w, http.StatusOK, envolope{"homework": items}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// toggleHomeworkHandler toggles the completed state of a homework item
// PATCH /v1/homework?id=<uuid>
func (app *application) toggleHomeworkHandler(w http.ResponseWriter, r *http.Request) {
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

	itemID, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid id format", http.StatusBadRequest)
		return
	}

	var input struct {
		Completed bool `json:"completed"`
	}

	err = app.readJson(w, r, &input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := app.models.HomeworkItem.Toggle(itemID, userUUID.String(), input.Completed)
	if err != nil {
		if errors.Is(err, data.ErrRecordNotFound) {
			http.Error(w, "Homework item not found", http.StatusNotFound)
			return
		}
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJson(w, http.StatusOK, envolope{"homework": result}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// deleteHomeworkHandler deletes a homework item
// DELETE /v1/homework?id=<uuid>
func (app *application) deleteHomeworkHandler(w http.ResponseWriter, r *http.Request) {
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

	itemID, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid id format", http.StatusBadRequest)
		return
	}

	err = app.models.HomeworkItem.Delete(itemID, userUUID.String())
	if err != nil {
		if errors.Is(err, data.ErrRecordNotFound) {
			http.Error(w, "Homework item not found", http.StatusNotFound)
			return
		}
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJson(w, http.StatusOK, envolope{"message": "homework item deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

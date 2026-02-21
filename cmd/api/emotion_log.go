package main

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"tranquara.net/internal/data"
	"tranquara.net/internal/validator"
)

func (app *application) GetEmotionLogs(w http.ResponseWriter, r *http.Request) {
	v := validator.New()
	qs := r.URL.Query()

	// Parse query parameters
	id, err := app.GetUserUUIDFromContext(r.Context())

	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	startStr := app.readString(qs, "start", "")
	endStr := app.readString(qs, "end", "")

	if startStr == "" || endStr == "" {
		app.badRequestResponse(w, r, nil)
		return
	}

	startTime, err := time.Parse(time.RFC3339, startStr)
	if err != nil {
		v.AddError("start", "must be a valid RFC3339 timestamp")
	}

	endTime, err := time.Parse(time.RFC3339, endStr)
	if err != nil {
		v.AddError("end", "must be a valid RFC3339 timestamp")
	}

	filter := app.readQueryFilter(qs, v, TimeRangeFilterOptions(
		"created_at",
		[]string{"created_at", "-created_at"},
		"created_at",
	))

	// Apply time range to filter
	filter.WithTimeRange(&startTime, &endTime, "created_at")

	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	// Query from database
	logs, metadata, err := app.models.EmotionLog.GetList(id, filter)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJson(w, http.StatusOK, envolope{
		"emotion_logs": logs,
		"metadata":     metadata,
	}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) CreateEmotionLog(w http.ResponseWriter, r *http.Request) {
	var input data.EmotionLog

	id, err := app.GetUserUUIDFromContext(r.Context())
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	err = app.readJson(w, r, &input)
	if err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	input.UserID = id

	if input.UserID == uuid.Nil || input.Emotion == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	newLog, err := app.models.EmotionLog.Insert(&input)
	if err != nil {
		http.Error(w, "Failed to insert emotion log: "+err.Error(), http.StatusInternalServerError)
		return
	}

	err = app.writeJson(w, http.StatusCreated, newLog, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

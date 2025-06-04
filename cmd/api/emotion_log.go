package main

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"tranquara.net/internal/data"
)

func (app *application) GetEmotionLogs(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	id, err := app.GetUserUUIDFromContext(r.Context())

	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	startStr := r.URL.Query().Get("start")
	endStr := r.URL.Query().Get("end")

	if startStr == "" || endStr == "" {
		http.Error(w, "Missing query parameters", http.StatusBadRequest)
		return
	}

	startTime, err := time.Parse(time.RFC3339, startStr)
	if err != nil {
		http.Error(w, "Invalid start time format", http.StatusBadRequest)
		return
	}

	endTime, err := time.Parse(time.RFC3339, endStr)
	if err != nil {
		http.Error(w, "Invalid end time format", http.StatusBadRequest)
		return
	}

	filter := data.TimeFilter{StartTime: startTime, EndTime: endTime}

	// Query from database
	logs, timeFilter, err := app.models.EmotionLog.GetList(id, filter)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJson(w, http.StatusOK, envolope{
		"emotion_logs": logs,
		"filter":       timeFilter,
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

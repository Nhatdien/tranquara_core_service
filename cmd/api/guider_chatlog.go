package main

import (
	"net/http"

	"github.com/google/uuid"
)

//	type UserGuidenceRequest struct {
//		current_week        string
//		chatbot_interaction string
//		emotion_tracking    string
//	}

func (app *application) getChatLogHandler(w http.ResponseWriter, r *http.Request) {
	qs := r.URL.Query()

	id, err := app.GetUserUUIDFromContext(r.Context())
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
	var input struct {
		journalId uuid.UUID
	}

	input.journalId, err = uuid.Parse(app.readString(qs, "journal_id", ""))
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	chatlogs, err := app.models.GuiderChatlog.GetList(id, input.journalId)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJson(w, http.StatusOK, envolope{"chat_logs": chatlogs}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

package main

import (
	"net/http"

	"tranquara.net/internal/data"
	"tranquara.net/internal/validator"
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
		data.Filter
	}

	v := validator.New()

	input.Page = app.readInt(qs, "page", 1, v)
	input.PageSize = app.readInt(qs, "page_size", 20, v)
	input.Sort = app.readString(qs, "sort", "created_at")

	input.SortSafelist = []string{"created_at"}

	if data.ValidateFilter(v, input.Filter); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	chatlogs, metadata, err := app.models.GuiderChatlog.GetList(id, input.Filter)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJson(w, http.StatusOK, envolope{"metadata": metadata, "chat_logs": chatlogs}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

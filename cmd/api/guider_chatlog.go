package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"tranquara.net/internal/data"
	"tranquara.net/internal/pubsub"
	"tranquara.net/internal/validator"
)

//	type UserGuidenceRequest struct {
//		current_week        string
//		chatbot_interaction string
//		emotion_tracking    string
//	}
func (app *application) ProvideGuidenceHandler(w http.ResponseWriter, r *http.Request) {
	// Publish the message to the RabbitMQ

	var input struct {
		CurrentWeek        int    `json:"current_week"`
		ChatbotInteraction string `json:"chatbot_interaction"`
		EmotionTracking    string `json:"emotion_tracking"`
	}

	var aiResponse struct {
		SuggestMindfulnessTip string `json:"suggest_mindfulness_tip"`
		Explaination          string `json:"explaination"`
	}

	err := app.readJson(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	err = pubsub.PublishJson(app.rabbitchannel, "", "ai_tasks", input)
	if err != nil {
		app.errorResponse(w, r, http.StatusInternalServerError, "Failed to publish message")
		return
	}

	// Synchronous Consumer: Wait for the response from the "ai_response" queue
	messages, err := app.rabbitchannel.Consume(
		"ai_response", // queue name
		"",            // consumer
		false,         // auto-ack
		false,         // exclusive
		false,         // no local
		false,         // no wait
		nil,           // arguments
	)
	if err != nil {
		app.errorResponse(w, r, http.StatusInternalServerError, "Failed to consume message")
		return
	}

	// headers := make(http.Header)
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	// w.Header().Set("Connection", "keep-alive")

	flusher, ok := w.(http.Flusher)
	if !ok {
		app.errorResponse(w, r, http.StatusInternalServerError, "Streaming not supported")
		return
	}
	// Wait for the message synchronously
	for msg := range messages {
		if err := json.Unmarshal(msg.Body, &aiResponse); err != nil {
			log.Printf("❌ Error parsing message: %v", err)
			continue
		}
		// Marshal the JSON response
		jsonData, err := json.Marshal(aiResponse)
		if err != nil {
			log.Printf("❌ Error marshalling JSON: %v", err)
			continue
		}

		// Send JSON as SSE data
		fmt.Fprintf(w, "data: %s\n\n", jsonData)
		flusher.Flush()
		log.Printf("✅ Streamed JSON response: %s", jsonData)
	}
}

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

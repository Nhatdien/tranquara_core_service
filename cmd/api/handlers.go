package main

import (
	"net/http"

	"tranquara.net/internal"
)

// type UserGuidenceRequest struct {
// 	current_week        string
// 	chatbot_interaction string
// 	emotion_tracking    string
// }

func (app *application) ProvideGuidenceHandler(w http.ResponseWriter, r *http.Request) {
	// Publish the message to the RabbitMQ

	var input struct {
		CurrentWeek        int    `json:"current_week"`
		ChatbotInteraction string `json:"chatbot_interaction"`
		EmotionTracking    string `json:"emotion_tracking"`
	}

	err := app.readJson(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	err = internal.PublishJson(app.rabbitchannel, "", "ai_tasks", input)
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

	headers := make(http.Header)
	headers.Set("Content-Type", "application/json")
	// Wait for the message synchronously
	select {
	case message := <-messages:
		// Write the received message directly to the response
		w.WriteHeader(http.StatusOK)
		app.writeJson(w, http.StatusOK, message.Body, headers)
		message.Ack(false)
	case <-r.Context().Done():
		http.Error(w, "Request cancelled", http.StatusRequestTimeout)
	}
}

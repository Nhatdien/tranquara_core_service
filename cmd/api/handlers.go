package main

import (
	"encoding/json"
	"io"
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

	err := internal.PublishJson(app.rabbitchannel, "", "ai_tasks", json.Unmarshal(io.ReadAll(r.Body)))
	if err != nil {
		http.Error(w, "Failed to publish message", http.StatusInternalServerError)
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
		http.Error(w, "Failed to consume message", http.StatusInternalServerError)
		return
	}

	// Wait for the message synchronously
	select {
	case message := <-messages:
		// Write the received message directly to the response
		w.WriteHeader(http.StatusOK)
		w.Write(message.Body)
		message.Ack(false)
	case <-r.Context().Done():
		http.Error(w, "Request cancelled", http.StatusRequestTimeout)
	}
}

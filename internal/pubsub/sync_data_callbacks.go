package pubsub

import (
	"encoding/json"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"tranquara.net/internal/data"
	"tranquara.net/internal/jsonlog"
)

func syncDataMessageCallback(message amqp.Delivery, models *data.Models) {
	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)
	var input struct {
		event     string
		payload   any
		timestamp time.Time
	}

	err := json.Unmarshal(message.Body, &input)
	if err != nil {
		logger.PrintError(err, nil)
		return
	}

	if input.event == "user_journal.create" {
		// Step 1: Marshal payload back to JSON
		payloadBytes, err := json.Marshal(input.payload)
		if err != nil {
			logger.PrintError(err, nil)
			return
		}

		// Step 2: Unmarshal into UserJournal
		var journal data.UserJournal
		err = json.Unmarshal(payloadBytes, &journal)
		if err != nil {
			logger.PrintError(err, nil)
			return
		}

		// Step 3: Use the journal object
		_, err = models.UserJournal.Insert(&journal)
		if err != nil {
			logger.PrintError(err, nil)
		}
	}

	if input.event == "emotion_log.create" {
		// Step 1: Marshal payload back to JSON
		payloadBytes, err := json.Marshal(input.payload)
		if err != nil {
			logger.PrintError(err, nil)
			return
		}

		// Step 2: Unmarshal into UserJournal
		var emotionLog data.EmotionLog
		err = json.Unmarshal(payloadBytes, &emotionLog)
		if err != nil {
			logger.PrintError(err, nil)
			return
		}

		// Step 3: Use the journal object
		_, err = models.EmotionLog.Insert(&emotionLog)
		if err != nil {
			logger.PrintError(err, nil)
		}
	}

	if input.event == "chat_log.create" {
		// Step 1: Marshal payload back to JSON
		payloadBytes, err := json.Marshal(input.payload)
		if err != nil {
			logger.PrintError(err, nil)
			return
		}

		// Step 2: Unmarshal into UserJournal
		var chatLog data.GuiderChatlog
		err = json.Unmarshal(payloadBytes, &chatLog)
		if err != nil {
			logger.PrintError(err, nil)
			return
		}

		// Step 3: Use the journal object
		_, err = models.GuiderChatlog.Insert(&chatLog)
		if err != nil {
			logger.PrintError(err, nil)
		}
	}

	message.Ack(false)

}

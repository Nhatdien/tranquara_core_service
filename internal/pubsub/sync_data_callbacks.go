package pubsub

import (
	"encoding/json"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"tranquara.net/internal/data"
	"tranquara.net/internal/jsonlog"
)

type CustomTime struct {
	time.Time
}

const customLayout = "2006-01-02T15:04:05.000000"

func (ct *CustomTime) UnmarshalJSON(b []byte) error {
	s := string(b)
	s = s[1 : len(s)-1] // remove the quotes
	t, err := time.Parse(customLayout, s)
	if err != nil {
		return err
	}
	ct.Time = t
	return nil
}

func syncDataMessageCallback(message amqp.Delivery, models *data.Models) {
	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)
	var input struct {
		Event     string     `json:"event"`
		Timestamp CustomTime `json:"timestamp"`
		Payload   any        `json:"payload"`
	}

	err := json.Unmarshal(message.Body, &input)
	if err != nil {
		logger.PrintError(err, nil)
		return
	}

	logger.PrintInfo("received message", map[string]string{
		"message": string(message.Body),
	})

	if input.Event == "user_journal.create" {
		// Step 1: Marshal payload back to JSON
		payloadBytes, err := json.Marshal(input.Payload)
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

	if input.Event == "emotion_log.create" {
		// Step 1: Marshal payload back to JSON
		payloadBytes, err := json.Marshal(input.Payload)
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

	if input.Event == "chatlog.create" {
		// Step 1: Marshal payload back to JSON
		payloadBytes, err := json.Marshal(input.Payload)
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

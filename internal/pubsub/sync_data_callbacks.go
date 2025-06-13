package pubsub

import (
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
	"tranquara.net/internal/data"
	"tranquara.net/internal/jsonlog"
)

func syncDataMessageCallback(message amqp.Delivery, models *data.Models) {
	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)

	logger.PrintInfo("get message", map[string]string{
		"message": string(message.Body),
	})

	message.Ack(false)
}

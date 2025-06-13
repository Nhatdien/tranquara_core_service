package pubsub

import (
	"context"
	"encoding/json"

	amqp "github.com/rabbitmq/amqp091-go"
	"tranquara.net/internal/data"
)

func PublishJson[T any](ch *amqp.Channel, exchange, key string, val T) error {
	dat, err := json.Marshal(val)

	if err != nil {
		return err
	}
	return ch.PublishWithContext(context.Background(), exchange, key, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        []byte(dat),
	})
}

func Consumer(ch *amqp.Channel, queue_name string, models *data.Models, callback func(message amqp.Delivery, models *data.Models)) error {
	messages, err := ch.Consume(
		queue_name, // queue name
		"",         // consumer
		false,      // auto-ack
		false,      // exclusive
		false,      // no local
		false,      // no wait
		nil,        // arguments
	)
	if err != nil {
		panic(err)
	}

	go func() {
		for message := range messages {
			// For example, show received message in a console.
			callback(message, models)
		}
	}()

	return err
}

func setupUnits(amqpChannel *amqp.Channel, models *data.Models) error {
	err := defineQueues(amqpChannel)

	if err != nil {
		return err
	}
	err = defineConsumers(amqpChannel, models)

	return err
}

func defineQueues(amqpChannel *amqp.Channel) error {
	_, err := amqpChannel.QueueDeclare("ai_tasks", false, false, false, false, nil)
	if err != nil {
		return err
	}

	_, err = amqpChannel.QueueDeclare("sync_data", false, false, false, false, nil)
	return err
}

func defineConsumers(amqpChannel *amqp.Channel, models *data.Models) error {
	err := Consumer(amqpChannel, "sync_data", models, syncDataMessageCallback)

	if err != nil {
		return err
	}
	err = Consumer(amqpChannel, "ai_tasks", models, aiTaskResponseMessageCallback)

	return err
}

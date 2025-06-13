package pubsub

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"tranquara.net/internal/data"
)

func Serve(models *data.Models) (*amqp.Channel, *amqp.Connection, error) {
	conUrl := "amqp://guest:guest@rabbitmq:5672/"
	conn, err := amqp.Dial(conUrl)
	if err != nil {
		return nil, nil, err
	}

	channel, err := conn.Channel()
	if err != nil {
		return nil, nil, err
	}

	err = setupUnits(channel, models)

	return channel, conn, err
}

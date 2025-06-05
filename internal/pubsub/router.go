package pubsub

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

func Serve() (*amqp.Channel, *amqp.Connection, error) {
	conUrl := "amqp://guest:guest@rabbitmq:5672/"
	conn, err := amqp.Dial(conUrl)
	if err != nil {
		return nil, nil, err
	}

	channel, err := conn.Channel()
	if err != nil {
		return nil, nil, err
	}

	err = setupUnits(channel)

	return channel, conn, err
}

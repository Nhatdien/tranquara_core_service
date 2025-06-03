package pubsub

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

func Serve() (*amqp.Channel, error) {
	conUrl := "amqp://guest:guest@rabbitmq:5672/"
	conn, err := amqp.Dial(conUrl)
	if err != nil {
		panic(err)
	}

	defer conn.Close()
	channel, err := conn.Channel()
	if err != nil {
		panic(err)
	}
	defer channel.Close()

	setupUnits(channel)

	return channel, nil
}

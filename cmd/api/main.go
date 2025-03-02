package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const version = "1.0.0"

type config struct {
	port int
	env  string
}

type application struct {
	config config
	logger *log.Logger
	rabbitchannel *amqp.Channel
}

func main() {
	var cfg config

	// get the arg from cmd
	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "environment", "development", "Environment (development|staging|production)")
	flag.Parse()

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

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

	app := &application{
		config: cfg,
		logger: logger,
		rabbitchannel: channel,
	}


	_, err = channel.QueueDeclare("ai_tasks", false, false, false, false, nil)

	if err != nil {
		log.Println(err)
	}
	_, err = channel.QueueDeclare("ai_response", false, false, false, false, nil)
	if err != nil {
		log.Println(err)
	}

	log.Println("Successfully connected to RabbitMQ")
	log.Println("Waiting for messages")

	mux := http.NewServeMux()
	mux.HandleFunc("/v1/healthcheck", app.healthcheckHandler)
	mux.HandleFunc("/v1/provide_guidence", app.ProvideGuidenceHandler)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      mux,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	logger.Printf("starting %s server on %s", cfg.env, srv.Addr)
	srv.ListenAndServe()
	logger.Fatal(err)
	logger.Printf("server started on %s", srv.Addr)

}

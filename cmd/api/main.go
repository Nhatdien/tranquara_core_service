package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
	amqp "github.com/rabbitmq/amqp091-go"
	"tranquara.net/internal/data"
	"tranquara.net/internal/jsonlog"
)

type envolope map[string]any

const version = "1.0.0"

type config struct {
	port int
	env  string
	db   struct {
		dsn          string
		maxOpenConns int
		maxIdleConns int
		maxIdleTime  string
	}
	limiter struct {
        rps     float64
        burst   int
        enabled bool
    }
}

type application struct {
	config        config
	logger        *jsonlog.Logger
	rabbitchannel *amqp.Channel
	models        data.Models
}

func main() {
	var cfg config

	// get the arg from cmd
	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "environment", "development", "Environment (development|staging|production)")

	flag.StringVar(&cfg.db.dsn, "db-dsn", os.Getenv("TRANQUARA_DB_DSN"), "postgres dsn")

	flag.IntVar(&cfg.db.maxOpenConns, "db-max-open-conns", 25, "Postgres max open connection")
	flag.IntVar(&cfg.db.maxIdleConns, "db-max-idle-conns", 25, "Postgres max idle connection")
	flag.StringVar(&cfg.db.maxIdleTime, "db-max-idle-time", "15m", "Postgres conn max idle time")


	flag.Float64Var(&cfg.limiter.rps, "limiter-rps", 2, "Rate limiter maximum requests per second")
    flag.IntVar(&cfg.limiter.burst, "limiter-burst", 4, "Rate limiter maximum burst")
    flag.BoolVar(&cfg.limiter.enabled, "limiter-enabled", true, "Enable rate limiter")
	flag.Parse()

	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)

	db, err := openDB(cfg)
	if err != nil {
		logger.PrintFatal(err, nil)
	}

	defer db.Close()

	logger.PrintInfo("connect to db successfully", nil)

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
		config:        cfg,
		logger:        logger,
		rabbitchannel: channel,
		models:        data.NewModels(db),
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

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		ErrorLog:     log.New(logger, "", 0),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	logger.PrintInfo("starting server", map[string]string{
		"addr": srv.Addr,
		"env":  cfg.env,
	})
	srv.ListenAndServe()
	logger.PrintFatal(err, nil)
}

func openDB(cfg config) (*sql.DB, error) {
	db, err := sql.Open("postgres", "host=db port=5432 user=postgres password=Nhatdien123 dbname=tranquara_core sslmode=disable")
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cfg.db.maxOpenConns)
	db.SetMaxIdleConns(cfg.db.maxIdleConns)

	idleTime, err := time.ParseDuration(cfg.db.maxIdleTime)

	if err != nil {
		return nil, err
	}

	db.SetConnMaxIdleTime(idleTime)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}
	return db, nil
}

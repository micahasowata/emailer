package main

import (
	"log"
	"net/http"

	"github.com/hibiken/asynq"
)

type app struct {
	client *asynq.Client
}

func main() {
	url := "127.0.0.1:6454"
	client := asynq.NewClient(asynq.RedisClientOpt{
		Addr: url,
	})

	app := &app{
		client: client,
	}
	srv := asynq.NewServer(asynq.RedisClientOpt{
		Addr: url,
	}, asynq.Config{
		Concurrency: 10,
	})

	http.HandleFunc("/", app.sendMail)

	mux := asynq.NewServeMux()
	mux.HandleFunc(typeDelivery, handleDeliveryTask)

	err := srv.Start(mux)
	if err != nil {
		log.Fatal(err.Error())
	}

	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		srv.Shutdown()
		log.Fatal(err.Error())
	}
}

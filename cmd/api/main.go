package main

import (
	"log"
	"log/slog"
	"net/http"

	"github.com/hibiken/asynq"
)

type app struct {
	c *asynq.Client
}

func main() {

	client := asynq.NewClient(asynq.RedisClientOpt{
		Addr: "127.0.0.1:6379",
	})

	app := &app{
		c: client,
	}
	srv := asynq.NewServer(asynq.RedisClientOpt{
		Addr: "127.0.0.1:6379",
	}, asynq.Config{
		Concurrency: 10,
	})

	http.HandleFunc("/", app.sendMail)

	mux := asynq.NewServeMux()
	mux.HandleFunc(TypeEmailDelivery, HandleEmailDeliveryTask)
	err := srv.Start(mux)
	if err != nil {
		slog.Error(err.Error())
		return
	}

	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		srv.Shutdown()
		log.Println("server error", err.Error())
	}
}

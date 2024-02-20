package main

import (
	"log/slog"
	"net/http"

	"github.com/hibiken/asynq"
)

func (app *app) sendMail(w http.ResponseWriter, r *http.Request) {
	task, err := NewEmailTask("barry@flash.com")
	if err != nil {
		slog.Error("task creation failed")
		return
	}

	info, err := app.c.Enqueue(task, asynq.MaxRetry(3))
	if err != nil {
		slog.Error("could not enqueue task")
		return
	}

	slog.Info("starting task queue", slog.String("id", info.ID))
}

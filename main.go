package main

import (
	"log/slog"

	"context"
	"fmt"

	"os"

	"github.com/bytedance/sonic"
	"github.com/hibiken/asynq"
	"github.com/joho/godotenv"

	mailer "github.com/wneessen/go-mail"
)

const (
	TypeEmailDelivery = "email:deliver"
)

type EmailDeliveryPayload struct {
	Email string
}

func NewEmailTask(email string) (*asynq.Task, error) {
	payload, err := sonic.Marshal(EmailDeliveryPayload{Email: email})
	if err != nil {
		return nil, err
	}

	return asynq.NewTask(TypeEmailDelivery, payload, asynq.MaxRetry(3)), nil
}

func HandleEmailDeliveryTask(ctx context.Context, t *asynq.Task) error {
	var p EmailDeliveryPayload

	err := sonic.Unmarshal(t.Payload(), &p)
	if err != nil {
		return fmt.Errorf("sonic: %v : %w", err, asynq.SkipRetry)
	}

	slog.Info("Sending email to user", slog.String("email", p.Email))

	m := mailer.NewMsg()

	err = m.From("toni@stark.com")
	if err != nil {
		return err
	}

	err = m.To(p.Email)
	if err != nil {
		return err
	}

	m.Subject("Hey babbbby, how's Asynq going")
	m.SetBodyString(mailer.TypeTextPlain, "Yooo, from asynq we poppping boy")

	c, err := mailer.NewClient(os.Getenv("HOST"), mailer.WithPort(465), mailer.WithSMTPAuth(mailer.SMTPAuthPlain), mailer.WithUsername(os.Getenv("USERNAME")), mailer.WithPassword(os.Getenv("PASSWORD")))

	if err != nil {
		return err
	}

	err = c.DialAndSend(m)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	client := asynq.NewClient(asynq.RedisClientOpt{
		Addr: "127.0.0.1:6379",
	})

	task, err := NewEmailTask("barry@flash.com")
	if err != nil {
		slog.Error("task creation failed")
		return
	}

	info, err := client.Enqueue(task)
	if err != nil {
		slog.Error("could not enqueue task")
		return
	}

	slog.Info("starting task queue", slog.String("id", info.ID))

	srv := asynq.NewServer(asynq.RedisClientOpt{
		Addr: "127.0.0.1:6379",
	}, asynq.Config{
		Concurrency: 10,
	})

	mux := asynq.NewServeMux()
	mux.HandleFunc(TypeEmailDelivery, HandleEmailDeliveryTask)

	err = srv.Run(mux)
	if err != nil {
		slog.Error(err.Error())
		return
	}
}

package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/bytedance/sonic"
	"github.com/hibiken/asynq"
	mailer "github.com/wneessen/go-mail"
)

const (
	TypeEmailDelivery = "email:deliver"
)

type EmailDeliveryPayload struct {
	Username string
}

func NewEmailTask(username string) (*asynq.Task, error) {
	payload, err := sonic.Marshal(EmailDeliveryPayload{Username: username})
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

	slog.Info("Sending email to user", slog.String("username", p.Username))

	m := mailer.NewMsg()

	err = m.From("toni@stark.com")
	if err != nil {
		return err
	}

	err = m.To("barry@flash.co")
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

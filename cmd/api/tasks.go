package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"

	"github.com/hibiken/asynq"
	mailer "github.com/wneessen/go-mail"
)

const (
	TypeEmailDelivery = "email:deliver"
)

type EmailDeliveryPayload struct {
	Email string
}

func generateEmail(email string) (*mailer.Msg, error) {
	m := mailer.NewMsg()
	if err := m.From("toni.sender@example.com"); err != nil {
		return nil, err
	}
	if err := m.To(email); err != nil {
		return nil, err
	}
	m.Subject("This is my first mail with go-mail!")
	m.SetBodyString(mailer.TypeTextPlain, "Do you like this mail? I certainly do!")
	m.SetBodyString(mailer.TypeTextHTML, "<p>Do you like this email? You no go like am keh</p>")
	return m, nil
}

func NewEmailTask(email string) (*asynq.Task, error) {
	payload, err := json.Marshal(EmailDeliveryPayload{Email: email})
	if err != nil {
		return nil, err
	}

	return asynq.NewTask(TypeEmailDelivery, payload, asynq.MaxRetry(3)), nil
}

func HandleEmailDeliveryTask(ctx context.Context, t *asynq.Task) error {
	var p EmailDeliveryPayload

	err := json.Unmarshal(t.Payload(), &p)
	if err != nil {
		return fmt.Errorf("sonic: %v : %w", err, asynq.SkipRetry)
	}

	slog.Info("Sending email to user", slog.String("email", p.Email))

	m, err := generateEmail(p.Email)
	if err != nil {
		return err
	}

	c, err := mailer.NewClient("sandbox.smtp.mailtrap.io", mailer.WithPort(25), mailer.WithSMTPAuth(mailer.SMTPAuthPlain),
		mailer.WithUsername("47203b77c8bab0"), mailer.WithPassword("32608c369c319e"))
	if err != nil {
		log.Fatalf("failed to create mail client: %s", err)
	}
	if err := c.DialAndSendWithContext(ctx, m); err != nil {
		log.Fatalf("failed to send mail: %s", err)
	}
	return nil
}

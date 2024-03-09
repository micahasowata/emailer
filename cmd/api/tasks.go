package main

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/hibiken/asynq"
	"github.com/wneessen/go-mail"
)

const (
	typeDelivery = "email:deliver"
)

type deliveryPayload struct {
	Receiver string
	Sender   string
}

func newDeliveryTask(receiver, sender string) (*asynq.Task, error) {
	payload, err := json.Marshal(deliveryPayload{Receiver: receiver, Sender: sender})
	if err != nil {
		return nil, err
	}

	return asynq.NewTask(typeDelivery, payload, asynq.MaxRetry(3)), nil
}

func handleDeliveryTask(ctx context.Context, t *asynq.Task) error {
	payload := deliveryPayload{}

	err := json.Unmarshal(t.Payload(), &payload)
	if err != nil {
		return err
	}

	slog.Info("Sending email", slog.String("receiver", payload.Receiver))

	message := mail.NewMsg()

	err = message.From(payload.Sender)
	if err != nil {
		return err
	}

	err = message.To(payload.Receiver)
	if err != nil {
		return err
	}

	message.Subject("Hope you get this!")
	message.SetBodyString(mail.TypeTextPlain, "Do you like this mail? I certainly do!")
	message.SetBodyString(mail.TypeTextHTML, "<p>Do you like this email? You no go like am ke</p>")

	client, err := mail.NewClient("sandbox.smtp.mailtrap.io", mail.WithPort(25), mail.WithSMTPAuth(mail.SMTPAuthPlain), mail.WithUsername("f5f0eac60ad891"), mail.WithPassword("27f809134ac57f"))
	if err != nil {
		return err
	}

	err = client.DialAndSendWithContext(ctx, message)
	if err != nil {
		return err
	}

	return nil
}

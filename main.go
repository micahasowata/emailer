package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"

	"github.com/bytedance/sonic"
	"github.com/hibiken/asynq"

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
	if err := m.From("toni.sender@example.com"); err != nil {
		return err
	}
	if err := m.To(p.Email); err != nil {
		return err
	}
	m.Subject("This is my first mail with go-mail!")
	m.SetBodyString(mailer.TypeTextPlain, "Do you like this mail? I certainly do!")
	m.SetBodyString(mailer.TypeTextHTML, "<p>Do you like this mail? I certainly do!</p>")

	c, err := mailer.NewClient("sandbox.smtp.mailtrap.io", mailer.WithPort(25), mailer.WithSMTPAuth(mailer.SMTPAuthPlain),
		mailer.WithUsername("6a88b61b95374d"), mailer.WithPassword("8eb315375d28f8"))
	if err != nil {
		log.Fatalf("failed to create mail client: %s", err)
	}
	if err := c.DialAndSendWithContext(ctx, m); err != nil {
		log.Fatalf("failed to send mail: %s", err)
	}
	return nil
}

func main() {

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

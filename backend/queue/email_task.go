package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"task-manager/config"

	"github.com/hibiken/asynq"
)

const TaskSendEmail = "email:send"

var asynqClient *asynq.Client

// EmailTaskPayload holds the email task data
type EmailTaskPayload struct {
	Email   string `json:"email"`
	Message string `json:"message"`
}

func InitAsynqClient(redisAddr string) {
	asynqClient = asynq.NewClient(asynq.RedisClientOpt{Addr: redisAddr})
}

func NewEmailTask(email, message string) (*asynq.Task, error) {
	payload, err := json.Marshal(EmailTaskPayload{Email: email, Message: message})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TaskSendEmail, payload), nil
}

func EnqueueEmailTask(payload EmailTaskPayload, runAt time.Time) error {
	if asynqClient == nil {
		return nil // or error - client not initialized
	}
	task, err := NewEmailTask(payload.Email, payload.Message)
	if err != nil {
		return err
	}
	_, err = asynqClient.Enqueue(task, asynq.ProcessAt(runAt))
	if err != nil {
		log.Println("Failed to enqueue email task:", err)
	}
	return err
}

func InitAsynq() {
	go func() {
		srv := asynq.NewServer(
			asynq.RedisClientOpt{Addr: config.RedisUrl},
			asynq.Config{Concurrency: 10},
		)

		mux := asynq.NewServeMux()
		mux.HandleFunc(TaskSendEmail, func(ctx context.Context, t *asynq.Task) error {
			var p EmailTaskPayload
			if err := json.Unmarshal(t.Payload(), &p); err != nil {
				return err
			}
			fmt.Println("Sending email to:", p.Email, "Message:", p.Message)
			// Here: integrate real email sending logic
			return nil
		})

		if err := srv.Run(mux); err != nil {
			fmt.Println("Asynq server error:", err)
		}
	}()
}

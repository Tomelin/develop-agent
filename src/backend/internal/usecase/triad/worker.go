package triad

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type QueuePublisher interface {
	Publish(ctx context.Context, routingKey string, body []byte) error
}

type QueueConsumer interface {
	Consume(queueName string, handler func(amqp.Delivery)) error
}

type Job struct {
	ProjectID   string `json:"project_id"`
	OwnerUserID string `json:"owner_user_id,omitempty"`
	PhaseNumber int    `json:"phase_number,omitempty"`
	Track       string `json:"track,omitempty"`
	Prompt      string `json:"prompt,omitempty"`
	Feedback    string `json:"feedback,omitempty"`
}

type Worker struct {
	QueueName    string
	Publisher    QueuePublisher
	Consumer     QueueConsumer
	Orchestrator *Orchestrator
	OnSuccess    func(ctx context.Context, job Job, refined string) error
}

func (w *Worker) Enqueue(ctx context.Context, routingKey string, job Job) error {
	payload, err := json.Marshal(job)
	if err != nil {
		return fmt.Errorf("marshal job: %w", err)
	}
	return w.Publisher.Publish(ctx, routingKey, payload)
}

func (w *Worker) Start() error {
	if w.QueueName == "" {
		return fmt.Errorf("queue name is required")
	}
	return w.Consumer.Consume(w.QueueName, func(msg amqp.Delivery) {
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		var job Job
		if err := json.Unmarshal(msg.Body, &job); err != nil {
			_ = msg.Nack(false, false)
			return
		}

		refined, err := w.Orchestrator.Run(ctx, ExecutionInput{
			ProjectID: job.ProjectID,
			Prompt:    job.Prompt,
			Feedback:  job.Feedback,
		})
		if err != nil {
			_ = msg.Nack(false, true)
			return
		}
		if w.OnSuccess != nil {
			if err := w.OnSuccess(ctx, job, refined); err != nil {
				_ = msg.Nack(false, true)
				return
			}
		}
		_ = msg.Ack(false)
	})
}

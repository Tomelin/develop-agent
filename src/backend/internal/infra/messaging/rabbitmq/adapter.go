package rabbitmq

import (
	"context"
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Adapter struct {
	Connection *amqp.Connection
	Channel    *amqp.Channel
	url        string
}

// NewAdapter creates a new RabbitMQ adapter and handles connection.
func NewAdapter(url string) (*Adapter, error) {
	a := &Adapter{url: url}
	if err := a.connect(); err != nil {
		return nil, err
	}

	// Start auto-reconnect routine
	go a.handleReconnect()

	return a, nil
}

func (a *Adapter) connect() error {
	conn, err := amqp.Dial(a.url)
	if err != nil {
		return err
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return err
	}

	a.Connection = conn
	a.Channel = ch

	// Setup topology
	if err := a.setupTopology(); err != nil {
		return err
	}

	return nil
}

func (a *Adapter) setupTopology() error {
	// Setup main exchange
	err := a.Channel.ExchangeDeclare(
		"agency.events", // name
		"topic",         // type
		true,            // durable
		false,           // auto-deleted
		false,           // internal
		false,           // no-wait
		nil,             // arguments
	)
	if err != nil {
		return err
	}

	// Setup dead letter exchange and queue
	err = a.Channel.ExchangeDeclare(
		"agency.events.dlx",
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	_, err = a.Channel.QueueDeclare(
		"dlq",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	err = a.Channel.QueueBind(
		"dlq",
		"#",
		"agency.events.dlx",
		false,
		nil,
	)
	if err != nil {
		return err
	}

	// Here you would also create the phase.1 ... phase.9 queues and bind them.
	for i := 1; i <= 9; i++ {
		queueName := fmt.Sprintf("phase.%d", i)
		args := amqp.Table{
			"x-dead-letter-exchange":    "agency.events.dlx",
			"x-dead-letter-routing-key": queueName,
			// Simplified DLQ logic here, real life would need x-delivery-limit setup using Quorum queues or similar
		}

		_, err = a.Channel.QueueDeclare(
			queueName,
			true,
			false,
			false,
			false,
			args,
		)
		if err != nil {
			return err
		}

		if bindErr := a.Channel.QueueBind(
			queueName,
			fmt.Sprintf("phase.%d.*", i), // routing key e.g. phase.1.started
			"agency.events",
			false,
			nil,
		); bindErr != nil {
			return bindErr
		}
	}

	return nil
}

func (a *Adapter) handleReconnect() {
	for {
		reason, ok := <-a.Connection.NotifyClose(make(chan *amqp.Error))
		if !ok {
			// Connection closed gracefully
			return
		}
		log.Printf("RabbitMQ connection closed: %v. Reconnecting...", reason)

		for {
			time.Sleep(5 * time.Second) // Basic exponential backoff could be implemented here
			if err := a.connect(); err == nil {
				log.Println("RabbitMQ reconnected successfully")
				break
			}
		}
	}
}

// Close disconnects gracefully.
func (a *Adapter) Close() error {
	if a.Channel != nil {
		a.Channel.Close()
	}
	if a.Connection != nil {
		return a.Connection.Close()
	}
	return nil
}

// Ping checks if the connection is open.
func (a *Adapter) Ping() (int64, error) {
	start := time.Now()
	if a.Connection == nil || a.Connection.IsClosed() {
		return 0, fmt.Errorf("connection is closed")
	}
	return time.Since(start).Milliseconds(), nil
}

// Publisher is an abstraction for publishing.
type Publisher struct {
	adapter *Adapter
}

func NewPublisher(adapter *Adapter) *Publisher {
	return &Publisher{adapter: adapter}
}

func (p *Publisher) Publish(ctx context.Context, routingKey string, body []byte) error {
	return p.adapter.Channel.PublishWithContext(
		ctx,
		"agency.events",
		routingKey,
		false,
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/json",
			Body:         body,
		},
	)
}

// Consumer is an abstraction for consuming.
type Consumer struct {
	adapter *Adapter
}

func NewConsumer(adapter *Adapter) *Consumer {
	return &Consumer{adapter: adapter}
}

func (c *Consumer) Consume(queueName string, handler func(amqp.Delivery)) error {
	msgs, err := c.adapter.Channel.Consume(
		queueName,
		"",    // consumer tag
		false, // auto-ack disabled (explicit ack needed)
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	go func() {
		for msg := range msgs {
			handler(msg)
		}
	}()

	return nil
}

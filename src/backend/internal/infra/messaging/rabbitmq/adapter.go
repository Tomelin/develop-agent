package rabbitmq

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Adapter struct {
	Connection *amqp.Connection
	Channel    *amqp.Channel
	url        string
	mu         sync.RWMutex
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
		_ = conn.Close()
		return err
	}

	// Setup topology in the new connection/channel first.
	a.mu.Lock()
	a.Connection = conn
	a.Channel = ch
	a.mu.Unlock()

	if err := a.setupTopology(); err != nil {
		_ = ch.Close()
		_ = conn.Close()
		return err
	}

	return nil
}

func (a *Adapter) setupTopology() error {
	a.mu.RLock()
	ch := a.Channel
	a.mu.RUnlock()

	if ch == nil {
		return fmt.Errorf("rabbitmq channel is not initialized")
	}

	// Setup main exchange
	err := ch.ExchangeDeclare(
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
	err = ch.ExchangeDeclare(
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

	_, err = ch.QueueDeclare(
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

	err = ch.QueueBind(
		"dlq",
		"#",
		"agency.events.dlx",
		false,
		nil,
	)
	if err != nil {
		return err
	}

	for i := 1; i <= 9; i++ {
		queueName := fmt.Sprintf("phase.%d", i)
		args := amqp.Table{
			"x-dead-letter-exchange":    "agency.events.dlx",
			"x-dead-letter-routing-key": queueName,
		}

		_, err = ch.QueueDeclare(
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

		if bindErr := ch.QueueBind(
			queueName,
			fmt.Sprintf("phase.%d.*", i),
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
	backoff := 1 * time.Second
	maxBackoff := 30 * time.Second

	for {
		a.mu.RLock()
		conn := a.Connection
		a.mu.RUnlock()

		if conn == nil {
			time.Sleep(backoff)
			continue
		}

		reason, ok := <-conn.NotifyClose(make(chan *amqp.Error))
		if !ok {
			// Connection closed gracefully
			return
		}
		log.Printf("RabbitMQ connection closed: %v. Reconnecting...", reason)

		for {
			time.Sleep(backoff)
			if err := a.connect(); err == nil {
				log.Println("RabbitMQ reconnected successfully")
				backoff = 1 * time.Second
				break
			}
			if backoff < maxBackoff {
				backoff *= 2
				if backoff > maxBackoff {
					backoff = maxBackoff
				}
			}
		}
	}
}

// Close disconnects gracefully.
func (a *Adapter) Close() error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.Channel != nil {
		_ = a.Channel.Close()
	}
	if a.Connection != nil {
		return a.Connection.Close()
	}
	return nil
}

// Ping checks if the connection is open.
func (a *Adapter) Ping() (int64, error) {
	start := time.Now()
	a.mu.RLock()
	defer a.mu.RUnlock()

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
	p.adapter.mu.RLock()
	ch := p.adapter.Channel
	p.adapter.mu.RUnlock()
	if ch == nil {
		return fmt.Errorf("rabbitmq channel is not initialized")
	}

	return ch.PublishWithContext(
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
	c.adapter.mu.RLock()
	ch := c.adapter.Channel
	c.adapter.mu.RUnlock()
	if ch == nil {
		return fmt.Errorf("rabbitmq channel is not initialized")
	}

	msgs, err := ch.Consume(
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

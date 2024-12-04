package queue

import (
	"context"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQ struct {
	conn *amqp.Connection
}

func NewRabbitMQ() *RabbitMQ {
	r := &RabbitMQ{}
	r.Init()
	return r
}

func (r *RabbitMQ) Init() error {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		return err
	}
	r.conn = conn
	return nil
}

func (r *RabbitMQ) Publish(logMessage string) error {
	ch, err := r.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"logs", // name of the queue
		false,  // durable
		false,  // delete when unused
		false,  // exclusive
		false,  // no-wait
		nil,    // arguments
	)
	if err != nil {
		return err
	}

	// Set up a context for timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Publish the message
	err = ch.PublishWithContext(ctx,
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(logMessage),
		})
	if err != nil {
		return err
	}

	log.Printf(" [x] Sent %s\n", logMessage)
	return nil
}

// Close closes the RabbitMQ connection.
func (r *RabbitMQ) Close() {
	if r.conn != nil {
		r.conn.Close()
	}
}

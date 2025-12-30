package rabbit

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/barbodimani81/Event-DDD.git/internal/domain"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitPublisher struct {
	Channel *amqp.Channel
	Queue   amqp.Queue
}

func NewRabbitPublisher(conn *amqp.Connection) (*RabbitPublisher, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("cannot create cahnnel for rabbit: %w", err)
	}

	q, err := ch.QueueDeclare(
		"message_queue",
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		return nil, fmt.Errorf("cannot declare queue for rabbit: %w", err)
	}

	return &RabbitPublisher{
		Channel: ch,
		Queue:   q,
	}, nil
}

func (r *RabbitPublisher) Publish(ctx context.Context, msg domain.Message) error {
	body, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("cannot marshal message: %w", err)
	}

	err = r.Channel.PublishWithContext(
		ctx,
		"",
		r.Queue.Name,
		false,
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/json",
			Body:         body,
		},
	)
	if err != nil {
		return fmt.Errorf("cannot publish on channel: %w", err)
	}

	log.Printf("Event created and published by rabbit by id: %s", msg.UserID)
	return nil
}

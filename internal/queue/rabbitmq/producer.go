package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/rabbitmq/amqp091-go"
)

type TransferEvent struct {
	TransferID    int64  `json:"transfer_id"`
	FromAccountID int64  `json:"from_account_id"`
	ToAccountID   int64  `json:"to_account_id"`
	Amount        int64  `json:"amount"`
	Currency      string `json:"currency"`
}

type Producer interface {
	PublishTransferEvent(ctx context.Context, event TransferEvent) error
	Close()
}

type RabbitProducer struct {
	conn    *amqp091.Connection
	channel *amqp091.Channel
}

func NewRabbitMqProducer(url string) (Producer, error) {
	conn, err := amqp091.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to rabbitmq: %w", err)
	}
	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}
	_, err = ch.QueueDeclare(
		"transfer_event",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}
	return &RabbitProducer{conn: conn, channel: ch}, nil
}
func (r *RabbitProducer) PublishTransferEvent(ctx context.Context, event TransferEvent) error {
	body, err := json.Marshal(event)
	if err != nil {
		return err
	}
	return r.channel.PublishWithContext(ctx,
		"",
		"transfer_event",
		false,
		false,
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}
func (r *RabbitProducer) Close() {
	r.channel.Close()
	r.conn.Close()
}

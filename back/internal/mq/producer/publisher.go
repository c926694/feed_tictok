package producer

import (
	"context"
	"errors"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Publisher struct {
	channel *amqp.Channel
}

func NewPublisher(channel *amqp.Channel) *Publisher {
	return &Publisher{channel: channel}
}

func (p *Publisher) Publish(ctx context.Context, exchange string, routingKey string, msg any) error {
	if p == nil || p.channel == nil {
		return errors.New("publisher channel is nil")
	}
	return p.channel.PublishWithContext(ctx, exchange, routingKey, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        []byte(msg.(string)),
	})
}

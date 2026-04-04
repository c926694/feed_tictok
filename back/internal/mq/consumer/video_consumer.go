package consumer

import (
	"encoding/json"
	"errors"
	"log"
	"simple_tiktok/internal/mq/event"
	"simple_tiktok/internal/pkg/upload"
	"simple_tiktok/internal/repository/mysql"

	amqp "github.com/rabbitmq/amqp091-go"
)

type VideoConsumer struct {
	channel   *amqp.Channel
	videoRepo *mysql.VideoRepo
}

func NewVideoConsumer(channel *amqp.Channel, videoRepo *mysql.VideoRepo) *VideoConsumer {
	return &VideoConsumer{
		channel:   channel,
		videoRepo: videoRepo,
	}
}

func (c *VideoConsumer) Declare(exchange string, exchangeType string, queue string, routingKey string) error {
	if c == nil || c.channel == nil {
		return errors.New("consumer channel is nil")
	}

	if err := c.channel.ExchangeDeclare(exchange, exchangeType, true, false, false, false, nil); err != nil {
		return err
	}

	declaredQueue, err := c.channel.QueueDeclare(queue, true, false, false, false, nil)
	if err != nil {
		return err
	}

	return c.channel.QueueBind(declaredQueue.Name, routingKey, exchange, false, nil)
}

func (c *VideoConsumer) ListenVideoConsumer(queue string, handler func(amqp.Delivery)) error {
	if c == nil || c.channel == nil {
		return errors.New("consumer channel is nil")
	}

	deliveries, err := c.channel.Consume(queue, "", false, false, false, false, nil)
	if err != nil {
		return err
	}

	for delivery := range deliveries {
		handler(delivery)
	}

	return nil
}

func (c *VideoConsumer) DeleteVideoHandler(msg amqp.Delivery) {
	if c == nil || c.channel == nil {
		log.Println("consumer channel is nil")
		_ = msg.Nack(false, false)
		return
	}
	var deleteVideoEvent event.DeleteVideoEvent
	if err := json.Unmarshal(msg.Body, &deleteVideoEvent); err != nil {
		log.Println(err)
		_ = msg.Nack(false, false)
		return
	}
	err := upload.Delete(upload.Video, deleteVideoEvent.PlayURL)
	if err != nil {
		log.Println(err)
		_ = msg.Nack(false, false)
		return
	}
	err = upload.Delete(upload.Cover, deleteVideoEvent.CoverURL)
	if err != nil {
		log.Println(err)
		_ = msg.Nack(false, false)
		return
	}
	_ = msg.Ack(false)
}

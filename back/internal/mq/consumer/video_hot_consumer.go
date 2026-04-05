package consumer

import (
	"encoding/json"
	"errors"
	"log"
	"simple_tiktok/internal/mq/event"
	"simple_tiktok/internal/pkg/constants"
	mysql2 "simple_tiktok/internal/repository/mysql"
	"simple_tiktok/internal/service"

	amqp "github.com/rabbitmq/amqp091-go"
	"gorm.io/gorm"
)

type VideoHotConsumer struct {
	channel      *amqp.Channel
	videoRepo    *mysql2.VideoRepo
	videoService *service.VideoService
}

func NewVideoHotConsumer(channel *amqp.Channel, videoRepo *mysql2.VideoRepo, videoService *service.VideoService) *VideoHotConsumer {
	return &VideoHotConsumer{
		channel:      channel,
		videoRepo:    videoRepo,
		videoService: videoService,
	}
}

func (c *VideoHotConsumer) Declare(exchange string, exchangeType string, queue string, routingKey string) error {
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

func (c *VideoHotConsumer) Listen(queue string, handler func(amqp.Delivery)) error {
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

func (c *VideoHotConsumer) HotUpdateHandler(msg amqp.Delivery) {
	if c == nil || c.channel == nil {
		log.Println("consumer channel is nil")
		_ = msg.Nack(false, false)
		return
	}

	var e event.VideoHotEvent
	if err := json.Unmarshal(msg.Body, &e); err != nil {
		log.Println(err)
		_ = msg.Nack(false, false)
		return
	}

	videoId := e.VideoId
	video, err := c.videoRepo.GetVideoById(videoId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			_ = msg.Ack(false)
			return
		}
		log.Println(err)
		_ = msg.Nack(false, true)
		return
	}

	score := c.videoService.HotScore(video.LikeCount, video.CommentCount, videoId)
	if err := c.videoService.UpdateZSet(constants.HotFeedVideoKey, score, videoId); err != nil {
		log.Println(err)
		_ = msg.Nack(false, true)
		return
	}

	_ = msg.Ack(false)
}

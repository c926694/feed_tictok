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
)

type LikeConsumer struct {
	channel      *amqp.Channel
	videoService *service.VideoService
	videoRepo    *mysql2.VideoRepo
	commentRepo  *mysql2.CommentRepo
}

func NewLikeConsumer(channel *amqp.Channel, videoService *service.VideoService, videoRepo *mysql2.VideoRepo, commentRepo *mysql2.CommentRepo) *LikeConsumer {
	return &LikeConsumer{
		channel:      channel,
		videoService: videoService,
		videoRepo:    videoRepo,
		commentRepo:  commentRepo,
	}
}

func (c *LikeConsumer) Declare(exchange string, exchangeType string, queue string, routingKey string) error {
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

func (c *LikeConsumer) ListenLikeConsumer(queue string, handler func(amqp.Delivery)) error {
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

func (c *LikeConsumer) LikeVideoHandler(msg amqp.Delivery) {
	if c == nil || c.channel == nil {
		log.Println("consumer channel is nil")
		_ = msg.Nack(false, false)
		return
	}
	var videoEvent event.LikeVideoEvent
	if err := json.Unmarshal(msg.Body, &videoEvent); err != nil {
		log.Println(err)
		_ = msg.Nack(false, false)
		return
	}
	videoId := videoEvent.VideoId
	eventType := videoEvent.EventType
	switch eventType {
	case event.Like:
		if err := c.videoRepo.IncVideoLikeCount(videoId); err != nil {
			log.Println(err)
			_ = msg.Nack(false, true)
			return
		}
	case event.Dislike:
		if err := c.videoRepo.DecVideoDislikeCount(videoId); err != nil {
			log.Println(err)
			_ = msg.Nack(false, true)
			return
		}
	default:
		log.Printf("unsupported like event type: %s", eventType)
		_ = msg.Nack(false, false)
		return
	}

	video, err := c.videoRepo.GetVideoById(videoId)
	if err != nil {
		log.Println(err)
		_ = msg.Nack(false, true)
		return
	}

	newScore := c.videoService.HotScore(video.LikeCount, video.CommentCount, videoId)
	if err := c.videoService.UpdateZSet(constants.HotFeedVideoKey, newScore, videoId); err != nil {
		log.Println(err)
		_ = msg.Nack(false, true)
		return
	}

	_ = msg.Ack(false)
}

func (c *LikeConsumer) LikeCommentHandler(msg amqp.Delivery) {
	if c == nil || c.channel == nil {
		log.Println("consumer channel is nil")
		_ = msg.Nack(false, false)
		return
	}
	var commentEvent event.LikeCommentEvent
	if err := json.Unmarshal(msg.Body, &commentEvent); err != nil {
		log.Println(err)
		_ = msg.Nack(false, false)
		return
	}

	commentId := commentEvent.CommentId
	switch commentEvent.EventType {
	case event.Like:
		if err := c.commentRepo.IncCommentLikeCount(commentId); err != nil {
			log.Println(err)
			_ = msg.Nack(false, true)
			return
		}
	case event.Dislike:
		if err := c.commentRepo.DecCommentLikeCount(commentId); err != nil {
			log.Println(err)
			_ = msg.Nack(false, true)
			return
		}
	default:
		log.Printf("unsupported comment like event type: %s", commentEvent.EventType)
		_ = msg.Nack(false, false)
		return
	}

	_ = msg.Ack(false)
}

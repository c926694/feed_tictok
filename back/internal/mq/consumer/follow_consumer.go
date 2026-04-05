package consumer

import (
	"encoding/json"
	"errors"
	"log"
	"simple_tiktok/internal/model"
	"simple_tiktok/internal/mq/event"
	"simple_tiktok/internal/repository/mysql"

	amqp "github.com/rabbitmq/amqp091-go"
	"gorm.io/gorm"
)

type FollowConsumer struct {
	channel    *amqp.Channel
	followRepo *mysql.FollowRepo
	userRepo   *mysql.UserRepo
}

func NewFollowConsumer(channel *amqp.Channel, followRepo *mysql.FollowRepo, userRepo *mysql.UserRepo) *FollowConsumer {
	return &FollowConsumer{channel: channel, followRepo: followRepo, userRepo: userRepo}
}

func (c *FollowConsumer) Declare(exchange, exchangeType, queue, routingKey string) error {
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

func (c *FollowConsumer) ListenFollowConsumer(queue string, handler func(amqp.Delivery)) error {
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

func (c *FollowConsumer) FollowHandler(msg amqp.Delivery) {
	if c == nil || c.channel == nil {
		log.Println("consumer channel is nil")
		_ = msg.Nack(false, false)
		return
	}

	var followEvent event.FollowEvent
	if err := json.Unmarshal(msg.Body, &followEvent); err != nil {
		log.Println(err)
		_ = msg.Nack(false, false)
		return
	}

	follow := &model.Follow{Following: followEvent.Following, Follower: followEvent.Follower}
	tx := c.followRepo.DB().Begin()
	var err error
	switch followEvent.EventType {
	case event.Follow:
		err = tx.Create(follow).Error
		if err == nil {
			err = tx.Model(&model.User{}).Where("id = ?", followEvent.Follower).Update("follow_count", gorm.Expr("follow_count + 1")).Error
		}
		if err == nil {
			err = tx.Model(&model.User{}).Where("id = ?", followEvent.Following).Update("follower_count", gorm.Expr("follower_count + 1")).Error
		}
	case event.Unfollow:
		err = tx.Where("follower = ? and following = ?", followEvent.Follower, followEvent.Following).Delete(&model.Follow{}).Error
		if err == nil {
			err = tx.Model(&model.User{}).Where("id = ?", followEvent.Follower).Update("follow_count", gorm.Expr("CASE WHEN follow_count > 0 THEN follow_count - 1 ELSE 0 END")).Error
		}
		if err == nil {
			err = tx.Model(&model.User{}).Where("id = ?", followEvent.Following).Update("follower_count", gorm.Expr("CASE WHEN follower_count > 0 THEN follower_count - 1 ELSE 0 END")).Error
		}
	default:
		_ = tx.Rollback()
		log.Printf("unsupported follow event type: %s", followEvent.EventType)
		_ = msg.Nack(false, false)
		return
	}

	if err != nil {
		_ = tx.Rollback()
		log.Println(err)
		_ = msg.Nack(false, true)
		return
	}
	if err = tx.Commit().Error; err != nil {
		log.Println(err)
		_ = msg.Nack(false, true)
		return
	}
	_ = msg.Ack(false)
}


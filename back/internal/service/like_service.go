package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"simple_tiktok/internal/dto/res"
	"simple_tiktok/internal/mq/event"
	"simple_tiktok/internal/pkg/constants"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
)

type LikeService struct {
	redisClient *redis.Client
	likeMQ      *amqp.Channel
}

var switchLikeScript = redis.NewScript(`
if redis.call("SISMEMBER", KEYS[1], ARGV[1]) == 1 then
    redis.call("SREM", KEYS[1], ARGV[1])
    return 0
end
redis.call("SADD", KEYS[1], ARGV[1])
return 1
`)

func NewLikeService(redisClient *redis.Client, ch *amqp.Channel) *LikeService {
	likeService := &LikeService{
		redisClient: redisClient,
		likeMQ:      ch,
	}
	return likeService
}

func (s *LikeService) LikeVideo(targetId uint64, userId uint64) (res.LikeVideoRes, error) {
	// 先原子切换 redis 点赞状态，再投递 mq 更新 mysql
	key := fmt.Sprintf(constants.LikeVideo, targetId)
	liked, err := s.switchLike(context.Background(), key, userId)
	if err != nil {
		return res.LikeVideoRes{}, err
	}

	eventType := event.Dislike
	rollback := func() error {
		return s.redisClient.SAdd(context.Background(), key, userId).Err()
	}
	if liked {
		eventType = event.Like
		rollback = func() error {
			return s.redisClient.SRem(context.Background(), key, userId).Err()
		}
	}

	msg, err := s.getLikeVideoEventMsg(targetId, eventType)
	if err != nil {
		return res.LikeVideoRes{}, err
	}

	log.Printf("like switch video_id=%d user_id=%d liked=%t", targetId, userId, liked)
	err = s.likeMQ.Publish(event.LikeVideoExchange, event.LikeVideoRoutingKey, false, false, msg)
	if err != nil {
		if rollbackErr := rollback(); rollbackErr != nil {
			log.Printf("like publish failed and rollback failed video_id=%d user_id=%d err=%v rollback_err=%v", targetId, userId, err, rollbackErr)
		}
		return res.LikeVideoRes{}, err
	}

	return res.LikeVideoRes{VideoId: targetId, IsLiked: liked}, nil
}

func (s *LikeService) switchLike(ctx context.Context, key string, userId uint64) (bool, error) {
	result, err := switchLikeScript.Run(ctx, s.redisClient, []string{key}, userId).Int()
	if err != nil {
		return false, err
	}
	return result == 1, nil
}
func (s *LikeService) getLikeVideoEventMsg(videoId uint64, eventType string) (amqp.Publishing, error) {
	e := event.LikeVideoEvent{VideoId: videoId, EventType: eventType}
	data, err := json.Marshal(e)
	if err != nil {
		return amqp.Publishing{}, err
	}
	msg := amqp.Publishing{
		ContentType: "application/json",
		Body:        data,
	}
	return msg, nil
}

func (s *LikeService) LikeComment(commentId uint64, userId uint64) (res.LikeCommentRes, error) {
	key := fmt.Sprintf(constants.LikeComment, commentId)
	liked, err := s.switchLike(context.Background(), key, userId)
	if err != nil {
		return res.LikeCommentRes{}, err
	}

	eventType := event.Dislike
	rollback := func() error {
		return s.redisClient.SAdd(context.Background(), key, userId).Err()
	}
	if liked {
		eventType = event.Like
		rollback = func() error {
			return s.redisClient.SRem(context.Background(), key, userId).Err()
		}
	}

	msg, err := s.getLikeCommentEventMsg(commentId, eventType)
	if err != nil {
		return res.LikeCommentRes{}, err
	}

	log.Printf("comment like switch comment_id=%d user_id=%d liked=%t", commentId, userId, liked)
	err = s.likeMQ.Publish(event.LikeCommentExchange, event.LikeCommentRoutingKey, false, false, msg)
	if err != nil {
		if rollbackErr := rollback(); rollbackErr != nil {
			log.Printf("comment like publish failed and rollback failed comment_id=%d user_id=%d err=%v rollback_err=%v", commentId, userId, err, rollbackErr)
		}
		return res.LikeCommentRes{}, err
	}

	return res.LikeCommentRes{CommentId: commentId, IsLiked: liked}, nil
}

func (s *LikeService) getLikeCommentEventMsg(commentId uint64, eventType string) (amqp.Publishing, error) {
	e := event.LikeCommentEvent{CommentId: commentId, EventType: eventType}
	data, err := json.Marshal(e)
	if err != nil {
		return amqp.Publishing{}, err
	}
	msg := amqp.Publishing{
		ContentType: "application/json",
		Body:        data,
	}
	return msg, nil
}

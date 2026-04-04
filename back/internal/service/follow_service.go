package service

import (
	"context"
	"fmt"
	"simple_tiktok/internal/model"
	"simple_tiktok/internal/pkg/constants"
	"simple_tiktok/internal/repository/mysql"

	"github.com/redis/go-redis/v9"
)

type FollowService struct {
	followRepo  *mysql.FollowRepo
	redisClient *redis.Client
}

func NewFollowService(followRepo *mysql.FollowRepo, redisClient *redis.Client) *FollowService {
	return &FollowService{
		followRepo:  followRepo,
		redisClient: redisClient,
	}
}

func (s *FollowService) Follow(id uint64, follower uint64) error {
	key := fmt.Sprintf(constants.FollowKey, follower)
	isFollowed, err := s.redisClient.SIsMember(context.Background(), key, id).Result()
	if err != nil {
		return err
	}
	follow := &model.Follow{Following: id, Follower: follower}
	if !isFollowed {
		err := s.redisClient.SAdd(context.Background(), key, id).Err()
		if err != nil {
			return err
		}
		return s.followRepo.Follow(follow)
	}
	err = s.redisClient.SRem(context.Background(), key, id).Err()
	if err != nil {
		return err
	}
	return s.followRepo.DeleteFollow(follow)
}

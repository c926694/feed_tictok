package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math"
	"simple_tiktok/internal/dto/req"
	"simple_tiktok/internal/dto/res"
	"simple_tiktok/internal/model"
	"simple_tiktok/internal/mq/event"
	"simple_tiktok/internal/pkg/constants"
	"simple_tiktok/internal/pkg/upload"
	"simple_tiktok/internal/pkg/util"
	mysql2 "simple_tiktok/internal/repository/mysql"
	"strconv"
	"strings"
	"time"

	//"github.com/redis/go-redis/v9"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
)

type VideoService struct {
	videoRepo   *mysql2.VideoRepo
	userRepo    *mysql2.UserRepo
	redisClient *redis.Client
	videoMQ     *amqp.Channel
	commentRepo *mysql2.CommentRepo
}

func NewVideoService(videoRepo *mysql2.VideoRepo, userRepo *mysql2.UserRepo, redisClient *redis.Client,
	videoMQ *amqp.Channel, commentRepo *mysql2.CommentRepo) *VideoService {
	return &VideoService{
		videoRepo:   videoRepo,
		userRepo:    userRepo,
		redisClient: redisClient,
		videoMQ:     videoMQ,
		commentRepo: commentRepo,
	}
}

func (s *VideoService) CreateVideo(req req.UploadVideoReq, userId uint64, nickName string) (res.VideoRes, error) {
	cover := req.Cover
	play := req.Play
	title := req.Title
	description := req.Description

	authorName := strings.TrimSpace(nickName)
	if s.userRepo != nil {
		user, userErr := s.userRepo.GetUserByID(userId)
		if userErr == nil && user != nil {
			currentNickName := strings.TrimSpace(user.NickName)
			if currentNickName != "" {
				authorName = currentNickName
			}
		}
	}

	coverPath, err := upload.UploadFile(cover, upload.Cover)
	if err != nil {
		return res.VideoRes{}, err
	}

	playPath, err := upload.UploadFile(play, upload.Video)
	if err != nil {
		err2 := upload.Delete(upload.Cover, coverPath)
		if err2 != nil {
			return res.VideoRes{}, err2
		}
		return res.VideoRes{}, err
	}

	video := model.Video{
		Title:       title,
		Description: description,
		AuthorID:    userId,
		PlayURL:     playPath,
		CoverURL:    coverPath,
		AuthorName:  authorName,
	}
	log.Println("nickName:", authorName)

	err = s.videoRepo.CreateVideo(&video)
	if err != nil {
		err2 := upload.Delete(upload.Video, playPath)
		if err2 != nil {
			return res.VideoRes{}, err2
		}
		return res.VideoRes{}, err
	}

	err = s.AddZSet(constants.FeedVideoKey, float64(video.CreatedAt.UnixMicro()), video.ID)
	if err != nil {
		return res.VideoRes{}, err
	}
	hotScore := s.HotScore(video.LikeCount, video.CommentCount, video.ID)
	err = s.AddZSet(constants.HotFeedVideoKey, hotScore, video.ID)
	return res.VideoRes{
		Id:  video.ID,
		Url: util.EnsureHTTPPath(video.PlayURL),
	}, nil
}

func (s *VideoService) AddZSet(key string, score float64, member uint64) error {
	value := redis.Z{
		Score:  score,
		Member: member,
	}
	_, err := s.redisClient.ZAdd(context.Background(), key, value).Result()
	if err != nil {
		return err
	}
	return nil
}

func (s *VideoService) RemZSet(key string, member uint64) error {
	_, err := s.redisClient.ZRem(context.Background(), key, member).Result()
	if err != nil {
		return err
	}
	return nil
}

func (s *VideoService) GetFeedVideos(limit uint64, lastScore float64, key string, userId uint64) ([]res.VideoInfoRes, float64, error) {
	ids, err := s.GetFeedVideoIds(limit, lastScore, key)
	if err != nil {
		return nil, 0.0, err
	}
	if len(ids) == 0 {
		return []res.VideoInfoRes{}, 0.0, nil
	}
	videoInfoList, err := s.getVideoInfoByIds(ids)
	if err != nil {
		return nil, 0.0, err
	}
	videoInfoList, err = s.getFinalVideoInfoResList(videoInfoList, ids)
	if err != nil {
		return nil, 0.0, err
	}
	s.fillVideoAuthorAvatar(videoInfoList)
	if err = s.fillVideoLikeStatus(videoInfoList, userId); err != nil {
		return nil, 0.0, err
	}
	if err = s.fillVideoFollowStatus(videoInfoList, userId); err != nil {
		return nil, 0.0, err
	}
	nextScore := float64(videoInfoList[len(videoInfoList)-1].CreatedAt.UnixMicro())
	return videoInfoList, nextScore, nil
}

func (s *VideoService) getFinalVideoInfoResList(videoInfoList []res.VideoInfoRes, ids []uint64) ([]res.VideoInfoRes, error) {
	videoMap := make(map[uint64]res.VideoInfoRes)
	for _, videoInfo := range videoInfoList {
		videoMap[videoInfo.Id] = videoInfo
	}
	result := make([]res.VideoInfoRes, len(ids))
	for i, id := range ids {
		result[i] = videoMap[id]
	}
	return result, nil
}

func (s *VideoService) GetFeedHotVideosAndLastCore(limit uint64, lastScore float64, key string, userId uint64) ([]res.VideoInfoRes, float64, error) {
	ids, err := s.GetFeedVideoIds(limit, lastScore, key)
	if err != nil {
		return nil, 0, err
	}
	if len(ids) == 0 {
		return []res.VideoInfoRes{}, 0, nil
	}
	videoInfoList, err := s.getVideoInfoByIds(ids)
	if err != nil {
		return nil, 0, err
	}
	videoInfoList, err = s.getFinalVideoInfoResList(videoInfoList, ids)
	if err != nil {
		return nil, 0, err
	}
	s.fillVideoAuthorAvatar(videoInfoList)
	if err = s.fillVideoLikeStatus(videoInfoList, userId); err != nil {
		return nil, 0, err
	}
	if err = s.fillVideoFollowStatus(videoInfoList, userId); err != nil {
		return nil, 0, err
	}
	member := &redis.ZRangeBy{
		Min:    "-inf",
		Max:    "+inf",
		Offset: 0,
		Count:  int64(limit),
	}
	if lastScore > 0 {
		member.Max = "(" + strconv.FormatFloat(lastScore, 'f', -1, 64)
	}
	ZList, err := s.redisClient.ZRevRangeByScoreWithScores(context.Background(), key, member).Result()
	if err != nil {
		return nil, 0, err
	}
	nextScore := 0.0
	if len(ZList) > 0 {
		// 从 Redis 结果中获取下一页游标（score）
		nextScore = ZList[len(ZList)-1].Score
	}
	return videoInfoList, nextScore, nil
}

func (s *VideoService) GetFollowFeedVideos(limit uint64, lastScore float64, userId uint64) ([]res.VideoInfoRes, float64, error) {
	followingIDs, err := s.getFollowingUserIDs(userId)
	if err != nil {
		return nil, 0, err
	}
	if len(followingIDs) == 0 {
		return []res.VideoInfoRes{}, 0, nil
	}

	var lastCreatedAt *time.Time
	if lastScore > 0 && lastScore < float64(math.MaxInt64) {
		t := time.UnixMicro(int64(lastScore))
		lastCreatedAt = &t
	}

	videoList, err := s.videoRepo.GetFollowFeedVideosByAuthors(followingIDs, limit, lastCreatedAt)
	if err != nil {
		return nil, 0, err
	}
	if len(videoList) == 0 {
		return []res.VideoInfoRes{}, 0, nil
	}

	videoInfoList := make([]res.VideoInfoRes, len(videoList))
	for i, v := range videoList {
		videoInfoList[i] = res.VideoInfoRes{
			Id:           v.ID,
			AuthorID:     v.AuthorID,
			AuthorName:   v.AuthorName,
			Title:        v.Title,
			Description:  v.Description,
			CoverURL:     util.EnsureHTTPPath(v.CoverURL),
			PlayURL:      util.EnsureHTTPPath(v.PlayURL),
			CreatedAt:    v.CreatedAt,
			LikeCount:    v.LikeCount,
			CommentCount: v.CommentCount,
		}
	}
	s.fillVideoAuthorAvatar(videoInfoList)
	if err = s.fillVideoLikeStatus(videoInfoList, userId); err != nil {
		return nil, 0, err
	}
	if err = s.fillVideoFollowStatus(videoInfoList, userId); err != nil {
		return nil, 0, err
	}
	nextScore := float64(videoInfoList[len(videoInfoList)-1].CreatedAt.UnixMicro())
	return videoInfoList, nextScore, nil
}

func (s *VideoService) GetMyVideos(userId uint64, limit uint64) ([]res.VideoInfoRes, error) {
	videoList, err := s.videoRepo.ListByAuthorID(userId, limit)
	if err != nil {
		return nil, err
	}
	videoInfoList := make([]res.VideoInfoRes, len(videoList))
	for i, v := range videoList {
		videoInfoList[i] = res.VideoInfoRes{
			Id:           v.ID,
			AuthorID:     v.AuthorID,
			AuthorName:   v.AuthorName,
			Title:        v.Title,
			Description:  v.Description,
			CoverURL:     util.EnsureHTTPPath(v.CoverURL),
			PlayURL:      util.EnsureHTTPPath(v.PlayURL),
			CreatedAt:    v.CreatedAt,
			LikeCount:    v.LikeCount,
			CommentCount: v.CommentCount,
		}
	}
	s.fillVideoAuthorAvatar(videoInfoList)
	if err = s.fillVideoLikeStatus(videoInfoList, userId); err != nil {
		return nil, err
	}
	if err = s.fillVideoFollowStatus(videoInfoList, userId); err != nil {
		return nil, err
	}
	return videoInfoList, nil
}

func (s *VideoService) getFollowingUserIDs(userId uint64) ([]uint64, error) {
	key := fmt.Sprintf(constants.FollowKey, userId)
	members, err := s.redisClient.SMembers(context.Background(), key).Result()
	if err != nil {
		return nil, err
	}
	result := make([]uint64, 0, len(members))
	for _, member := range members {
		id, parseErr := strconv.ParseUint(member, 10, 64)
		if parseErr != nil {
			return nil, parseErr
		}
		result = append(result, id)
	}
	return result, nil
}

func (s *VideoService) fillVideoAuthorAvatar(videoInfoList []res.VideoInfoRes) {
	if s.userRepo == nil {
		return
	}
	type authorProfile struct {
		name   string
		avatar string
	}
	cache := make(map[uint64]authorProfile)
	for i := range videoInfoList {
		authorID := videoInfoList[i].AuthorID
		if authorID == 0 {
			continue
		}
		if profile, ok := cache[authorID]; ok {
			if profile.name != "" {
				videoInfoList[i].AuthorName = profile.name
			}
			videoInfoList[i].AuthorAvatar = profile.avatar
			continue
		}
		user, err := s.userRepo.GetUserByID(authorID)
		if err != nil || user == nil {
			continue
		}
		profile := authorProfile{
			name:   user.NickName,
			avatar: util.EnsureHTTPPath(user.AvatarURL),
		}
		cache[authorID] = profile
		if profile.name != "" {
			videoInfoList[i].AuthorName = profile.name
		}
		videoInfoList[i].AuthorAvatar = profile.avatar
	}
}

func (s *VideoService) fillVideoLikeStatus(videoInfoList []res.VideoInfoRes, userId uint64) error {
	for i := range videoInfoList {
		likeKey := fmt.Sprintf(constants.LikeVideo, videoInfoList[i].Id)
		isLiked, err := s.redisClient.SIsMember(context.Background(), likeKey, userId).Result()
		if err != nil {
			return err
		}
		videoInfoList[i].IsLiked = isLiked
	}
	return nil
}

func (s *VideoService) fillVideoFollowStatus(videoInfoList []res.VideoInfoRes, userId uint64) error {
	followKey := fmt.Sprintf(constants.FollowKey, userId)
	for i := range videoInfoList {
		if videoInfoList[i].AuthorID == 0 || videoInfoList[i].AuthorID == userId {
			videoInfoList[i].IsFollow = false
			continue
		}
		isFollow, err := s.redisClient.SIsMember(context.Background(), followKey, videoInfoList[i].AuthorID).Result()
		if err != nil {
			return err
		}
		videoInfoList[i].IsFollow = isFollow
	}
	return nil
}

func (s *VideoService) HotScore(likeCount, commentCount int64, videoId uint64) float64 {
	return float64(likeCount*2+commentCount) + float64(videoId)/1e10
}

func (s *VideoService) GetFeedVideoIds(limit uint64, lastScore float64, key string) ([]uint64, error) {
	member := &redis.ZRangeBy{
		Min:    "-inf",
		Max:    "+inf",
		Offset: 0,
		Count:  int64(limit),
	}
	if lastScore > 0 {
		member.Max = "(" + strconv.FormatFloat(lastScore, 'f', -1, 64)
	}
	//log.Println("max===", member.Max)
	idsStr, err := s.redisClient.ZRevRangeByScore(context.Background(), key, member).Result()
	if err != nil {
		return nil, err
	}
	videoIds := make([]uint64, len(idsStr))
	for i, v := range idsStr {
		id, err := strconv.ParseUint(v, 10, 64)
		if err != nil {
			return nil, err
		}
		videoIds[i] = id
	}
	return videoIds, nil
}

func (s *VideoService) GetVideoInfo(id uint64, userId uint64) (res.VideoInfoRes, error) {
	video, err := s.videoRepo.GetVideoById(id)
	if err != nil {
		return res.VideoInfoRes{}, err
	}
	videoInfo := res.VideoInfoRes{
		Id:           video.ID,
		AuthorID:     video.AuthorID,
		AuthorName:   video.AuthorName,
		Title:        video.Title,
		Description:  video.Description,
		CoverURL:     util.EnsureHTTPPath(video.CoverURL),
		PlayURL:      util.EnsureHTTPPath(video.PlayURL),
		CommentCount: video.CommentCount,
		LikeCount:    video.LikeCount,
		CreatedAt:    video.CreatedAt,
	}
	videoInfoList := []res.VideoInfoRes{videoInfo}
	s.fillVideoAuthorAvatar(videoInfoList)
	if err = s.fillVideoLikeStatus(videoInfoList, userId); err != nil {
		return res.VideoInfoRes{}, err
	}
	if err = s.fillVideoFollowStatus(videoInfoList, userId); err != nil {
		return res.VideoInfoRes{}, err
	}
	return videoInfoList[0], nil
}

func (s *VideoService) getVideoInfoByIds(ids []uint64) ([]res.VideoInfoRes, error) {
	videoList, err := s.videoRepo.GetFeedVideos(ids)
	if err != nil {
		return nil, err
	}
	videoInfoList := make([]res.VideoInfoRes, len(videoList))
	for i, v := range videoList {
		videoInfoList[i] = res.VideoInfoRes{
			Id:           v.ID,
			AuthorID:     v.AuthorID,
			AuthorName:   v.AuthorName,
			Title:        v.Title,
			Description:  v.Description,
			CoverURL:     util.EnsureHTTPPath(v.CoverURL),
			PlayURL:      util.EnsureHTTPPath(v.PlayURL),
			CreatedAt:    v.CreatedAt,
			LikeCount:    v.LikeCount,
			CommentCount: v.CommentCount,
		}
	}
	return videoInfoList, nil
}

func (s *VideoService) DeleteVideo(id uint64, userId uint64) error {
	tx := s.videoRepo.DB().Begin()
	video, err := s.videoRepo.GetVideoById(id)
	if err != nil {
		return err
	}
	if video.AuthorID != userId {
		_ = tx.Rollback()
		return errors.New("no permission to delete this video")
	}

	commentRepo := s.commentRepo.WithTx(tx)
	commentList, err := commentRepo.ListByVideoId(id)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	videoRepo := s.videoRepo.WithTx(tx)
	if err := videoRepo.DeleteVideoById(id); err != nil {
		_ = tx.Rollback()
		return err
	}
	if err := commentRepo.DeleteByVideoId(id); err != nil {
		_ = tx.Rollback()
		return err
	}
	if err := tx.Commit().Error; err != nil {
		return err
	}

	if err := s.RemZSet(constants.FeedVideoKey, id); err != nil {
		return err
	}
	if err := s.RemZSet(constants.HotFeedVideoKey, id); err != nil {
		return err
	}

	likeKey := fmt.Sprintf(constants.LikeVideo, id)
	if err := s.redisClient.Del(context.Background(), likeKey).Err(); err != nil {
		return err
	}
	if len(commentList) > 0 {
		commentLikeKeys := make([]string, 0, len(commentList))
		for _, comment := range commentList {
			commentLikeKeys = append(commentLikeKeys, fmt.Sprintf(constants.LikeComment, comment.ID))
		}
		if err := s.redisClient.Del(context.Background(), commentLikeKeys...).Err(); err != nil {
			return err
		}
	}

	videoProducer := s.videoMQ
	msg, err := s.getDeleteVideoEvent(video)
	if err != nil {
		return err
	}
	err = videoProducer.Publish(event.DeleteVideoExchange,
		event.DeleteVideoRoutingKey, false, false, msg)
	if err != nil {
		return err
	}
	return nil
}

func (s *VideoService) getDeleteVideoEvent(video model.Video) (amqp.Publishing, error) {
	deleteVideoEvent := event.DeleteVideoEvent{
		PlayURL:  video.PlayURL,
		CoverURL: video.CoverURL,
	}
	data, err := json.Marshal(deleteVideoEvent)
	if err != nil {
		return amqp.Publishing{}, err
	}
	msg := amqp.Publishing{
		ContentType: "application/json",
		Body:        data,
	}
	return msg, nil
}

func (s *VideoService) UpdateZSet(key string, score float64, id uint64) error {
	member := redis.Z{
		Score:  score,
		Member: id,
	}
	return s.redisClient.ZAdd(context.Background(), key, member).Err()
}

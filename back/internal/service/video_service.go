package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"simple_tiktok/internal/dto/req"
	"simple_tiktok/internal/dto/res"
	"simple_tiktok/internal/model"
	"simple_tiktok/internal/mq/event"
	"simple_tiktok/internal/pkg/constants"
	"simple_tiktok/internal/pkg/upload"
	mysql2 "simple_tiktok/internal/repository/mysql"
	"strconv"

	//"github.com/redis/go-redis/v9"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
)

type VideoService struct {
	videoRepo   *mysql2.VideoRepo
	redisClient *redis.Client
	videoMQ     *amqp.Channel
	commentRepo *mysql2.CommentRepo
}

func NewVideoService(videoRepo *mysql2.VideoRepo, redisClient *redis.Client,
	videoMQ *amqp.Channel, commentRepo *mysql2.CommentRepo) *VideoService {
	return &VideoService{
		videoRepo:   videoRepo,
		redisClient: redisClient,
		videoMQ:     videoMQ,
		commentRepo: commentRepo,
	}
}

// 上传视频
func (s *VideoService) CreateVideo(req req.UploadVideoReq, userId uint64, nickName string) (res.VideoRes, error) {
	cover := req.Cover
	play := req.Play
	title := req.Title
	//上传封面
	coverPath, err := upload.UploadFile(cover, upload.Cover)
	if err != nil {
		return res.VideoRes{}, err
	}
	//上传视频
	playPath, err := upload.UploadFile(play, upload.Video)
	if err != nil {
		//删除封面
		err2 := upload.Delete(upload.Cover, coverPath)
		if err2 != nil {
			return res.VideoRes{}, err2
		}
		return res.VideoRes{}, err
	}
	//数据库保存
	video := model.Video{
		Title:      title,
		AuthorID:   userId,
		PlayURL:    playPath,
		CoverURL:   coverPath,
		AuthorName: nickName,
	}
	log.Println("nickName:", nickName)
	err = s.videoRepo.CreateVideo(&video)
	if err != nil {
		//删除视频
		err2 := upload.Delete(upload.Video, playPath)
		if err2 != nil {
			return res.VideoRes{}, err2
		}
		return res.VideoRes{}, err
	}
	//更新到redis

	// 1.普通feed
	err = s.AddZSet(constants.FeedVideoKey, float64(video.CreatedAt.UnixMicro()), video.ID)
	if err != nil {
		return res.VideoRes{}, err
	}
	// 2.hot feed
	hotScore := s.HotScore(video.LikeCount, video.CommentCount, video.ID)
	err = s.AddZSet(constants.HotFeedVideoKey, hotScore, video.ID)
	return res.VideoRes{
		Id:  video.ID,
		Url: video.PlayURL,
	}, nil
}

// 添加zSet工具函数
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

// 添加zRem工具函数
func (s *VideoService) RemZSet(key string, member uint64) error {
	_, err := s.redisClient.ZRem(context.Background(), key, member).Result()
	if err != nil {
		return err
	}
	return nil
}

// 获取普通feed视频列表
func (s *VideoService) GetFeedVideos(limit uint64, lastScore float64, key string) ([]res.VideoInfoRes, float64, error) {
	//视频ids

	ids, err := s.GetFeedVideoIds(limit, lastScore, key)
	if err != nil {
		return nil, 0.0, err
	}
	if len(ids) == 0 {
		return []res.VideoInfoRes{}, 0.0, nil
	}
	//无序视频列表
	videoInfoList, err := s.getVideoInfoByIds(ids)
	if err != nil {
		return nil, 0.0, err
	}
	//转成有序视频列表并获取lastScore
	videoInfoList, err = s.getFinalVideoInfoResList(videoInfoList, ids)
	nextScore := float64(videoInfoList[len(videoInfoList)-1].CreatedAt.UnixMicro())
	return videoInfoList, nextScore, nil
}

// 获取有序视频列表
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

// 获取热榜feed视频列表
func (s *VideoService) GetFeedHotVideosAndLastCore(limit uint64, lastScore float64, key string) ([]res.VideoInfoRes, float64, error) {
	ids, err := s.GetFeedVideoIds(limit, lastScore, key)
	if err != nil {
		return nil, 0, err
	}
	videoInfoList, err := s.getVideoInfoByIds(ids)
	if err != nil {
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
		//从redis获取nextScore
		nextScore = ZList[len(ZList)-1].Score
	}
	return videoInfoList, nextScore, nil
}

// 热度计算公式
func (s *VideoService) HotScore(likeCount, commentCount int64, videoId uint64) float64 {
	return float64(likeCount*2+commentCount) + float64(videoId)/1e10
}

// 获取feed ids
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

func (s *VideoService) GetVideoInfo(id uint64) (res.VideoInfoRes, error) {
	video, err := s.videoRepo.GetVideoById(id)
	if err != nil {
		return res.VideoInfoRes{}, err
	}
	return res.VideoInfoRes{
		Id:           video.ID,
		AuthorName:   video.AuthorName,
		CoverURL:     video.CoverURL,
		PlayURL:      video.PlayURL,
		CommentCount: video.CommentCount,
		LikeCount:    video.LikeCount,
		CreatedAt:    video.CreatedAt,
	}, nil
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
			AuthorName:   v.AuthorName,
			CoverURL:     v.CoverURL,
			PlayURL:      v.PlayURL,
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
	//删除video
	videoRepo := s.videoRepo.WithTx(tx)
	if err := videoRepo.DeleteVideoById(id); err != nil {
		tx.Rollback()
		return err
	}
	if err := s.commentRepo.DeleteByVideoId(id); err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	//删除redis

	// 1.feed
	if err := s.RemZSet(constants.FeedVideoKey, id); err != nil {
		return err
	}
	// 2.hot feed
	if err := s.RemZSet(constants.HotFeedVideoKey, id); err != nil {
		return err
	}
	// 3.like
	likeKey := fmt.Sprintf(constants.LikeVideo, id)
	if err := s.redisClient.SRem(context.Background(), likeKey, userId).Err(); err != nil {
		return err
	}
	//异步删除文件
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

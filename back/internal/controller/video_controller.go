package controller

import (
	"math"
	"net/http"
	"simple_tiktok/internal/dto/req"
	"simple_tiktok/internal/dto/res"
	"simple_tiktok/internal/middleware"
	"simple_tiktok/internal/pkg/constants"
	"simple_tiktok/internal/pkg/response"
	"simple_tiktok/internal/pkg/type_convert"
	"simple_tiktok/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type VideoController struct {
	service *service.VideoService
}

func NewVideoController(videoService *service.VideoService) *VideoController {
	return &VideoController{
		service: videoService,
	}
}

func (ctl *VideoController) Publish(c *gin.Context) {
	response.Fail(c, http.StatusNotImplemented, "not implemented")
}

func (ctl *VideoController) Delete(c *gin.Context) {
	response.Fail(c, http.StatusNotImplemented, "not implemented")
}

func (ctl *VideoController) CreateVideo(c *gin.Context) {
	var createVideoReq req.UploadVideoReq
	play, _ := c.FormFile("play")
	cover, _ := c.FormFile("cover")
	title := c.PostForm("title")
	description := c.PostForm("description")
	createVideoReq = req.UploadVideoReq{
		Title:       title,
		Description: description,
		Play:        play,
		Cover:       cover,
	}
	videoRes, err := ctl.service.CreateVideo(
		createVideoReq, c.MustGet(middleware.UserCtx).(uint64), c.MustGet(middleware.UserNickName).(string))
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.OK(c, videoRes)
}

func (ctl *VideoController) GetFeedVideos(c *gin.Context) {
	lastScore, err := strconv.ParseFloat(c.DefaultQuery("last_score", strconv.FormatFloat(math.MaxFloat64, 'f', -1, 64)), 64)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, err.Error())
	}
	limit, err := strconv.ParseUint(c.DefaultQuery("limit", "3"), 10, 64)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, err.Error())
	}
	userId, err := type_convert.AnyToUint64(c.MustGet(middleware.UserCtx))
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}
	//log.Printf("last_id=====%d  create_at====%v", feedVideoReq.LastId, feedVideoReq.CreatedAt)
	videoInfoResList, nextScore, err := ctl.service.GetFeedVideos(limit, lastScore, constants.FeedVideoKey, userId)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.OK(c, &res.FeedVideoRes{
		FeedVideoList: videoInfoResList,
		LastScore:     nextScore,
	})
}

func (ctl *VideoController) GetVideoInfo(c *gin.Context) {
	rawID := c.Param("id")
	if rawID == "me" {
		limit, err := strconv.ParseUint(c.DefaultQuery("limit", "60"), 10, 64)
		if err != nil {
			response.Fail(c, http.StatusBadRequest, "invalid limit")
			return
		}
		userId, err := type_convert.AnyToUint64(c.MustGet(middleware.UserCtx))
		if err != nil {
			response.Fail(c, http.StatusUnauthorized, "invalid user")
			return
		}
		videoInfoResList, err := ctl.service.GetMyVideos(userId, limit)
		if err != nil {
			response.Fail(c, http.StatusInternalServerError, err.Error())
			return
		}
		response.OK(c, videoInfoResList)
		return
	}

	videoId, err := strconv.ParseUint(rawID, 10, 64)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid video id")
		return
	}
	userId, err := type_convert.AnyToUint64(c.MustGet(middleware.UserCtx))
	if err != nil {
		response.Fail(c, http.StatusUnauthorized, "invalid user")
		return
	}
	videoInfoRes, err := ctl.service.GetVideoInfo(videoId, userId)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.OK(c, videoInfoRes)
}

func (ctl *VideoController) GetFeedHotVideos(c *gin.Context) {
	interval, err := strconv.Atoi(c.DefaultQuery("interval", "60"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid interval")
		return
	}
	limit, err := strconv.ParseUint(c.DefaultQuery("limit", "3"), 10, 64)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid limit")
		return
	}
	offset, err := strconv.ParseUint(c.DefaultQuery("offset", "0"), 10, 64)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid offset")
		return
	}
	userId, err := type_convert.AnyToUint64(c.MustGet(middleware.UserCtx))
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}
	videoInfoResList, nextOffset, hasMore, err := ctl.service.GetFeedHotVideos(limit, offset, interval, userId)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.OK(c, &res.HotFeedVideoRes{
		FeedVideoList: videoInfoResList,
		NextOffset:    nextOffset,
		HasMore:       hasMore,
		Interval:      interval,
	})
}

func (ctl *VideoController) GetFollowFeedVideos(c *gin.Context) {
	lastScore, err := strconv.ParseFloat(c.DefaultQuery("last_score", strconv.FormatFloat(math.MaxFloat64, 'f', -1, 64)), 64)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}
	limit, err := strconv.ParseUint(c.DefaultQuery("limit", "3"), 10, 64)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}
	userId, err := type_convert.AnyToUint64(c.MustGet(middleware.UserCtx))
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}
	videoInfoResList, nextScore, err := ctl.service.GetFollowFeedVideos(limit, lastScore, userId)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.OK(c, &res.FeedVideoRes{
		FeedVideoList: videoInfoResList,
		LastScore:     nextScore,
	})
}

func (ctl *VideoController) GetMyVideos(c *gin.Context) {
	limit, err := strconv.ParseUint(c.DefaultQuery("limit", "60"), 10, 64)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}
	userId, err := type_convert.AnyToUint64(c.MustGet(middleware.UserCtx))
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}
	videoInfoResList, err := ctl.service.GetMyVideos(userId, limit)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.OK(c, videoInfoResList)
}

func (ctl *VideoController) DeleteVideos(c *gin.Context) {
	videoId, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, err.Error())
		return
	}
	userId, err := type_convert.AnyToUint64(c.MustGet(middleware.UserCtx))
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, err.Error())
	}
	err = ctl.service.DeleteVideo(videoId, userId)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.OK(c, nil)
}

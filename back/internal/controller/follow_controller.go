package controller

import (
	"net/http"
	"simple_tiktok/internal/middleware"
	"simple_tiktok/internal/pkg/response"
	"simple_tiktok/internal/pkg/type_convert"
	"simple_tiktok/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type FollowController struct {
	service *service.FollowService
}

func NewFollowController(followService *service.FollowService) *FollowController {
	return &FollowController{service: followService}
}

func (ctl *FollowController) Follow(c *gin.Context) {
	targetUserID, err := strconv.ParseUint(c.Param("follower"), 10, 64)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, err.Error())
		return
	}
	currentUserID, err := type_convert.AnyToUint64(c.MustGet(middleware.UserCtx))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, err.Error())
		return
	}
	result, err := ctl.service.Follow(targetUserID, currentUserID)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(c, result)
}

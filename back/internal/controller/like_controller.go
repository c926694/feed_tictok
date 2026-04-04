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

type LikeController struct {
	service *service.LikeService
}

func NewLikeController(likeService *service.LikeService) *LikeController {
	return &LikeController{service: likeService}
}

func (ctl *LikeController) Action(c *gin.Context) {
	response.Fail(c, http.StatusNotImplemented, "not implemented")
}

func (ctl *LikeController) LikeVideo(c *gin.Context) {
	targetId, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}
	userId, err := type_convert.AnyToUint64(c.MustGet(middleware.UserCtx))
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}
	result, err := ctl.service.LikeVideo(targetId, userId)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.OK(c, result)
}

func (ctl *LikeController) LikeComment(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}
	userId, err := type_convert.AnyToUint64(c.MustGet(middleware.UserCtx))
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}
	result, err := ctl.service.LikeComment(id, userId)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.OK(c, result)
}

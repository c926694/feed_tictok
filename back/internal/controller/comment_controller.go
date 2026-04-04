package controller

import (
	"net/http"
	"simple_tiktok/internal/dto/req"
	"simple_tiktok/internal/middleware"
	"simple_tiktok/internal/pkg/response"
	"simple_tiktok/internal/pkg/type_convert"
	"simple_tiktok/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CommentController struct {
	service *service.CommentService
}

func NewCommentController(commentService *service.CommentService) *CommentController {
	return &CommentController{service: commentService}
}

func (ctl *CommentController) Create(c *gin.Context) {
	var commentReq req.CommentReq
	if err := c.ShouldBind(&commentReq); err != nil {
		response.Fail(c, http.StatusBadRequest, err.Error())
	}
	userId, err := type_convert.AnyToUint64(c.MustGet(middleware.UserCtx))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, err.Error())
	}
	commentRes, err := ctl.service.CreateComment(userId, commentReq)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, err.Error())
	}
	response.OK(c, commentRes)
}

func (ctl *CommentController) Delete(c *gin.Context) {
	userId, err := type_convert.AnyToUint64(c.MustGet(middleware.UserCtx))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, err.Error())
	}
	commentId, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, err.Error())
	}
	err = ctl.service.DeleteComment(userId, commentId)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, err.Error())
	}
	response.OK(c, nil)
}

func (ctl *CommentController) List(c *gin.Context) {
	videoId, err := strconv.ParseUint(c.Param("videoId"), 10, 64)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, err.Error())
	}
	commentList, err := ctl.service.ListByVideoId(videoId)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, err.Error())
	}
	response.OK(c, commentList)
}

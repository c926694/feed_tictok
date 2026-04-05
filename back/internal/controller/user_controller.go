package controller

import (
	"errors"
	"mime/multipart"
	"net/http"
	"simple_tiktok/internal/dto/req"
	"simple_tiktok/internal/middleware"
	"simple_tiktok/internal/pkg/response"
	"simple_tiktok/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserController struct {
	service *service.UserService
}

func NewUserController(service *service.UserService) *UserController {
	return &UserController{
		service: service,
	}
}
func (ctl *UserController) Register(c *gin.Context) {
	var registerReq req.RegisterReq
	if err := c.ShouldBind(&registerReq); err != nil {
		response.Fail(c, http.StatusBadRequest, err.Error())
		return
	}
	userId, err := ctl.service.Register(c, registerReq)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			response.Fail(c, http.StatusConflict, errors.New("user already exists").Error())
		}
		response.Fail(c, http.StatusInternalServerError, "注册失败")
		return
	}
	response.OK(c, userId)
}

func (ctl *UserController) Login(c *gin.Context) {
	var loginReq req.LoginReq
	if err := c.ShouldBind(&loginReq); err != nil {
		response.Fail(c, http.StatusBadRequest, err.Error())
		return
	}
	token, err := ctl.service.Login(loginReq.Username, loginReq.Password)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.OK(c, token)
}

func (ctl *UserController) GetUserInfo(c *gin.Context) {
	userId := c.MustGet(middleware.UserCtx).(uint64)
	userInfoRes, err := ctl.service.GetUserInfo(userId)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "获取个人信息失败")
		return
	}
	response.OK(c, userInfoRes)
}

func (ctl *UserController) Logout(c *gin.Context) {
	userId := c.MustGet(middleware.UserCtx).(uint64)
	err := ctl.service.Logout(userId)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "退出失败")
		return
	}
	response.OK(c, nil)
}

func (ctl *UserController) UpdateProfile(c *gin.Context) {
	var profileReq req.UpdateUserProfileReq
	_ = c.ShouldBind(&profileReq)
	var avatar *multipart.FileHeader
	file, err := c.FormFile("avatar")
	if err == nil {
		avatar = file
	}
	userId := c.MustGet(middleware.UserCtx).(uint64)
	userInfo, err := ctl.service.UpdateProfile(userId, profileReq.Nickname, avatar)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.OK(c, userInfo)
}

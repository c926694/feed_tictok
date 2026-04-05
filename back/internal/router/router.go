package router

import (
	controller "simple_tiktok/internal/controller"
	"simple_tiktok/internal/middleware"
	"simple_tiktok/internal/mq/event"
	mysql2 "simple_tiktok/internal/repository/mysql"
	service "simple_tiktok/internal/service"

	"github.com/gin-gonic/gin"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func InitRouter(db *gorm.DB, redisClient *redis.Client, conn *amqp.Connection) (*gin.Engine, error) {
	r := gin.Default()
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	if err = ch.ExchangeDeclare(event.LikeVideoExchange, event.LikeVideoExchangeType, true, false, false, false, nil); err != nil {
		return nil, err
	}
	if err = ch.ExchangeDeclare(event.LikeCommentExchange, event.LikeCommentExchangeType, true, false, false, false, nil); err != nil {
		return nil, err
	}
	if err = ch.ExchangeDeclare(event.DeleteVideoExchange, event.DeleteVideoExchangeType, true, false, false, false, nil); err != nil {
		return nil, err
	}
	if err = ch.ExchangeDeclare(event.FollowExchange, event.FollowExchangeType, true, false, false, false, nil); err != nil {
		return nil, err
	}
	if err = ch.ExchangeDeclare(event.VideoHotExchange, event.VideoHotExchangeType, true, false, false, false, nil); err != nil {
		return nil, err
	}
	videoMQ, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	likeMQ, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	followMQ, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	hotMQ, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	userGroup := r.Group("users")
	UserRepo := mysql2.NewUserRepo(db)
	videoRepo := mysql2.NewVideoRepo(db)
	userService := service.NewUserService(UserRepo, videoRepo, redisClient)
	userController := controller.NewUserController(userService)
	{
		userGroup.POST("/login", userController.Login)
		userGroup.POST("/register", userController.Register)
		userGroup.DELETE("/logout", middleware.JWTAuth(redisClient), userController.Logout)
		userGroup.GET("/me", middleware.JWTAuth(redisClient), userController.GetUserInfo)
		userGroup.POST("/me", middleware.JWTAuth(redisClient), userController.UpdateProfile)
	}

	videoGroup := r.Group("videos")
	commentRepo := mysql2.NewCommentRepo(db)
	videoService := service.NewVideoService(videoRepo, UserRepo, redisClient, videoMQ, commentRepo)
	videoController := controller.NewVideoController(videoService)
	{
		videoGroup.POST("/create", middleware.JWTAuth(redisClient), videoController.CreateVideo)
		videoGroup.DELETE("/:id", middleware.JWTAuth(redisClient), videoController.DeleteVideos)
		videoGroup.GET("/me", middleware.JWTAuth(redisClient), videoController.GetMyVideos)
		videoGroup.GET("/feed", middleware.JWTAuth(redisClient), videoController.GetFeedVideos)
		videoGroup.GET("/feed/hot", middleware.JWTAuth(redisClient), videoController.GetFeedHotVideos)
		videoGroup.GET("/feed/follow", middleware.JWTAuth(redisClient), videoController.GetFollowFeedVideos)
		videoGroup.GET("/:id", middleware.JWTAuth(redisClient), videoController.GetVideoInfo)
	}
	likeGroup := r.Group("likes")
	likeService := service.NewLikeService(redisClient, likeMQ)
	likeController := controller.NewLikeController(likeService)
	{
		likeGroup.POST("/video/switchLike/:id", middleware.JWTAuth(redisClient), likeController.LikeVideo)
		likeGroup.POST("/comment/switchLike/:id", middleware.JWTAuth(redisClient), likeController.LikeComment)
	}
	followGroup := r.Group("follows")
	followService := service.NewFollowService(redisClient, followMQ)
	followController := controller.NewFollowController(followService)
	{
		followGroup.POST("/switchFollow/:follower", middleware.JWTAuth(redisClient), followController.Follow)
	}

	commentGroup := r.Group("comments")
	commentService := service.NewCommentService(commentRepo, videoRepo, UserRepo, redisClient, hotMQ)
	commentController := controller.NewCommentController(commentService)
	{
		commentGroup.POST("", middleware.JWTAuth(redisClient), commentController.Create)
		commentGroup.DELETE("/:id", middleware.JWTAuth(redisClient), commentController.Delete)
		commentGroup.GET("/list/:videoId", middleware.JWTAuth(redisClient), commentController.List)
	}
	return r, nil
}

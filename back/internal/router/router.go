package router

import (
	controller2 "simple_tiktok/internal/controller"
	"simple_tiktok/internal/middleware"
	"simple_tiktok/internal/mq/event"
	mysql2 "simple_tiktok/internal/repository/mysql"
	service2 "simple_tiktok/internal/service"

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
	videoMQ, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	likeMQ, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	userGroup := r.Group("users")
	UserRepo := mysql2.NewUserRepo(db)
	userService := service2.NewUserService(UserRepo, redisClient)
	userController := controller2.NewUserController(userService)
	{
		userGroup.POST("/login", userController.Login)
		userGroup.POST("/register", userController.Register)
		userGroup.DELETE("/logout", middleware.JWTAuth(redisClient), userController.Logout)
		userGroup.GET("/me", middleware.JWTAuth(redisClient), userController.GetUserInfo)
	}

	videoGroup := r.Group("videos")
	videoRepo := mysql2.NewVideoRepo(db)
	commentRepo := mysql2.NewCommentRepo(db)
	videoServide := service2.NewVideoService(videoRepo, redisClient, videoMQ, commentRepo)
	videoController := controller2.NewVideoController(videoServide)
	{
		videoGroup.POST("/create", middleware.JWTAuth(redisClient), videoController.CreateVideo)
		videoGroup.DELETE("/:id", middleware.JWTAuth(redisClient), videoController.DeleteVideos)
		videoGroup.GET("/feed", videoController.GetFeedVideos)
		videoGroup.GET("/feed/hot", videoController.GetFeedHotVideos)
		videoGroup.GET("/:id", videoController.GetVideoInfo)
	}
	likeGroup := r.Group("likes")
	likeService := service2.NewLikeService(redisClient, likeMQ)
	likeController := controller2.NewLikeController(likeService)
	{
		likeGroup.POST("/video/switchLike/:id", middleware.JWTAuth(redisClient), likeController.LikeVideo)
		likeGroup.POST("/comment/switchLike/:id", middleware.JWTAuth(redisClient), likeController.LikeComment)
	}
	followGroup := r.Group("follows")
	followRepo := mysql2.NewFollowRepo(db)
	followService := service2.NewFollowService(followRepo, redisClient)
	followController := controller2.NewFollowController(followService)
	{
		followGroup.POST("/switchFollow/:follower", middleware.JWTAuth(redisClient), followController.Follow)
	}

	commentGroup := r.Group("comments")
	commentService := service2.NewCommentService(commentRepo, videoRepo)
	commentController := controller2.NewCommentController(commentService)
	{
		commentGroup.POST("/", middleware.JWTAuth(redisClient), commentController.Create)
		commentGroup.DELETE("/:id", middleware.JWTAuth(redisClient), commentController.Delete)
		commentGroup.GET("/list/:videoId", commentController.List)
	}
	return r, nil
}

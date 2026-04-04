package middleware

import (
	"context"
	"fmt"
	"net/http"
	"simple_tiktok/internal/initialize"
	"simple_tiktok/internal/pkg/jwt"
	"simple_tiktok/internal/pkg/response"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type JWT struct {
	redisClient *redis.Client
}

const (
	TokenKey     = "jwt_token:%d"
	UserCtx      = "userId"
	UserNickName = "userNickName"
)

func JWTAuth(redisClient *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := strings.TrimSpace(c.GetHeader("Authorization"))
		if authHeader == "" {
			response.Fail(c, http.StatusUnauthorized, "missing authorization header")
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			response.Fail(c, http.StatusUnauthorized, "invalid authorization format")
			c.Abort()
			return
		}

		claims, err := jwt.ParseToken(parts[1])
		if err != nil {
			response.Fail(c, http.StatusUnauthorized, "invalid token")
			c.Abort()
			return
		}
		userId := claims.UserID
		nickName := claims.UserNickName
		key := fmt.Sprintf(TokenKey, userId)
		ctx := context.Background()
		token, err := redisClient.Get(ctx, key).Result()
		if err != nil {
			response.Fail(c, http.StatusUnauthorized, "redis无token")
			c.Abort()
		}
		//刷新token
		expire := initialize.AppConfig.JWT.ExpireHours
		_, err = redisClient.Set(ctx, key, token, time.Duration(expire)*time.Hour).Result()
		if err != nil {
			response.Fail(c, http.StatusUnauthorized, "设置token失败")
			c.Abort()
		}
		c.Set(UserCtx, userId)
		c.Set(UserNickName, nickName)
		c.Next()
	}
}

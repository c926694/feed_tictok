package jwt

import (
	"fmt"
	"simple_tiktok/internal/initialize"
	"time"

	gjwt "github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID       uint64 `json:"user_id"`
	UserNickName string `json:"user_nick_name"`
	gjwt.RegisteredClaims
}

func GenerateToken(userID uint64, nickName string) (string, error) {
	secret, expireHours, err := loadJWTConfig()
	if err != nil {
		return "", err
	}

	now := time.Now()
	claims := Claims{
		UserID:       userID,
		UserNickName: nickName,
		RegisteredClaims: gjwt.RegisteredClaims{
			IssuedAt:  gjwt.NewNumericDate(now),
			ExpiresAt: gjwt.NewNumericDate(now.Add(time.Duration(expireHours) * time.Hour)),
		},
	}

	token := gjwt.NewWithClaims(gjwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func ParseToken(tokenString string) (*Claims, error) {
	secret, _, err := loadJWTConfig()
	if err != nil {
		return nil, err
	}

	token, err := gjwt.ParseWithClaims(tokenString, &Claims{}, func(token *gjwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*gjwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	return claims, nil
}

func loadJWTConfig() (string, int64, error) {
	if initialize.AppConfig == nil {
		return "", 0, fmt.Errorf("app config is not initialized")
	}
	secret := initialize.AppConfig.JWT.Secret
	if secret == "" {
		return "", 0, fmt.Errorf("jwt secret is empty")
	}
	expireHours := initialize.AppConfig.JWT.ExpireHours
	if expireHours <= 0 {
		expireHours = 72
	}
	return secret, expireHours, nil
}

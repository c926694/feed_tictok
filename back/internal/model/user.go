package model

import "time"

type User struct {
	ID            uint64    `gorm:"primaryKey"`
	Username      string    `gorm:"size:64;uniqueIndex;not null"`
	Password      string    `gorm:"size:128;not null"`
	NickName      string    `gorm:"size:64;not null"`
	AvatarURL     string    `gorm:"size:255"`
	FollowCount   int64     `gorm:"default:0"`
	FollowerCount int64     `gorm:"default:0"`
	CreatedAt     time.Time `gorm:"autoCreateTime;not null"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime;not null"`
}

package model

import "time"

type Comment struct {
	ID        uint64    `gorm:"primaryKey"`
	Commenter uint64    `gorm:"commenter" `
	VideoID   uint64    `gorm:"index;not null"`
	Content   string    `gorm:"size:500;not null"`
	LikeCount int64     `gorm:"default:0"`
	CreatedAt time.Time `gorm:"autoCreateTime;not null"`
	UpdatedAt time.Time ` gorm:"autoUpdateTime;not null"`
}

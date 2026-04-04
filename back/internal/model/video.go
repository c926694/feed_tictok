package model

import "time"

type Video struct {
	ID           uint64    `gorm:"primaryKey"`
	AuthorID     uint64    `gorm:"index;not null;"`
	AuthorName   string    `gorm:"size:255;not null;"`
	PlayURL      string    `gorm:"size:255;not null"`
	CoverURL     string    `gorm:"size:255;not null"`
	Title        string    `gorm:"size:255;not null"`
	LikeCount    int64     `gorm:"default:0"`
	CommentCount int64     `gorm:"default:0"`
	CreatedAt    time.Time `gorm:"index" gorm:"autoCreateTime;not null"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"autoUpdateTime;not null"`
}

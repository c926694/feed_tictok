package model

import "time"

type Follow struct {
	ID        uint64    `gorm:"primaryKey" json:"id"`
	Following uint64    `gorm:"not null;uniqueIndex:idx_user_follow" json:"user_id"`
	Follower  uint64    `gorm:"not null;uniqueIndex:idx_user_follow" json:"follow_user_id"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime;not null"`
}

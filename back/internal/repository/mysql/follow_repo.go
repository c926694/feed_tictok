package mysql

import (
	"simple_tiktok/internal/model"

	"gorm.io/gorm"
)

type FollowRepo struct {
	db *gorm.DB
}

func NewFollowRepo(db *gorm.DB) *FollowRepo {
	return &FollowRepo{db: db}
}

func (r *FollowRepo) Follow(follow *model.Follow) error {
	return r.db.Save(follow).Error
}

func (r *FollowRepo) DeleteFollow(follow *model.Follow) error {
	return r.db.Where("follower = ? and following = ?", follow.Follower, follow.Following).
		Delete(&model.Follow{}).Error
}

func (r *FollowRepo) DB() *gorm.DB {
	return r.db
}


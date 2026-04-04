package initialize

import (
	model2 "simple_tiktok/internal/model"

	"gorm.io/gorm"
)

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&model2.User{},
		&model2.Video{},
		&model2.Comment{},
		&model2.Follow{},
	)
}

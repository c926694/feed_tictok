package mysql

import (
	"context"
	"errors"
	"simple_tiktok/internal/model"
	"simple_tiktok/internal/pkg/random_nickname"

	"gorm.io/gorm"
)

type UserRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) Create(ctx context.Context, user *model.User) error {
	_ = ctx
	_ = user
	return errors.New("not implemented")
}

func (r *UserRepo) GetUserByID(userID uint64) (*model.User, error) {
	var user model.User
	err := r.db.Where("id = ?", userID).First(&user).Error
	return &user, err
}

func (r *UserRepo) GetUserByUserNameAndPassword(username string) (*model.User, error) {
	var user model.User
	err := r.db.Model(&user).Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepo) CreateUser(username string, password string, avatar string) (*model.User, error) {
	user := model.User{
		Username:  username,
		Password:  password,
		AvatarURL: avatar,
		NickName:  random_nickname.GenerateNickname(),
	}
	err := r.db.Model(&user).Create(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepo) IncFollowCount(userID uint64) error {
	return r.db.Model(&model.User{}).Where("id = ?", userID).Update("follow_count", gorm.Expr("follow_count + 1")).Error
}

func (r *UserRepo) DecFollowCount(userID uint64) error {
	return r.db.Model(&model.User{}).Where("id = ?", userID).Update("follow_count", gorm.Expr("CASE WHEN follow_count > 0 THEN follow_count - 1 ELSE 0 END")).Error
}

func (r *UserRepo) IncFollowerCount(userID uint64) error {
	return r.db.Model(&model.User{}).Where("id = ?", userID).Update("follower_count", gorm.Expr("follower_count + 1")).Error
}

func (r *UserRepo) DecFollowerCount(userID uint64) error {
	return r.db.Model(&model.User{}).Where("id = ?", userID).Update("follower_count", gorm.Expr("CASE WHEN follower_count > 0 THEN follower_count - 1 ELSE 0 END")).Error
}

func (r *UserRepo) UpdateProfile(userID uint64, updates map[string]any) error {
	if len(updates) == 0 {
		return nil
	}
	return r.db.Model(&model.User{}).Where("id = ?", userID).Updates(updates).Error
}

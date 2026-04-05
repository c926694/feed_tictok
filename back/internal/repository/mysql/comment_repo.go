package mysql

import (
	"simple_tiktok/internal/model"

	"gorm.io/gorm"
)

type CommentRepo struct {
	db *gorm.DB
}

func NewCommentRepo(db *gorm.DB) *CommentRepo {
	return &CommentRepo{db: db}
}

func (r *CommentRepo) Save(comment *model.Comment) error {
	return r.db.Save(comment).Error
}

func (r *CommentRepo) WithTx(tx *gorm.DB) *CommentRepo {
	return &CommentRepo{db: tx}
}

func (r *CommentRepo) GetById(id uint64) (model.Comment, error) {
	var comment model.Comment
	err := r.db.Where("id = ?", id).First(&comment).Error
	return comment, err
}

func (r *CommentRepo) DeleteComment(id uint64) error {
	return r.db.Where("id = ?", id).Delete(&model.Comment{}).Error
}

func (r *CommentRepo) ListByVideoId(videoId uint64) ([]model.Comment, error) {
	var comments []model.Comment
	err := r.db.Where("video_id = ?", videoId).Find(&comments).Error
	return comments, err
}

func (r *CommentRepo) DB() *gorm.DB {
	return r.db
}

func (r *CommentRepo) DeleteByVideoId(videoId uint64) error {
	return r.db.Where("video_id = ?", videoId).Delete(&model.Comment{}).Error
}

func (r *CommentRepo) IncCommentLikeCount(id uint64) error {
	return r.db.Model(&model.Comment{}).Where("id = ?", id).Update("like_count", gorm.Expr("like_count+1")).Error
}

func (r *CommentRepo) DecCommentLikeCount(id uint64) error {
	return r.db.Model(&model.Comment{}).Where("id = ?", id).Update("like_count", gorm.Expr("like_count-1")).Error
}

package mysql

import (
	"simple_tiktok/internal/model"
	"time"

	"gorm.io/gorm"
)

type VideoRepo struct {
	db *gorm.DB
}

func NewVideoRepo(db *gorm.DB) *VideoRepo {
	return &VideoRepo{db: db}
}

func (r *VideoRepo) CreateVideo(video *model.Video) error {
	return r.db.Create(video).Error
}

func (r *VideoRepo) WithTx(tx *gorm.DB) *VideoRepo {
	return &VideoRepo{db: tx}
}

func (r *VideoRepo) GetFeedVideos(ids []uint64) ([]model.Video, error) {
	videoList := []model.Video{}
	err := r.db.Model(&model.Video{}).Where("id in ?", ids).Find(&videoList).Error
	return videoList, err
}

func (r *VideoRepo) GetFollowFeedVideosByAuthors(authorIDs []uint64, limit uint64, lastCreatedAt *time.Time) ([]model.Video, error) {
	videoList := make([]model.Video, 0, limit)
	query := r.db.Model(&model.Video{}).Where("author_id in ?", authorIDs)
	if lastCreatedAt != nil {
		query = query.Where("created_at < ?", *lastCreatedAt)
	}
	err := query.Order("created_at desc").Limit(int(limit)).Find(&videoList).Error
	return videoList, err
}

func (r *VideoRepo) ListByAuthorID(authorID uint64, limit uint64) ([]model.Video, error) {
	videoList := make([]model.Video, 0, limit)
	err := r.db.Model(&model.Video{}).Where("author_id = ?", authorID).Order("created_at desc").Limit(int(limit)).Find(&videoList).Error
	return videoList, err
}

func (r *VideoRepo) CountByAuthorID(authorID uint64) (int64, error) {
	var count int64
	err := r.db.Model(&model.Video{}).Where("author_id = ?", authorID).Count(&count).Error
	return count, err
}

func (r *VideoRepo) GetVideoById(id uint64) (model.Video, error) {
	video := model.Video{}
	err := r.db.First(&video, id).Error
	return video, err
}

func (r *VideoRepo) IncVideoLikeCount(id uint64) error {
	return r.db.Model(&model.Video{}).Where("id = ?", id).Update("like_count", gorm.Expr("like_count+1")).Error
}
func (r *VideoRepo) DecVideoDislikeCount(id uint64) error {
	return r.db.Model(&model.Video{}).Where("id = ?", id).Update("like_count", gorm.Expr("like_count-1")).Error
}

func (r *VideoRepo) DeleteVideoById(id uint64) error {
	return r.db.Delete(&model.Video{}, id).Error
}

func (r *VideoRepo) DB() *gorm.DB {
	return r.db
}

func (r *VideoRepo) UpdateCommentCount(id uint64) error {
	return r.db.Model(&model.Video{}).Where("id = ?", id).Update("comment_count", gorm.Expr("comment_count+1")).Error
}

func (r *VideoRepo) DeleteCommentCount(id uint64) error {
	return r.db.Model(&model.Video{}).Where("id = ?", id).Update("comment_count", gorm.Expr("comment_count-1")).Error
}

func (r *VideoRepo) UpdateAuthorNameByAuthorID(authorID uint64, authorName string) error {
	return r.db.Model(&model.Video{}).Where("author_id = ?", authorID).Update("author_name", authorName).Error
}

package service

import (
	"errors"
	"simple_tiktok/internal/dto/req"
	"simple_tiktok/internal/dto/res"
	"simple_tiktok/internal/model"
	mysql2 "simple_tiktok/internal/repository/mysql"
)

type CommentService struct {
	commentRepo *mysql2.CommentRepo
	videoRepo   *mysql2.VideoRepo
}

func NewCommentService(commentRepo *mysql2.CommentRepo, videoRepo *mysql2.VideoRepo) *CommentService {
	return &CommentService{
		commentRepo: commentRepo,
		videoRepo:   videoRepo,
	}
}

func (s *CommentService) CreateComment(id uint64, req req.CommentReq) (res.CommentRes, error) {
	tx := s.commentRepo.DB()
	tx.Begin()
	comment := &model.Comment{
		Content:   req.Content,
		VideoID:   req.VideoId,
		Commenter: id,
	}
	err := s.commentRepo.Save(comment)
	if err != nil {
		tx.Rollback()
		return res.CommentRes{}, err
	}
	//更新comment_count
	err = s.videoRepo.UpdateCommentCount(comment.VideoID)
	if err != nil {
		tx.Rollback()
		return res.CommentRes{}, err
	}
	tx.Commit()
	commentRes := res.CommentRes{
		Id:        comment.ID,
		VideoId:   comment.VideoID,
		Commenter: comment.Commenter,
		Content:   comment.Content,
		CreatedAt: comment.CreatedAt,
	}
	return commentRes, nil
}

func (s *CommentService) DeleteComment(userId uint64, id uint64) error {
	comment, err := s.commentRepo.GetById(id)
	if err != nil {
		return err
	}
	if comment.Commenter != userId {
		return errors.New("无法删除他人评论")
	}
	tx := s.commentRepo.DB()
	tx.Begin()

	err = s.commentRepo.DeleteComment(id)
	if err != nil {
		tx.Rollback()
		return err
	}
	err = s.videoRepo.DeleteCommentCount(comment.VideoID)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (s *CommentService) ListByVideoId(videoId uint64) ([]res.CommentRes, error) {
	commentList, err := s.commentRepo.ListByVideoId(videoId)
	if err != nil {
		return nil, err
	}
	commentResList := []res.CommentRes{}
	for _, comment := range commentList {
		commentRes := res.CommentRes{
			Id:        comment.ID,
			VideoId:   comment.VideoID,
			Commenter: comment.Commenter,
			Content:   comment.Content,
			CreatedAt: comment.CreatedAt,
		}
		commentResList = append(commentResList, commentRes)
	}
	return commentResList, nil
}

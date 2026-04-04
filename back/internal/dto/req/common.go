package req

import (
	"mime/multipart"
)

type EmptyRequest struct{}

type LoginReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
type RegisterReq struct {
	LoginReq
	RePassword string `json:"re_password"`
}

type UploadVideoReq struct {
	Title string                `json:"title"`
	Cover *multipart.FileHeader `json:"cover"`
	Play  *multipart.FileHeader `json:"play"`
}

type CommentReq struct {
	VideoId uint64 `json:"video_id"`
	Content string `json:"content"`
}

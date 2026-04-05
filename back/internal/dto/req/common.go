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

type UpdateUserProfileReq struct {
	Nickname string                `form:"nickname"`
	Avatar   *multipart.FileHeader `form:"avatar"`
}

type UploadVideoReq struct {
	Title       string                `json:"title" form:"title"`
	Description string                `json:"description" form:"description"`
	Cover       *multipart.FileHeader `json:"cover" form:"cover"`
	Play        *multipart.FileHeader `json:"play" form:"play"`
}

type CommentReq struct {
	VideoId uint64 `json:"video_id"`
	Content string `json:"content"`
}

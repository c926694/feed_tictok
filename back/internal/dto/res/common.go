package res

import "time"

type EmptyData struct{}

type UserInfoRes struct {
	Nickname      string `json:"nickname"`
	AvatarURL     string `json:"avatar_URL"`
	FollowCount   int64  `json:"follow_count"`
	FollowerCount int64  `json:"follower_count"`
}

type VideoRes struct {
	Id  uint64 `json:"id"`
	Url string `json:"url"`
}

type VideoInfoRes struct {
	Id           uint64    `json:"id"`
	AuthorName   string    `json:"author_name"`
	CoverURL     string    `json:"coverURL"`
	PlayURL      string    `json:"playURL"`
	CommentCount int64     `json:"comment_count"`
	LikeCount    int64     `json:"like_count"`
	CreatedAt    time.Time `json:"created_at"`
}

type FeedVideoRes struct {
	FeedVideoList []VideoInfoRes `json:"feed_video_list"`
	LastScore     float64        `json:"last_score"`
}

type CommentRes struct {
	Id        uint64    `json:"id"`
	Commenter uint64    `json:"commenter"`
	VideoId   uint64    `json:"video_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

type LikeVideoRes struct {
	VideoId uint64 `json:"video_id"`
	IsLiked bool   `json:"is_liked"`
}

type LikeCommentRes struct {
	CommentId uint64 `json:"comment_id"`
	IsLiked   bool   `json:"is_liked"`
}

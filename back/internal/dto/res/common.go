package res

import "time"

type EmptyData struct{}

type UserInfoRes struct {
	UserID        uint64 `json:"user_id"`
	Username      string `json:"username"`
	Nickname      string `json:"nickname"`
	AvatarURL     string `json:"avatar_URL"`
	FollowCount   int64  `json:"follow_count"`
	FollowerCount int64  `json:"follower_count"`
	VideoCount    int64  `json:"video_count"`
}

type VideoRes struct {
	Id  uint64 `json:"id"`
	Url string `json:"url"`
}

type VideoInfoRes struct {
	Id           uint64    `json:"id"`
	AuthorID     uint64    `json:"author_id"`
	AuthorName   string    `json:"author_name"`
	AuthorAvatar string    `json:"author_avatar"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	CoverURL     string    `json:"coverURL"`
	PlayURL      string    `json:"playURL"`
	CommentCount int64     `json:"comment_count"`
	LikeCount    int64     `json:"like_count"`
	IsLiked      bool      `json:"is_liked"`
	IsFollow     bool      `json:"is_follow"`
	CreatedAt    time.Time `json:"created_at"`
}

type FeedVideoRes struct {
	FeedVideoList []VideoInfoRes `json:"feed_video_list"`
	LastScore     float64        `json:"last_score"`
}

type HotFeedVideoRes struct {
	FeedVideoList []VideoInfoRes `json:"feed_video_list"`
	NextOffset    uint64         `json:"next_offset"`
	HasMore       bool           `json:"has_more"`
	Interval      int            `json:"interval"`
}

type CommentRes struct {
	Id        uint64      `json:"id"`
	Commenter uint64      `json:"commenter"`
	VideoId   uint64      `json:"video_id"`
	Content   string      `json:"content"`
	LikeCount int64       `json:"like_count"`
	IsLiked   bool        `json:"is_liked"`
	Author    UserInfoRes `json:"author"`
	CreatedAt time.Time   `json:"created_at"`
}

type LikeVideoRes struct {
	VideoId uint64 `json:"video_id"`
	IsLiked bool   `json:"is_liked"`
}

type LikeCommentRes struct {
	CommentId uint64 `json:"comment_id"`
	IsLiked   bool   `json:"is_liked"`
}

type FollowRes struct {
	Following uint64 `json:"following"`
	IsFollow  bool   `json:"is_follow"`
}

package event

type LikeVideoEvent struct {
	EventType string `json:"event"`
	VideoId   uint64 `json:"videoId"`
}

type LikeCommentEvent struct {
	EventType string `json:"event"`
	CommentId uint64 `json:"commentId"`
}

type DeleteVideoEvent struct {
	PlayURL  string `json:"playUrl"`
	CoverURL string `json:"coverUrl"`
}

const (
	Like    = "like"
	Dislike = "dislike"
)

const (
	LikeVideoExchange     = "like.video.exchange"
	LikeVideoRoutingKey   = "like.video"
	LikeVideoQueue        = "like.video.queue"
	LikeVideoExchangeType = "direct"
)

const (
	LikeCommentExchange     = "like.comment.exchange"
	LikeCommentRoutingKey   = "like.comment"
	LikeCommentQueue        = "like.comment.queue"
	LikeCommentExchangeType = "direct"
)

const (
	DeleteVideoExchange     = "video.delete.exchange"
	DeleteVideoRoutingKey   = "video.delete"
	DeleteVideoQueue        = "video.delete.queue"
	DeleteVideoExchangeType = "direct"
)

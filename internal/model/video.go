package model

import (
	"mime/multipart"
	"time"
)

// PublishVideoInput 投稿操作 Input
type PublishVideoInput struct {
	Id        int64                 `json:"id"`
	AuthorId  int64                 `json:"author_id"`
	Title     string                `json:"title"`
	PlayUrl   string                `json:"play_url"`
	CoverUrl  string                `json:"cover_url"`
	VideoData *multipart.FileHeader `json:"video_data"`
}

// PublishVideoOutput 投稿操作 Output
type PublishVideoOutput struct {
}

// VideoFeedInput 视频流 Input
type VideoFeedInput struct {
	LatestTime time.Time `json:"latest_time"`
	UserId     int64     `json:"user_id"`
}

// VideoFeedOutput 视频流 output
type VideoFeedOutput struct {
	NextTime  int64        `json:"next_time"`
	VideoList []*VideoItem `json:"video_list"`
}

type VideoItem struct {
	Id            int64     `json:"id"`
	AuthorId      int64     `json:"author_id"`
	FavoriteCount int64     `json:"favorite_count"`
	CommentCount  int64     `json:"comment_count"`
	Title         string    `json:"title"`
	PlayUrl       string    `json:"play_url"`
	CoverUrl      string    `json:"cover_url"`
	IsFavorite    bool      `json:"is_favorite"`
	Author        *UserItem `json:"author"`
	CreatedAt     time.Time `json:"-"`
}

// GetVideoPublishedListInput 发布视频列表 Input
type GetVideoPublishedListInput struct {
	UserId int64 `json:"user_id"`
}

// GetVideoPublishedListOutput 发布视频列表 Output
type GetVideoPublishedListOutput struct {
	VideoList []*VideoItem `json:"video_list"`
}

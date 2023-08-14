package v1

import (
	"mime/multipart"
	"simple_tiktok_rime/internal/consts"
	"time"
)

// PublishVideoReq 投稿视频请求
type PublishVideoReq struct {
	Title string                `json:"title"`
	Token string                `json:"token"`
	Data  *multipart.FileHeader `json:"data,omitempty"`
}

// PublishVideoResp 投稿视频响应
type PublishVideoResp struct {
	*consts.ResponseData
}

// VideoFeedReq 视频流请求
type VideoFeedReq struct {
	LatestTime time.Time `json:"latest_time"`
	Token      string    `json:"token"`
}

// VideoFeedResp 视频流响应
type VideoFeedResp struct {
	*consts.ResponseData
	NextTime  time.Time `json:"next_time"`
	VideoList any       `json:"video_list"`
}

// GetVidePublishedListReq 视频发布列表请求
type GetVidePublishedListReq struct {
	UserId int64  `json:"user_id"`
	Token  string `json:"token"`
}

// GetVidePublishedListResp 视频发布列表响应
type GetVidePublishedListResp struct {
	*consts.ResponseData
	VideoList any `json:"video_list"`
}

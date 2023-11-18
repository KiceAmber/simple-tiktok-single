package v1

import "simple_tiktok_single/internal/consts"

// FavoriteActionReq 点赞操作请求
type FavoriteActionReq struct {
	Token      string `json:"token"`
	VideoId    int64  `json:"video_id"`
	ActionType string `json:"action_type"`
}

// FavoriteActionResp 点赞操作响应
type FavoriteActionResp struct {
	*consts.ResponseData
}

// GetFavoriteVideoListReq 获取视频点赞列表请求
type GetFavoriteVideoListReq struct {
	Token  string `json:"token"`
	UserId int64  `json:"user_id"`
}

// GetFavoriteVideoListResp 获取视频点赞列表响应
type GetFavoriteVideoListResp struct {
	//*consts.ResponseData
	StatusCode string `json:"status_code"`
	StatusMsg  string `json:"status_msg"`
	VideoList  any    `json:"video_list"`
}

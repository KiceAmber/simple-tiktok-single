package v1

import "simple_tiktok_single/internal/consts"

// CommentActionReq 评论操作请求
type CommentActionReq struct {
	Token       string `json:"token"`
	ActionType  string `json:"action_type"`
	CommentText string `json:"comment_text"`
	VideoId     int64  `json:"video_id"`
	CommentId   int64  `json:"comment_id"`
}

// CommentActionReq 评论操作响应
type CommentActionResp struct {
	*consts.ResponseData
	Comment any `json:"comment"`
}

// GetCommentList 显示视频评论请求
type GetCommentListReq struct {
	Token   string `json:"token"`
	VideoId int64  `json:"video_id"`
}

// GetCommentListResp 显示视频评论请求
type GetCommentListResp struct {
	*consts.ResponseData
	CommentList any `json:"comment_list"`
}

package v1

import "simple_tiktok_rime/internal/consts"

// MessageActionReq 消息操作请求
type MessageActionReq struct {
	Token      string `json:"token"`
	ToUserId   int64  `json:"to_user_id"`
	ActionType string `json:"action_type"`
	Content    string `json:"content"`
}

// MessageActionResp 消息操作响应
type MessageActionResp struct {
	*consts.ResponseData
}

// GetMessageListReq 获取消息列表请求
type GetMessageListReq struct {
	Token    string `json:"token"`
	ToUserId int64  `json:"to_user_id"`
}

// GetMessageListResp 获取消息列表响应
type GetMessageListResp struct {
	*consts.ResponseData
	MessageList any `json:"message_list"`
}

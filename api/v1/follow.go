package v1

import "simple_tiktok_rime/internal/consts"

// FollowActionReq 关注操作请求
type FollowActionReq struct {
	Token      string `json:"token"`
	ActionType string `json:"action_type"`
	ToUserId   int64  `json:"to_user_id"`
}

// FollowActionResp 关注操作响应
type FollowActionResp struct {
	*consts.ResponseData
}

// GetFollowListReq 获取用户关注列表请求
type GetFollowListReq struct {
	Token  string `json:"token"`
	UserId int64  `json:"user_id"`
}

// GetFollowListResp 获取用户关注列表响应
type GetFollowListResp struct {
	//*consts.ResponseData
	StatusCode string `json:"status_code"`
	StatusMsg  string `json:"status_msg"`
	FollowList any    `json:"user_list"`
}

// GetFollowerListReq 获取用户粉丝列表请求
type GetFollowerListReq struct {
	Token  string `json:"token"`
	UserId int64  `json:"user_id"`
}

// GetFollowerListResp 获取用户粉丝列表响应
type GetFollowerListResp struct {
	//*consts.ResponseData
	StatusCode string `json:"status_code"`
	StatusMsg  string `json:"status_msg"`
	FollowList any    `json:"user_list"`
}

// GetFriendListReq 获取好友列表请求
type GetFriendListReq struct {
	Token  string `json:"token"`
	UserId int64  `json:"user_id"`
}

// GetFriendListResp 获取好友列表响应
type GetFriendListResp struct {
	//*consts.ResponseData
	StatusCode string `json:"status_code"`
	StatusMsg  string `json:"status_msg"`
	FriendList any    `json:"user_list"`
}

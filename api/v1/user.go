package v1

import "simple_tiktok_single/internal/consts"

// UserRegisterReq 用户注册的请求结构体
type UserRegisterReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// UserRegisterResp 用户注册的响应结构体
type UserRegisterResp struct {
	*consts.ResponseData
	UserId int64  `json:"user_id"`
	Token  string `json:"token"`
}

// UserLoginReq 用户登录的请求结构体
type UserLoginReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// UserLoginResp 用户登录的响应结构体
type UserLoginResp struct {
	*consts.ResponseData
	UserId int64  `json:"user_id"`
	Token  string `json:"token"`
}

// GetUserInfoReq 显示用户信息请求结构体
type GetUserInfoReq struct {
	UserId int64  `json:"user_id"`
	Token  string `json:"token"`
}

// GetUserInfoResp 显示用户信息响应结构体
type GetUserInfoResp struct {
	*consts.ResponseData
	User any `json:"user"`
}

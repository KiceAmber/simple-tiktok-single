package consts

import (
	"errors"
	"time"
)

// 封装常用的常量
const (
	TokenExpireDuration = time.Hour * 2  // Token 的过期时间
	CtxUserIdKey        = "ctxUserIdKey" // 用户 ID
)

// 封装响应码
const (
	// 基础响应码
	CodeSuccess RespCode = 0 + iota
	CodeServerBusy
	CodeInvalidParam

	// 用户响应码
	CodeNeedLogin
	CodeInvalidToken
	CodeUserExists
	CodeLoginFailed

	// 视频响应码
	CodeInvalidFileHeader
	CodeInvalidTimeStamp

	// 点赞响应码
	CodeUserFavoritedVideo
	CodeUserNotFavoriteVideo
)

// 封装错误码
var (
	ErrInvalidVideoExt      = errors.New("无效的视频扩展名")
	ErrUserNotExists        = errors.New("用户不存在")
	ErrUserExists           = errors.New("用户已存在")
	ErrUserFavoritedVideo   = errors.New("用户已对该视频点赞")
	ErrUserNotFavoriteVideo = errors.New("用户未对该视频点赞")
)

// 其他常量
const (
	DefaultBackgroundImage = "http://ryr42bm4i.hn-bkt.clouddn.com/bg.jpeg"    // 默认的个人主页顶部大图
	DefaultAvatar          = "http://ryr42bm4i.hn-bkt.clouddn.com/avatar.png" // 默认的头像图片
	DefaultSignature       = "这个人很懒，什么都没有填写"                                  // 默认的个人简介
)

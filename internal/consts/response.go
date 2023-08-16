package consts

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// RespCode 响应码封装
type RespCode int32

var codeMsg = map[RespCode]string{
	// 基础 Code
	CodeSuccess:      "成功",
	CodeServerBusy:   "服务繁忙",
	CodeInvalidParam: "无效的参数",

	// 用户 Code
	CodeNeedLogin:    "需要登录",
	CodeInvalidToken: "无效的token",
	CodeUserExists:   "用户名已存在",
	CodeLoginFailed:  "用户名或密码错误",

	// 视频 Code
	CodeInvalidFileHeader: "无效的视频文件头",
	CodeInvalidTimeStamp:  "无效的时间戳",

	// 点赞 Code
	CodeUserFavoritedVideo:   "用户已对该视频点赞",
	CodeUserNotFavoriteVideo: "用户未对该视频点赞",

	// 关注 Code
	CodeUserFollowedTargetUser: "已对关注用户",
	CodeUserNotFollowTargetUser: "未关注该用户",
}

func (c RespCode) GetMsg() string {
	msg, ok := codeMsg[c]
	if !ok {
		msg = codeMsg[c]
	}
	return msg
}

// ResponseData 返回的响应信息结构体
//type ResponseData struct {
//	StatusCode RespCode    `json:"status_code"`
//	StatusMsg  string      `json:"status_msg"`
//	Data       interface{} `json:"data,omitempty"`
//}

// ResponseData 返回的响应信息结构体
type ResponseData struct {
	StatusCode RespCode `json:"status_code"`
	StatusMsg  string   `json:"status_msg"`
}

// ResponseError 错误响应
func ResponseErrorData(code RespCode) *ResponseData {
	return &ResponseData{
		StatusCode: code,
		StatusMsg:  code.GetMsg(),
	}
}

// ResponseSuccess 成功响应
func ResponseSuccessData(message string) *ResponseData {
	responseSuccessData := ResponseData{
		StatusCode: CodeSuccess,
		StatusMsg:  CodeSuccess.GetMsg(),
	}
	if message != "" {
		responseSuccessData.StatusMsg = message
	}
	return &responseSuccessData
}

// ResponseError 错误响应
func ResponseError(ctx *gin.Context, data interface{}) {
	ctx.JSON(http.StatusOK, data)
}

// ResponseSuccess 成功响应
func ResponseSuccess(ctx *gin.Context, data interface{}) {
	ctx.JSON(http.StatusOK, data)
}

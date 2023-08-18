package controller

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	v1 "simple_tiktok_rime/api/v1"
	"simple_tiktok_rime/internal/consts"
	"simple_tiktok_rime/internal/model"
	"simple_tiktok_rime/internal/service"
	"simple_tiktok_rime/pkg/jwt"
	"strconv"
)

// MessageAction 消息操作
func MessageAction(ctx *gin.Context) {

	var req = new(v1.MessageActionReq)

	// 接收参数
	req.Token = ctx.Query("token")
	req.Content = ctx.Query("content")
	req.ActionType = ctx.Query("action_type")
	toUserId, err := strconv.ParseInt(ctx.Query("to_user_id"), 10, 64)
	if err != nil {
		zap.L().Error("GetMessage strconv.ParseINt Failed", zap.Error(err))
		consts.ResponseError(ctx, v1.MessageActionResp{
			ResponseData: consts.ResponseErrorData(consts.CodeInvalidParam),
		})
		return
	}
	req.ToUserId = toUserId

	// 解析 Token，获取 UserId
	myClaims, err := jwt.ParseToken(req.Token)
	if err != nil {
		zap.L().Error("jwt.ParseToken Failed", zap.Error(err))
		consts.ResponseError(ctx, &v1.MessageActionResp{
			ResponseData: consts.ResponseErrorData(consts.CodeInvalidToken),
		})
		return
	}

	// 业务处理
	_, err = service.Chat().MessageAction(&model.MessageActionInput{
		UserId:     myClaims.Id,
		ToUserId:   req.ToUserId,
		Content:    req.Content,
		ActionType: req.ActionType,
	})
	if err != nil {
		zap.L().Error("service.Chat().MessageAction Failed", zap.Error(err))
		consts.ResponseError(ctx, &v1.MessageActionResp{
			ResponseData: consts.ResponseErrorData(consts.CodeInvalidToken),
		})
		return
	}

	// 返回响应
	consts.ResponseSuccess(ctx, v1.MessageActionResp{
		ResponseData: consts.ResponseSuccessData("发送消息成功"),
	})
}

// MessageList 获取消息列表
func MessageList(ctx *gin.Context) {
	var req = new(v1.GetMessageListReq)

	// 接收参数
	req.Token = ctx.Query("token")

	toUserId, err := strconv.ParseInt(ctx.Query("to_user_id"), 10, 64)
	if err != nil {
		zap.L().Error("MessageList strconv.ParseInt Failed", zap.Error(err))
		consts.ResponseError(ctx, v1.GetMessageListResp{
			ResponseData: consts.ResponseErrorData(consts.CodeInvalidParam),
		})
		return
	}
	req.ToUserId = toUserId

	// 解析 Token，获取 UserId
	myClaims, err := jwt.ParseToken(req.Token)
	if err != nil {
		zap.L().Error("jwt.ParseToken Failed", zap.Error(err))
		consts.ResponseError(ctx, &v1.MessageActionResp{
			ResponseData: consts.ResponseErrorData(consts.CodeInvalidToken),
		})
		return
	}

	// 业务处理
	out, err := service.Chat().GetMessageList(&model.GetMessageListInput{
		UserId:   myClaims.Id,
		ToUserId: req.ToUserId,
	})
	if err != nil {
		zap.L().Error("MessageList strconv.ParseINt Failed", zap.Error(err))
		consts.ResponseError(ctx, v1.GetMessageListResp{
			ResponseData: consts.ResponseErrorData(consts.CodeServerBusy),
		})
		return
	}

	// 返回响应
	consts.ResponseSuccess(ctx, v1.GetMessageListResp{
		ResponseData: consts.ResponseSuccessData("获取消息列表成功"),
		MessageList:  out,
	})
}

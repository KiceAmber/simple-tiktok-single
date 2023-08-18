package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	v1 "simple_tiktok_rime/api/v1"
	"simple_tiktok_rime/internal/consts"
	"simple_tiktok_rime/internal/model"
	"simple_tiktok_rime/internal/service"
	"simple_tiktok_rime/pkg/jwt"
	"strconv"
)

// FollowAction 关注操作
func FollowAction(ctx *gin.Context) {
	var req = new(v1.FollowActionReq)

	// 接收参数
	req.Token = ctx.Query("token")
	req.ActionType = ctx.Query("action_type")

	// 接收 ToUserId 参数
	toUserId, err := strconv.ParseInt(ctx.Query("to_user_id"), 10, 64)
	if err != nil {
		zap.L().Error("FollowAction Parse ToUserId Failed", zap.Error(err))
		consts.ResponseError(ctx, v1.FollowActionResp{
			ResponseData: consts.ResponseErrorData(consts.CodeInvalidParam),
		})
		return
	}
	req.ToUserId = toUserId

	myClaims, err := jwt.ParseToken(req.Token)
	if err != nil {
		zap.L().Error("jwt.ParseToken Failed", zap.Error(err))
		consts.ResponseError(ctx, &v1.FollowActionResp{
			ResponseData: consts.ResponseErrorData(consts.CodeInvalidToken),
		})
		return
	}

	// 业务处理
	_, err = service.Follow().FollowAction(&model.FollowActionInput{
		ActionType: req.ActionType,
		ToUserId:   req.ToUserId,
		UserId:     myClaims.Id,
	})
	if err != nil {
		zap.L().Error("Service FollowAction Failed", zap.Error(err))
		if errors.Is(err, consts.ErrUserFollowedTargetUser) {
			consts.ResponseError(ctx, v1.FollowActionResp{
				ResponseData: consts.ResponseErrorData(consts.CodeUserFollowedTargetUser),
			})
			return
		} else if errors.Is(err, consts.ErrUserNotFollowTargetUser) {
			consts.ResponseError(ctx, v1.FollowActionResp{
				ResponseData: consts.ResponseErrorData(consts.CodeUserNotFollowTargetUser),
			})
			return
		}
		consts.ResponseError(ctx, v1.FollowActionResp{
			ResponseData: consts.ResponseErrorData(consts.CodeServerBusy),
		})
		return
	}

	// 返回响应
	if req.ActionType == "1" {
		consts.ResponseError(ctx, v1.FollowActionResp{
			ResponseData: consts.ResponseSuccessData("关注成功"),
		})
		return
	}
	consts.ResponseError(ctx, v1.FollowActionResp{
		ResponseData: consts.ResponseSuccessData("取消关注成功"),
	})
}

// GetFollowList 获取关注列表
func GetFollowList(ctx *gin.Context) {
	var req = new(v1.GetFollowListReq)

	// 接收参数
	req.Token = ctx.Query("token")

	// 接收 UserId 参数
	userId, err := strconv.ParseInt(ctx.Query("user_id"), 10, 64)
	if err != nil {
		zap.L().Error("GetFollowList Parse UserId Failed", zap.Error(err))
		consts.ResponseError(ctx, v1.GetFollowListResp{
			ResponseData: consts.ResponseErrorData(consts.CodeInvalidParam),
		})
		return
	}
	req.UserId = userId

	_, err = jwt.ParseToken(req.Token)
	if err != nil {
		zap.L().Error("jwt.ParseToken Failed", zap.Error(err))
		consts.ResponseError(ctx, &v1.GetFollowListResp{
			ResponseData: consts.ResponseErrorData(consts.CodeInvalidToken),
		})
		return
	}

	// 业务处理
	out, err := service.Follow().GetFollowList(&model.GetFollowListInput{
		UserId: userId,
	})
	if err != nil {
		zap.L().Error("service.Follow().GetFollowList Failed", zap.Error(err))
		consts.ResponseError(ctx, v1.GetFollowListResp{
			ResponseData: consts.ResponseErrorData(consts.CodeServerBusy),
		})
		return
	}

	// 返回响应
	consts.ResponseError(ctx, v1.GetFollowListResp{
		ResponseData: consts.ResponseSuccessData("查询用户关注列表成功"),
		FollowList:   out.UserList,
	})
}

// GetFollowerList 获取粉丝列表
func GetFollowerList(ctx *gin.Context) {
	var req = new(v1.GetFollowerListReq)

	// 接收参数
	req.Token = ctx.Query("token")

	// 接收 UserId 参数
	userId, err := strconv.ParseInt(ctx.Query("user_id"), 10, 64)
	if err != nil {
		zap.L().Error("GetFollowerList Parse UserId Failed", zap.Error(err))
		consts.ResponseError(ctx, v1.GetFollowerListResp{
			ResponseData: consts.ResponseErrorData(consts.CodeInvalidParam),
		})
		return
	}
	req.UserId = userId

	_, err = jwt.ParseToken(req.Token)
	if err != nil {
		zap.L().Error("jwt.ParseToken Failed", zap.Error(err))
		consts.ResponseError(ctx, &v1.GetFollowerListResp{
			ResponseData: consts.ResponseErrorData(consts.CodeInvalidToken),
		})
		return
	}

	// 业务处理
	out, err := service.Follow().GetFollowerList(&model.GetFollowerListInput{
		UserId: userId,
	})
	if err != nil {
		zap.L().Error("service.Follow().GetFollowList Failed", zap.Error(err))
		consts.ResponseError(ctx, v1.GetFollowerListResp{
			ResponseData: consts.ResponseErrorData(consts.CodeServerBusy),
		})
		return
	}

	// 返回响应
	consts.ResponseError(ctx, v1.GetFollowerListResp{
		ResponseData: consts.ResponseSuccessData("查询用户粉丝列表成功"),
		FollowList:   out.UserList,
	})
}

// GetFriendList 获取好友列表
func GetFriendList(ctx *gin.Context) {
	var req = new(v1.GetFriendListReq)

	// 接收参数
	req.Token = ctx.Query("token")

	// 接收 UserId 参数
	userId, err := strconv.ParseInt(ctx.Query("user_id"), 10, 64)
	if err != nil {
		zap.L().Error("GetFriendList Parse UserId Failed", zap.Error(err))
		consts.ResponseError(ctx, v1.GetFriendListResp{
			ResponseData: consts.ResponseErrorData(consts.CodeInvalidParam),
		})
		return
	}
	req.UserId = userId

	_, err = jwt.ParseToken(req.Token)
	if err != nil {
		zap.L().Error("jwt.ParseToken Failed", zap.Error(err))
		consts.ResponseError(ctx, &v1.GetFriendListResp{
			ResponseData: consts.ResponseErrorData(consts.CodeInvalidToken),
		})
		return
	}

	// 业务处理
	out, err := service.Follow().GetFriendList(&model.GetFriendListInput{
		UserId: userId,
	})
	if err != nil {
		zap.L().Error("service.Follow().GetFriendList Failed", zap.Error(err))
		consts.ResponseError(ctx, v1.GetFriendListResp{
			ResponseData: consts.ResponseErrorData(consts.CodeServerBusy),
		})
		return
	}

	// 返回响应
	consts.ResponseError(ctx, v1.GetFriendListResp{
		ResponseData: consts.ResponseSuccessData("查询用户好友列表成功"),
		FriendList:   out.UserList,
	})
}

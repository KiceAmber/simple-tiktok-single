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

// FavoriteAction 点赞操作
func FavoriteAction(ctx *gin.Context) {
	var req = new(v1.FavoriteActionReq)

	// 接收参数
	req.Token = ctx.Query("token")
	req.ActionType = ctx.Query("action_type")

	// 接收 VideoId 参数
	videoId, err := strconv.ParseInt(ctx.Query("video_id"), 10, 64)
	if err != nil {
		zap.L().Error("CommentAction Parse VideoId Failed", zap.Error(err))
		consts.ResponseError(ctx, v1.FavoriteActionResp{
			ResponseData: consts.ResponseErrorData(consts.CodeInvalidParam),
		})
		return
	}
	req.VideoId = videoId

	myClaims, err := jwt.ParseToken(req.Token)
	if err != nil {
		zap.L().Error("jwt.ParseToken Failed", zap.Error(err))
		consts.ResponseError(ctx, &v1.FavoriteActionResp{
			ResponseData: consts.ResponseErrorData(consts.CodeInvalidToken),
		})
		return
	}

	// 业务处理
	_, err = service.Favorite().FavoriteAction(&model.FavoriteActionInput{
		ActionType: req.ActionType,
		VideoId:    req.VideoId,
		UserId:     myClaims.Id,
	})
	if err != nil {
		zap.L().Error("Service CommentAction Failed", zap.Error(err))
		if errors.Is(err, consts.ErrUserFavoritedVideo) {
			consts.ResponseError(ctx, v1.FavoriteActionResp{
				ResponseData: consts.ResponseErrorData(consts.CodeUserFavoritedVideo),
			})
			return
		} else if errors.Is(err, consts.ErrUserNotFavoriteVideo) {
			consts.ResponseError(ctx, v1.FavoriteActionResp{
				ResponseData: consts.ResponseErrorData(consts.CodeUserNotFavoriteVideo),
			})
			return
		}
		consts.ResponseError(ctx, v1.FavoriteActionResp{
			ResponseData: consts.ResponseErrorData(consts.CodeServerBusy),
		})
		return
	}

	// 返回响应
	if req.ActionType == "1" {
		consts.ResponseError(ctx, v1.FavoriteActionResp{
			ResponseData: consts.ResponseSuccessData("点赞成功"),
		})
		return
	}
	consts.ResponseError(ctx, v1.FavoriteActionResp{
		ResponseData: consts.ResponseSuccessData("取消点赞成功"),
	})
}

// GetUserFavoriteList 点赞列表(喜欢列表)
func GetUserFavoriteList(ctx *gin.Context) {
	var req = new(v1.GetFavoriteVideoListReq)

	// 接收参数
	req.Token = ctx.Query("token")

	// 接收 VideoId 参数
	userId, err := strconv.ParseInt(ctx.Query("user_id"), 10, 64)
	if err != nil {
		zap.L().Error("GetUserFavoriteList Parse UserId Failed", zap.Error(err))
		consts.ResponseError(ctx, v1.GetFavoriteVideoListResp{
			ResponseData: consts.ResponseErrorData(consts.CodeInvalidParam),
		})
		return
	}
	req.UserId = userId

	_, err = jwt.ParseToken(req.Token)
	if err != nil {
		zap.L().Error("jwt.ParseToken Failed", zap.Error(err))
		consts.ResponseError(ctx, &v1.GetFavoriteVideoListResp{
			ResponseData: consts.ResponseErrorData(consts.CodeInvalidToken),
		})
		return
	}

	// 业务处理
	out, err := service.Favorite().GetFavoriteVideoList(&model.GetFavoriteVideoListInput{
		UserId: userId,
	})
	if err != nil {
		zap.L().Error("service.Favorite().GetFavoriteVideoList Failed", zap.Error(err))
		consts.ResponseError(ctx, v1.GetFavoriteVideoListResp{
			ResponseData: consts.ResponseErrorData(consts.CodeServerBusy),
		})
		return
	}

	// 返回响应
	consts.ResponseError(ctx, v1.GetFavoriteVideoListResp{
		ResponseData: consts.ResponseSuccessData("查询用户点赞视频列表成功"),
		VideoList:    out.VideoList,
	})
}

package controller

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	v1 "simple_tiktok_single/api/v1"
	"simple_tiktok_single/internal/consts"
	"simple_tiktok_single/internal/model"
	"simple_tiktok_single/internal/service"
	"simple_tiktok_single/pkg/jwt"
	"strconv"
)

// CommentAction 评论操作，包含添加评论和删除评论
func CommentAction(ctx *gin.Context) {
	var req = new(v1.CommentActionReq)

	// 接收参数
	req.Token = ctx.Query("token")
	req.ActionType = ctx.Query("action_type")
	req.CommentText = ctx.DefaultQuery("comment_text", "")

	// 接收 VideoId 参数
	videoId, err := strconv.ParseInt(ctx.Query("video_id"), 10, 64)
	if err != nil {
		zap.L().Error("CommentAction Parse VideoId Failed", zap.Error(err))
		consts.ResponseError(ctx, v1.CommentActionResp{
			ResponseData: consts.ResponseErrorData(consts.CodeInvalidParam),
		})
		return
	}
	req.VideoId = videoId

	// 接收 CommentId 参数，用于删除评论操作，添加评论不需要
	// 如果接收到的值为 -1，则表示要添加评论，因为没有传递这个参数
	if req.ActionType != "1" {
		commentId, err := strconv.ParseInt(ctx.DefaultQuery("comment_id", "-1"), 10, 64)
		if err != nil {
			zap.L().Error("CommentAction Parse CommentId Failed", zap.Error(err))
			consts.ResponseError(ctx, v1.CommentActionResp{
				ResponseData: consts.ResponseErrorData(consts.CodeInvalidParam),
			})
			return
		}
		req.CommentId = commentId
	}
	myClaims, err := jwt.ParseToken(req.Token)
	if err != nil {
		zap.L().Error("jwt.ParseToken Failed", zap.Error(err))
		consts.ResponseError(ctx, &v1.CommentActionResp{
			ResponseData: consts.ResponseErrorData(consts.CodeInvalidToken),
		})
		return
	}

	// 业务处理
	out, err := service.Comment().CommentAction(&model.CommentActionInput{
		ActionType: req.ActionType,
		Content:    req.CommentText,
		AuthorId:   myClaims.Id,
		VideoId:    req.VideoId,
	})
	if err != nil {
		zap.L().Error("Service CommentAction Failed", zap.Error(err))
		consts.ResponseError(ctx, v1.CommentActionResp{
			ResponseData: consts.ResponseErrorData(consts.CodeServerBusy),
		})
		return
	}

	// 返回响应
	if req.ActionType == "1" {
		consts.ResponseError(ctx, v1.CommentActionResp{
			ResponseData: consts.ResponseSuccessData("添加评论成功"),
			Comment:      out.CommentItem,
		})
		return
	}
	consts.ResponseError(ctx, v1.CommentActionResp{
		ResponseData: consts.ResponseSuccessData("删除评论成功"),
	})
}

// GetCommentList 获取视频的列表
func GetCommentList(ctx *gin.Context) {
	var req = new(v1.GetCommentListReq)

	// 接收参数
	req.Token = ctx.Query("token")

	videoId, err := strconv.ParseInt(ctx.Query("video_id"), 10, 64)
	if err != nil {
		zap.L().Error("GetCommentList Parse VideoId Failed", zap.Error(err))
		consts.ResponseError(ctx, v1.GetCommentListResp{
			ResponseData: consts.ResponseErrorData(consts.CodeInvalidParam),
		})
		return
	}
	req.VideoId = videoId

	// 业务处理
	out, err := service.Comment().GetCommentList(&model.GetCommentListInput{VideoId: req.VideoId})
	if err != nil {
		zap.L().Error("Service GetCommentList Failed", zap.Error(err))
		consts.ResponseError(ctx, v1.GetCommentListResp{
			ResponseData: consts.ResponseErrorData(consts.CodeServerBusy),
		})
		return
	}

	// 返回响应
	consts.ResponseSuccess(ctx, v1.GetCommentListResp{
		ResponseData: consts.ResponseSuccessData("获取视频评论列表成功"),
		CommentList:  out.CommentList,
	})
}

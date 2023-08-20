package controller

import (
	v1 "simple_tiktok_rime/api/v1"
	"simple_tiktok_rime/internal/consts"
	"simple_tiktok_rime/internal/model"
	"simple_tiktok_rime/internal/service"
	"simple_tiktok_rime/pkg/jwt"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// PublishVideo 视频投稿操作
func PublishVideo(ctx *gin.Context) {

	var req = new(v1.PublishVideoReq)

	// 接收参数
	fileData, err := ctx.FormFile("data")
	if err != nil {
		zap.L().Error("PublishVideo ctx.FormFile Failed", zap.Error(err))
		consts.ResponseError(ctx, &v1.PublishVideoResp{
			ResponseData: consts.ResponseErrorData(consts.CodeInvalidFileHeader),
		})
		return
	}

	req.Data = fileData
	req.Token = ctx.PostForm("token")
	req.Title = ctx.PostForm("title")

	// 解析 Token，检测 Token 是否合法
	myClaims, err := jwt.ParseToken(req.Token)
	if err != nil {
		zap.L().Error("jwt.ParseToken Failed", zap.Error(err))
		consts.ResponseError(ctx, &v1.PublishVideoResp{
			ResponseData: consts.ResponseErrorData(consts.CodeInvalidToken),
		})
		return
	}

	// 业务处理
	if err = service.Video().PublishVideo(&model.PublishVideoInput{
		AuthorId:  myClaims.Id,
		Title:     req.Title,
		VideoData: req.Data,
	}); err != nil {
		zap.L().Error("service.Video().PublishVideo Failed", zap.Error(err))
		consts.ResponseError(ctx, &v1.PublishVideoResp{
			ResponseData: consts.ResponseErrorData(consts.CodeServerBusy),
		})
		return
	}

	// 返回响应
	consts.ResponseSuccess(ctx, &v1.PublishVideoResp{
		ResponseData: consts.ResponseSuccessData("发布视频成功"),
	})
}

// VideoFeed 视频流
func VideoFeed(ctx *gin.Context) {
	var req = new(v1.VideoFeedReq)
	var userId int64 = -1 // -1 表示未登录用户

	// 接收参数，这里要注意时间问题
	timeString := ctx.DefaultQuery("latest_time", "-1")
	if timeString == "-1" || timeString == "" {
		latestTime := time.Now().Unix()
		req.LatestTime = time.Unix(latestTime, 0) // 将 int64 转化为 time.Time
	} else {
		if timeString == "0" {
			// 如果等于 0 说明是刷新了视频，则要将时间戳更新为上一次视频列表的最后一个视频所处的时间戳
			// 所以这里需要一个全局时间戳来存储上一次视频流的最后一个视频的时间戳
			req.LatestTime = time.Unix(consts.NextTimeStamp, 0)
		} else {
			latestTime, err := strconv.ParseInt(timeString, 10, 64)
			if err != nil {
				zap.L().Error("Parse LatestTime Failed", zap.Error(err))
				consts.ResponseError(ctx, &v1.VideoFeedResp{
					ResponseData: consts.ResponseErrorData(consts.CodeInvalidTimeStamp),
				})
				return
			}
			if latestTime > time.Now().Unix() {
				latestTime = time.Now().Unix()
			}
			req.LatestTime = time.Unix(latestTime, 0) // 将 int64 转化为 time.Time
		}
	}

	req.Token = ctx.DefaultQuery("token", "")
	if req.Token != "" {
		// 解析 Token，检测 Token 是否合法
		myClaims, err := jwt.ParseToken(req.Token)
		if err != nil {
			zap.L().Error("jwt.ParseToken Failed", zap.Error(err))
			consts.ResponseError(ctx, &v1.PublishVideoResp{
				ResponseData: consts.ResponseErrorData(consts.CodeInvalidToken),
			})
			return
		}
		userId = myClaims.Id
	}

	// 业务处理
	out, err := service.Video().GetVideoFeed(&model.VideoFeedInput{
		LatestTime: req.LatestTime,
		UserId:     userId,
	})
	if err != nil {
		zap.L().Error("service.Video().GetVideoFeed Failed", zap.Error(err))
		consts.ResponseError(ctx, &v1.VideoFeedResp{
			ResponseData: consts.ResponseErrorData(consts.CodeServerBusy),
		})
		return
	}

	// 返回响应
	consts.NextTimeStamp = out.NextTime
	consts.ResponseSuccess(ctx, &v1.VideoFeedResp{
		ResponseData: consts.ResponseSuccessData("获取视频流成功"),
		NextTime:     out.NextTime,
		VideoList:    out.VideoList,
	})
}

// GetVideoPublishedList 获取视频发布列表
func GetVideoPublishedList(ctx *gin.Context) {
	var req = new(v1.GetVidePublishedListReq)

	// 接收参数，这里要注意时间问题
	userId, err := strconv.ParseInt(ctx.Query("user_id"), 10, 64)
	if err != nil {
		zap.L().Error("service.Video().GetVideoFeed Failed", zap.Error(err))
		consts.ResponseError(ctx, &v1.GetVidePublishedListResp{
			ResponseData: consts.ResponseErrorData(consts.CodeInvalidParam),
		})
		return
	}
	req.UserId = userId
	req.Token = ctx.DefaultQuery("token", "")

	// 业务处理
	out, err := service.Video().GetVideoPublishedList(&model.GetVideoPublishedListInput{
		UserId: req.UserId,
	})
	if err != nil {
		zap.L().Error("service.Video().GetVideoPublishedList Failed", zap.Error(err))
		consts.ResponseError(ctx, &v1.GetVidePublishedListResp{
			ResponseData: consts.ResponseErrorData(consts.CodeServerBusy),
		})
		return
	}

	// 返回响应
	consts.ResponseSuccess(ctx, &v1.GetVidePublishedListResp{
		ResponseData: consts.ResponseSuccessData("获取视频发布列表成功"),
		VideoList:    out.VideoList,
	})
}

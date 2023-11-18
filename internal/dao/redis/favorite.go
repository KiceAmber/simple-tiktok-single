package redis

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"simple_tiktok_single/internal/dao/mysql"
	"simple_tiktok_single/internal/model"
	"strconv"
)

type dFavorite struct {
}

var (
	favorite *dFavorite
)

func Favorite() *dFavorite {
	if favorite == nil {
		once.Do(func() {
			favorite = &dFavorite{}
		})
	}
	return favorite
}

// IsUserFavorite 判断用户是否点赞了视频
func (*dFavorite) IsUserFavorite(userId int64, videoId int64) (bool, error) {

	userVideoKey := fmt.Sprintf("user:%d:favorite_video", userId)
	return rdb.SIsMember(context.Background(), userVideoKey, videoId).Result()
}

// VideoFavoriteAction 添加或减少视频的点赞数
func (*dFavorite) VideoFavoriteAction(videoId int64, increment int) error {

	videoFavoriteSet := "video_favorite"
	videoKey := fmt.Sprintf("video:%d", videoId)

	return rdb.ZIncrBy(context.Background(), videoFavoriteSet, float64(increment), videoKey).Err()
}

// UserFavoriteActionToVideo 用户与视频的点赞关系设置
func (*dFavorite) UserFavoriteActionToVideo(actionType string, userId int64, videoId int64) error {

	userVideoKey := fmt.Sprintf("user:%d:favorite_video", userId)

	if actionType == "1" {
		return rdb.SAdd(context.Background(), userVideoKey, videoId).Err()
	} else {
		return rdb.SRem(context.Background(), userVideoKey, videoId).Err()
	}
}

// GetFavoriteList 获取点赞(喜欢)列表
func (*dFavorite) GetFavoriteList(in *model.GetFavoriteVideoListInput) (*model.GetFavoriteVideoListOutput, error) {

	var out = &model.GetFavoriteVideoListOutput{
		FavoriteList: nil,
		VideoList:    []*model.VideoItem{},
	}

	videoIdList, err := rdb.SMembers(context.Background(), fmt.Sprintf("user:%d:favorite_video", in.UserId)).Result()
	if err != nil {
		return nil, err
	}

	for _, videoIdString := range videoIdList {
		videoId, err := strconv.ParseInt(videoIdString, 10, 64)
		if err != nil {
			zap.L().Error("Redis.Favorite().GetFavoriteList Failed", zap.Error(err))
			continue
		}
		video, err := mysql.Video().GetVideoInfoByVideoId(videoId)
		if err != nil {
			zap.L().Error("mysql.Video().GetVideoInfoByVideoId() Failed", zap.Error(err))
			continue
		}
		out.VideoList = append(out.VideoList, &model.VideoItem{
			Id:            video.Id,
			AuthorId:      video.AuthorId,
			FavoriteCount: video.FavoriteCount,
			CommentCount:  video.CommentCount,
			Title:         video.Title,
			PlayUrl:       video.PlayUrl,
			CoverUrl:      video.CoverUrl,
			IsFavorite:    video.IsFavorite,
			CreatedAt:     video.CreatedAt,
		})
	}
	return out, nil
}

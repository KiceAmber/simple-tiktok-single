package redis

import (
	"context"
	"fmt"
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

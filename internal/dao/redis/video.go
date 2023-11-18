package redis

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"simple_tiktok_single/internal/model"
)

type dVideo struct {
}

var (
	video *dVideo
)

func Video() *dVideo {
	if video == nil {
		once.Do(func() {
			video = &dVideo{}
		})
	}
	return video
}

// IsUserFavoriteVideo 用户是否点赞该视频
func (*dVideo) IsUserFavoriteVideo(userId int64, videoList []*model.VideoItem) error {

	for _, video := range videoList {
		exists, err := rdb.SIsMember(
			context.Background(),
			fmt.Sprintf("user:%d:favorite_video", userId),
			video.Id,
		).Result()
		if err != nil {
			zap.L().Error("rdb.SIsMember Failed", zap.Error(err))
			continue
		}

		// 如果存在则用户点赞了该视频，设置 is_favorite = true
		if exists {
			video.IsFavorite = true
		}
	}
	return nil
}

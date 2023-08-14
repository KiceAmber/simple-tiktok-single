package cron

import (
	"context"
	"errors"
	"fmt"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"simple_tiktok_rime/internal/dao/mysql"
	"simple_tiktok_rime/internal/dao/redis"
	"simple_tiktok_rime/internal/model/entity"
	"simple_tiktok_rime/pkg/snowflake"
	"strconv"
	"strings"
)

// Init 初始化定时任务
func Init() error {
	// 初始化两个定时器，一个40秒为间隔，一个1分钟为间隔
	crontabBy40S := cron.New(cron.WithSeconds())
	if _, err := crontabBy40S.AddFunc("*/40 * * * * ?", CronTabBy40S); err != nil {
		zap.L().Error("crontabBy40S.AddFunc Failed", zap.Error(err))
		return err
	}

	cronTabBy1M := cron.New()
	if _, err := cronTabBy1M.AddFunc("*/1 * * * ?", CronTabBy1M); err != nil {
		zap.L().Error("cronTabBy1M.AddFunc Failed", zap.Error(err))
		return err
	}

	crontabBy40S.Start()
	cronTabBy1M.Start()

	return nil
}

// 40 秒定时任务执行的操作
func CronTabBy40S() {
	if err := syncFavoriteCountForVideo(); err != nil {
		zap.L().Error("syncFavoriteCountForVideo Failed", zap.Error(err))
		return
	}
	if err := syncUserFavoriteVideo(); err != nil {
		zap.L().Error("syncUserFavoriteVideo Failed", zap.Error(err))
		return
	}
}

// 1 分钟定时任务执行的操作
func CronTabBy1M() {
	if err := syncUserTotalFavoriteCount(); err != nil {
		zap.L().Error("syncUserTotalFavoriteCount failed", zap.Error(err))
		return
	}
}

// 同步视频点赞数量
func syncFavoriteCountForVideo() error {

	var (
		rdb = redis.New()
		mdb = mysql.New()
	)
	ctx := context.Background()

	videoFavoriteSet := "video_favorite"
	// 首先对比 redis 的数据是否和 mysql 的数据一致

	videoSet, err := rdb.ZRange(ctx, videoFavoriteSet, 0, -1).Result()
	// videoSet: [video:1002 video:1003 video:1001]
	if err != nil {
		return err
	}
	videoIdArr := []int64{}
	// 解析出所有的视频ID
	for _, video := range videoSet {
		videoId, _ := strconv.ParseInt(strings.Split(video, ":")[1], 10, 64)
		videoIdArr = append(videoIdArr, videoId)
	}
	// videoIdArr: [1002 1003 1001]

	for _, videoId := range videoIdArr {
		// 根据 videoId 依次查询 mysql 中的点赞数量
		var video = &entity.Video{}
		mdb.Where("id = ?", videoId).First(video)

		result, err := rdb.ZScore(ctx, videoFavoriteSet, fmt.Sprintf("video:%d", videoId)).Result()
		if err != nil {
			zap.L().Error("rdb.ZScore Failed", zap.Error(err))
			continue
		}
		// 如果不一致，则同步数据
		if video.FavoriteCount != int64(result) {
			err := mdb.Model(&entity.Video{}).Where("id = ?", videoId).Update("favorite_count", result).Error
			if err != nil {
				zap.L().Error("mdb update favorite_count failed", zap.Error(err))
				continue
			}
		}
	}
	return nil
}

// 同步用户与视频的点赞关系
func syncUserFavoriteVideo() error {

	var (
		rdb = redis.New()
		mdb = mysql.New()
	)
	ctx := context.Background()
	// 查询出 mysql 所有的用户
	userList := []*entity.User{}
	err := mdb.Find(&userList).Error
	if err != nil {
		return err
	}

	// 遍历 userList,针对于每一个 user 进行比较
	for _, user := range userList {
		favoriteVideoList, err := rdb.SMembers(ctx, fmt.Sprintf("user:%d:favorite_video", user.Id)).Result()
		// favoriteVideoList: [1001 1002 1003]
		if err != nil {
			zap.L().Error("rdb.SMembers", zap.Error(err))
			continue
		}

		for _, favoriteVideo := range favoriteVideoList {
			favoriteVideoId, err := strconv.ParseInt(favoriteVideo, 10, 64)
			if err != nil {
				zap.L().Error("strconv.ParseInt favoriteVideo Failed", zap.Error(err))
				continue
			}
			// 在 mysql 查询 favorite 表是否有这样的对应关系
			var favorite = &entity.Favorite{}
			err = mdb.Where("user_id = ? AND video_id = ?", user.Id, favoriteVideoId).First(favorite).Error
			if err != nil {
				// 如果找不到这样的数据，则表示没有这样的对应关系，就直接插入 favorite 数据
				if errors.Is(err, gorm.ErrRecordNotFound) {
					favorite = &entity.Favorite{
						Id:      snowflake.GenID(),
						UserId:  user.Id,
						VideoId: favoriteVideoId,
					}
					err = mdb.Create(favorite).Error
					if err != nil {
						zap.L().Error("mdb Create Favorite Failed", zap.Error(err))
						continue
					}
				}
				zap.L().Error("rdb Select Favorite Failed", zap.Error(err))
				continue
			}
			// 如果没有报 gorm.ErrRecordNotFound 说明找到了数据，则不用插入新的数据
		}
	}
	return nil
}

// 同步用户的总被点赞量
func syncUserTotalFavoriteCount() error {

	var mdb = mysql.New()

	userList := []*entity.User{}
	err := mdb.Find(&userList).Error
	if err != nil {
		return err
	}

	// 查询每一个用户各个视频点赞量，
	for _, user := range userList {
		var totalFavoriteCount int64 = 0
		var videoList = []*entity.Video{}
		err := mdb.Where("author_id = ?", user.Id).Find(&videoList).Error
		if err != nil {
			zap.L().Error("Query user VideoLists Failed", zap.Error(err))
			continue
		}

		// 计算视频点赞量总和
		for _, video := range videoList {
			totalFavoriteCount += video.FavoriteCount
		}

		// 最后将总的点赞数量更新到用户表中
		err = mdb.Model(&entity.User{}).Where("id = ?", user.Id).Update("total_favorited", totalFavoriteCount).Error
		if err != nil {
			zap.L().Error("Update user TotalFavoriteCount Failed", zap.Error(err))
			continue
		}
	}

	return nil
}

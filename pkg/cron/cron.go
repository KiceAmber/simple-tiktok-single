package cron

import (
	"context"
	"errors"
	"fmt"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"simple_tiktok_single/internal/dao/mysql"
	"simple_tiktok_single/internal/dao/redis"
	"simple_tiktok_single/internal/model/entity"
	"simple_tiktok_single/pkg/snowflake"
	"strconv"
	"strings"
)

// Init 初始化定时任务
func Init() error {
	// 初始化两个定时器，一个40秒为间隔，一个1分钟为间隔
	crontabBy5S := cron.New(cron.WithSeconds())
	if _, err := crontabBy5S.AddFunc("*/40 * * * * ?", CronTabBy5S); err != nil {
		zap.L().Error("crontabBy40S.AddFunc Failed", zap.Error(err))
		return err
	}

	cronTabBy8S := cron.New()
	if _, err := cronTabBy8S.AddFunc("*/1 * * * ?", CronTabBy8S); err != nil {
		zap.L().Error("cronTabBy1M.AddFunc Failed", zap.Error(err))
		return err
	}

	go crontabBy5S.Start()
	go cronTabBy8S.Start()

	return nil
}

// CronTabBy5S 5 秒定时任务执行的操作
func CronTabBy5S() {
	if err := syncFavoriteCountForVideo(); err != nil {
		zap.L().Error("syncFavoriteCountForVideo Failed", zap.Error(err))
		return
	}
	if err := addUserFavoriteVideo(); err != nil {
		zap.L().Error("syncUserFavoriteVideo Failed", zap.Error(err))
		return
	}

	if err := syncUserFollowerCount(); err != nil {
		zap.L().Error("syncUserFollowerCount Failed", zap.Error(err))
		return
	}
	if err := syncUserFollowCount(); err != nil {
		zap.L().Error("syncUserFollowCount Failed", zap.Error(err))
		return
	}
	if err := addUserFollowRelation(); err != nil {
		zap.L().Error("syncUserFollowRelation failed", zap.Error(err))
		return
	}
}

// CronTabBy8S 1 分钟定时任务执行的操作
func CronTabBy8S() {
	if err := syncUserTotalFavoriteCount(); err != nil {
		zap.L().Error("syncUserTotalFavoriteCount failed", zap.Error(err))
		return
	}

	if err := delUserFavoriteVideo(); err != nil {
		zap.L().Error("delUserFavoriteVideo failed", zap.Error(err))
		return
	}

	if err := delUserFollowRelation(); err != nil {
		zap.L().Error("delUserFollowRelation failed", zap.Error(err))
		return
	}
}

// ================= 点赞模块定时任务 =================

// syncFavoriteCountForVideo 同步视频点赞数量到 mysql
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
			zap.L().Error("rdb.ZScore videoFavoriteSet Failed", zap.Error(err))
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

// syncUserTotalFavoriteCount 同步用户的总被点赞量到 mysql
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

// addUserFavoriteVideo 添加用户与视频的点赞关系到 mysql
func addUserFavoriteVideo() error {

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
			zap.L().Error("rdb.SMembers favoriteVideoList", zap.Error(err))
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

// delUserFollowRelation 删除用户与视频的点赞关系到 mysql
func delUserFavoriteVideo() error {

	// 查找所有用户
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

	// 查询 mysql 用户与视频的点赞关系
	for _, user := range userList {
		userVideoList := []*entity.Favorite{}
		err := mdb.Where("user_id = ?", user.Id).Find(&userVideoList).Error
		if err != nil {
			zap.L().Error("delUserFavoriteVideo Query userVideoList Failed", zap.Error(err))
			continue
		}
		// 对比在 redis 是否存在
		for _, video := range userVideoList {
			exists, err := rdb.SIsMember(ctx, fmt.Sprintf("user:%d:favorite_video", user.Id), video.VideoId).Result()
			if err != nil {
				zap.L().Error("delUserFavoriteVideo rdb.SIsMember Failed", zap.Error(err))
				continue
			}

			// 如果不存在就删除
			if !exists {
				mdb.Delete(video)
			}
		}
	}
	return nil
}

// ================= 关注模块定时任务 =================

// syncUserFollowerCount 同步用户的粉丝数到 mysql
func syncUserFollowerCount() error {

	var (
		rdb = redis.New()
		mdb = mysql.New()
	)
	ctx := context.Background()
	userFollowerSet := "user_follower"
	// 从 redis 查询粉丝数量
	userSet, err := rdb.ZRange(ctx, userFollowerSet, 0, -1).Result()
	if err != nil {
		return err
	}

	for _, user := range userSet {
		userId, _ := strconv.ParseInt(strings.Split(user, ":")[1], 10, 64)
		result, err := rdb.ZScore(ctx, userFollowerSet, fmt.Sprintf("user:%d", userId)).Result()
		if err != nil {
			zap.L().Error("rdb.ZScore userFollowerSet Failed", zap.Error(err))
			continue
		}
		// 插入数据库
		err = mdb.Model(&entity.User{}).Where("id = ?", userId).Update("follower_count", result).Error
		if err != nil {
			zap.L().Error("mdb update follower_count failed", zap.Error(err))
			continue
		}
	}
	return nil
}

// syncUserFollowCount 同步用户的关注数到 mysql
func syncUserFollowCount() error {

	var (
		rdb = redis.New()
		mdb = mysql.New()
	)
	ctx := context.Background()
	userFollowSet := "user_follow"

	// 从 redis 查询粉丝数量
	userSet, err := rdb.ZRange(ctx, userFollowSet, 0, -1).Result()
	if err != nil {
		return err
	}

	for _, user := range userSet {
		userId, _ := strconv.ParseInt(strings.Split(user, ":")[1], 10, 64)
		result, err := rdb.ZScore(ctx, userFollowSet, fmt.Sprintf("user:%d", userId)).Result()
		if err != nil {
			zap.L().Error("rdb.ZScore userFollowSet Failed", zap.Error(err))
			continue
		}
		// 插入数据库
		err = mdb.Model(&entity.User{}).Where("id = ?", userId).Update("follow_count", result).Error
		if err != nil {
			zap.L().Error("mdb update follow_count failed", zap.Error(err))
			continue
		}
	}
	return nil
}

// addUserFollowRelation 添加用户之间的关注关系到 mysql
func addUserFollowRelation() error {

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
		followerList, err := rdb.SMembers(ctx, fmt.Sprintf("user:%d:follower_user", user.Id)).Result()
		// followerList: [1001 1002 1003]
		if err != nil {
			zap.L().Error("rdb.SMembers followerList", zap.Error(err))
			continue
		}

		for _, follower := range followerList {
			followerId, err := strconv.ParseInt(follower, 10, 64)
			if err != nil {
				zap.L().Error("strconv.ParseInt followerId Failed", zap.Error(err))
				continue
			}
			// 在 mysql 查询 follow 表是否有这样的对应关系
			var follow = &entity.Follow{}
			err = mdb.Where("user_id = ? AND follower_id = ?", user.Id, followerId).First(follow).Error
			if err != nil {
				// 如果找不到这样的数据，则表示没有这样的对应关系，就直接插入 follow 数据
				if errors.Is(err, gorm.ErrRecordNotFound) {
					newFollow := &entity.Follow{
						Id:         snowflake.GenID(),
						UserId:     user.Id,
						FollowerId: followerId,
					}
					err = mdb.Create(newFollow).Error
					if err != nil {
						zap.L().Error("mdb Create Follow Failed", zap.Error(err))
						continue
					}
				} else {
					zap.L().Error("rdb Select Follow Failed", zap.Error(err))
					continue
				}
			}
			// 如果没有报 gorm.ErrRecordNotFound 说明找到了数据，则不用插入新的数据
		}
	}
	return nil
}

// delUserFollowRelation 删除用户之间的关注关系到 mysql
func delUserFollowRelation() error {
	// 查找所有用户
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

	// 查询 mysql 用户之间的关注关系
	for _, user := range userList {
		followList := []*entity.Follow{}
		err := mdb.Where("user_id = ?", user.Id).Find(&followList).Error
		if err != nil {
			zap.L().Error("delUserFollowRelation Query followList Failed", zap.Error(err))
			continue
		}
		// 对比在 redis 是否存在
		for _, follow := range followList {
			exists, err := rdb.SIsMember(ctx, fmt.Sprintf("user:%d:follower_user", user.Id), follow.FollowerId).Result()
			if err != nil {
				zap.L().Error("delUserFollowRelation rdb.SIsMember Failed", zap.Error(err))
				continue
			}

			// 如果不存在就删除
			if !exists {
				mdb.Delete(follow)
			}
		}
	}
	return nil
}

package favorite

import (
	"simple_tiktok_single/internal/consts"
	"simple_tiktok_single/internal/dao/mysql"
	"simple_tiktok_single/internal/dao/redis"
	"simple_tiktok_single/internal/model"
	"simple_tiktok_single/internal/service"
)

type sFavorite struct{}

func init() {
	service.RegisterFavorite(New())
}

func New() *sFavorite {
	return &sFavorite{}
}

// FavoriteAction 赞操作
func (*sFavorite) FavoriteAction(in *model.FavoriteActionInput) (out *model.FavoriteActionOutput, err error) {

	if in.ActionType == "1" {
		return addFavorite(in)
	} else {
		return cancelFavorite(in)
	}
}

// addFavorite 点赞
func addFavorite(in *model.FavoriteActionInput) (out *model.FavoriteActionOutput, err error) {
	// 首先判断用户是否已经对视频点赞了
	exists, err := redis.Favorite().IsUserFavorite(in.UserId, in.VideoId)
	if err != nil {
		return nil, err
	}
	if exists { // 如果已经点赞了，就报错
		return nil, consts.ErrUserFavoritedVideo
	}

	// 点赞视频，点赞数量 +1
	if err = redis.Favorite().VideoFavoriteAction(in.VideoId, 1); err != nil {
		return nil, err
	}

	// 添加用户与视频的点赞关系
	if err := redis.Favorite().UserFavoriteActionToVideo(in.ActionType, in.UserId, in.VideoId); err != nil {
		return nil, err
	}
	return
}

// cancelFavorite 取消点赞
func cancelFavorite(in *model.FavoriteActionInput) (out *model.FavoriteActionOutput, err error) {

	// 判断用户是否已经对视频点赞
	exists, err := redis.Favorite().IsUserFavorite(in.UserId, in.VideoId)
	if err != nil {
		return nil, err
	}
	// 未点赞，取消点赞操作报错
	if !exists {
		return nil, consts.ErrUserNotFavoriteVideo
	}

	// 取消点赞，点赞数 -1
	if err = redis.Favorite().VideoFavoriteAction(in.VideoId, -1); err != nil {
		return nil, err
	}

	// 取消用户与视频的点赞关系
	if err := redis.Favorite().UserFavoriteActionToVideo(in.ActionType, in.UserId, in.VideoId); err != nil {
		return nil, err
	}
	return
}

// GetFavoriteVideoList 获取点赞的视频列表
func (*sFavorite) GetFavoriteVideoList(in *model.GetFavoriteVideoListInput) (out *model.GetFavoriteVideoListOutput, err error) {

	// 首先从 redis 查询，如果 redis 出现问题才从 mysql 查询
	out, err = getFavoriteVideoListByRedis(in)
	if err == nil {
		return
	}

	// 如果有 err 那么说明 redis 查询失败了，这时候从 mysql 中查询
	out, err = getFavoriteVideoListByMysql(in)
	return
}

// getFavoriteVideoListByRedis 由 redis 查询点赞的视频列表
func getFavoriteVideoListByRedis(in *model.GetFavoriteVideoListInput) (out *model.GetFavoriteVideoListOutput, err error) {

	out, err = redis.Favorite().GetFavoriteList(in)
	if err != nil {
		return nil, err
	}

	// 然后查询视频的作者信息
	for _, video := range out.VideoList {
		userOut, err := mysql.User().QueryUserById(&model.GetUserInfoInput{UserId: video.AuthorId})
		if err != nil {
			continue
		}
		video.Author = userOut.UserItem
	}
	return
}

// getFavoriteVideoListByMysql 由 mysql 查询点赞的视频列表
func getFavoriteVideoListByMysql(in *model.GetFavoriteVideoListInput) (out *model.GetFavoriteVideoListOutput, err error) {
	// 获取到点赞的视频ID列表
	out, err = mysql.Favorite().GetFavoriteList(in)
	if err != nil {
		return nil, err
	}

	// 再通过视频ID 获取到视频的信息，这需要调用视频服务
	for _, favorite := range out.FavoriteList {
		video, err := mysql.Video().GetVideoInfoByVideoId(favorite.VideoId)
		if err != nil {
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

	// 再通过用户ID查询用户信息
	for _, video := range out.VideoList {
		userOut, err := mysql.User().QueryUserById(&model.GetUserInfoInput{UserId: video.AuthorId})
		if err != nil {
			continue
		}
		video.Author = userOut.UserItem
	}
	return
}

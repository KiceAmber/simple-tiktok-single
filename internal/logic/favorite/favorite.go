package favorite

import (
	"simple_tiktok_rime/internal/consts"
	"simple_tiktok_rime/internal/dao/mysql"
	"simple_tiktok_rime/internal/dao/redis"
	"simple_tiktok_rime/internal/model"
	"simple_tiktok_rime/internal/service"
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
		return AddFavorite(in)
	} else {
		return CancelFavorite(in)
	}
}

// AddFavorite 点赞
func AddFavorite(in *model.FavoriteActionInput) (out *model.FavoriteActionOutput, err error) {
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

// CancelFavorite 取消点赞
func CancelFavorite(in *model.FavoriteActionInput) (out *model.FavoriteActionOutput, err error) {

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

// 获取点赞的视频列表
func (*sFavorite) GetFavoriteVideoList(in *model.GetFavoriteVideoListInput) (out *model.GetFavoriteVideoListOutput, err error) {

	// 获取到点赞的视频ID列表
	out, err = mysql.Favorite().GetFavoriteList(in)
	if err != nil {
		return nil, err
	}

	// 再通过视频ID 获取到视频的信息，这需要调用视频服务
	for _, favorite := range out.FavoriteList {
		videoList, err := mysql.Video().GetVideoListByVideoId(favorite.VideoId)
		if err != nil {
			continue
		}
		for _, video := range videoList {
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

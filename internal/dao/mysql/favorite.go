package mysql

import (
	"simple_tiktok_rime/internal/model"
	"simple_tiktok_rime/internal/model/entity"
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

// GetFavoriteList 获取到点赞列表
func (*dFavorite) GetFavoriteList(in *model.GetFavoriteVideoListInput) (*model.GetFavoriteVideoListOutput, error) {

	var out = &model.GetFavoriteVideoListOutput{
		FavoriteList: []*model.FavoriteItem{},
		VideoList:    []*model.VideoItem{},
	}

	// 这里获取到用户点赞的视频ID
	var favoriteList = []*entity.Favorite{}
	result := engine.Where("user_id = ?", in.UserId).Find(&favoriteList)
	if result.Error != nil {
		return nil, result.Error
	}

	for _, favorite := range favoriteList {
		out.FavoriteList = append(out.FavoriteList, &model.FavoriteItem{
			UserId:  favorite.UserId,
			VideoId: favorite.VideoId,
		})
	}
	return out, nil
}

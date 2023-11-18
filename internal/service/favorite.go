package service

import "simple_tiktok_single/internal/model"

type IFavorite interface {
	FavoriteAction(in *model.FavoriteActionInput) (out *model.FavoriteActionOutput, err error)
	GetFavoriteVideoList(in *model.GetFavoriteVideoListInput) (out *model.GetFavoriteVideoListOutput, err error)
}

var (
	localFavorite IFavorite
)

func Favorite() IFavorite {
	if localFavorite == nil {
		panic("implement not found for interface IFavorite, forgot register?")
	}
	return localFavorite
}

func RegisterFavorite(i IFavorite) {
	localFavorite = i
}

package model

// FavoriteActionInput 点赞操作 Input
type FavoriteActionInput struct {
	ActionType string `json:"action_type"`
	VideoId    int64  `json:"video_id"`
	UserId     int64  `json:"user_id"`
}

// FavoriteActionOutput 点赞操作 Output
type FavoriteActionOutput struct {
}

// GetFavoriteVideoListInput 获取点赞视频列表 Input
type GetFavoriteVideoListInput struct {
	UserId int64 `json:"user_id"`
}

// GetFavoriteVideoListOutput 获取点赞列表 Output
type GetFavoriteVideoListOutput struct {
	FavoriteList []*FavoriteItem `json:"favorite_list"`
	VideoList    []*VideoItem    `json:"video_list"`
}

type FavoriteItem struct {
	UserId  int64 `json:"user_id"`
	VideoId int64 `json:"video_id"`
}

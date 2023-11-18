package mysql

import (
	"gorm.io/gorm"
	"simple_tiktok_single/internal/model"
	"simple_tiktok_single/internal/model/entity"
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

// InsertVideoInfo 插入视频信息
func (*dVideo) InsertVideoInfo(in *model.PublishVideoInput) error {

	engine.Transaction(func(tx *gorm.DB) error {
		// 插入视频数据
		newVideo := &entity.Video{
			Id:       in.Id,
			AuthorId: in.AuthorId,
			Title:    in.Title,
			PlayUrl:  in.PlayUrl,
			CoverUrl: in.CoverUrl,
		}
		if err := engine.Create(newVideo).Error; err != nil {
			return err
		}

		// 同步更新用户表的作品数量
		user := &entity.User{}
		tx.Where("id = ?", in.AuthorId).First(user)
		if err := tx.Model(&entity.User{}).Where("id = ?", in.AuthorId).Update("work_count", user.WorkCount+1).Error; err != nil {
			return err
		}
		return nil
	})

	return nil
}

// GetVideoList 获取视频流数据
func (*dVideo) GetVideoList(in *model.VideoFeedInput) (*model.VideoFeedOutput, error) {
	var videoList = make([]*entity.Video, 0, 30) // 一次性最多获取 30 个视频
	var out = &model.VideoFeedOutput{
		VideoList: make([]*model.VideoItem, 0, 30),
	}
	engine.Where("created_at < ?", in.LatestTime).Limit(30).Find(&videoList)

	for _, video := range videoList {
		var videoItem = &model.VideoItem{
			Id:            video.Id,
			AuthorId:      video.AuthorId,
			FavoriteCount: video.FavoriteCount,
			CommentCount:  video.CommentCount,
			Title:         video.Title,
			PlayUrl:       video.PlayUrl,
			CoverUrl:      video.CoverUrl,
			Author:        nil,
			CreatedAt:     video.CreatedAt,
		}
		out.VideoList = append(out.VideoList, videoItem)
	}

	return out, nil
}

// GetVideoPublishedList 获取用户视频发布列表
func (*dVideo) GetVideoPublishedList(in *model.GetVideoPublishedListInput) (*model.GetVideoPublishedListOutput, error) {
	var videoList = []*entity.Video{}
	var out = &model.GetVideoPublishedListOutput{
		VideoList: []*model.VideoItem{},
	}
	// 查询出视频
	engine.Where("author_id = ?", in.UserId).Find(&videoList)

	for _, video := range videoList {
		var videoItem = &model.VideoItem{
			Id:            video.Id,
			AuthorId:      video.AuthorId,
			FavoriteCount: video.FavoriteCount,
			CommentCount:  video.CommentCount,
			Title:         video.Title,
			PlayUrl:       video.PlayUrl,
			CoverUrl:      video.CoverUrl,
			IsFavorite:    video.IsFavorite,
			Author:        nil,
			CreatedAt:     video.CreatedAt,
		}
		out.VideoList = append(out.VideoList, videoItem)
	}

	return out, nil
}

// GetVideoInfoByVideoId 根据视频ID获取视频信息
func (*dVideo) GetVideoInfoByVideoId(videoId int64) (*entity.Video, error) {

	var video = &entity.Video{}

	result := engine.Where("id = ?", videoId).Find(video)
	if result.Error != nil {
		return nil, result.Error
	}

	return video, nil
}

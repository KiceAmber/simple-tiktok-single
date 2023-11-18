package service

import "simple_tiktok_single/internal/model"

type IVideo interface {
	PublishVideo(in *model.PublishVideoInput) (err error)
	GetVideoFeed(in *model.VideoFeedInput) (out *model.VideoFeedOutput, err error)
	GetVideoPublishedList(in *model.GetVideoPublishedListInput) (out *model.GetVideoPublishedListOutput, err error)
}

var (
	localVideo IVideo
)

func Video() IVideo {
	if localVideo == nil {
		panic("implement not found for interface IVideo, forgot register?")
	}
	return localVideo
}

func RegisterVideo(i IVideo) {
	localVideo = i
}

package service

import "simple_tiktok_rime/internal/model"

type IFollow interface {
	FollowAction(in *model.FollowActionInput) (out *model.FollowActionOutput, err error)
	GetFollowList(in *model.GetFollowListInput) (out *model.GetFollowListOutput, err error)
	GetFollowerList(in *model.GetFollowerListInput) (out *model.GetFollowerListOutput, err error)
}

var (
	localFollow IFollow
)

func Follow() IFollow {
	if localFollow == nil {
		panic("implement not found for interface IFollow, forgot register?")
	}
	return localFollow
}

func RegisterFollow(i IFollow) {
	localFollow = i
}

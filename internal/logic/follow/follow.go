package follow

import (
	"simple_tiktok_single/internal/consts"
	"simple_tiktok_single/internal/dao/mysql"
	"simple_tiktok_single/internal/dao/redis"
	"simple_tiktok_single/internal/model"
	"simple_tiktok_single/internal/service"
)

type sFollow struct{}

func init() {
	service.RegisterFollow(New())
}

func New() *sFollow {
	return &sFollow{}
}

// FollowAction 关注操作，根据 action_type 里决定是关注还是取消关注
func (*sFollow) FollowAction(in *model.FollowActionInput) (out *model.FollowActionOutput, err error) {
	if in.ActionType == "1" {
		return AddFollow(in)
	} else {
		return CancelFollow(in)
	}
}

// AddFollow 关注操作
func AddFollow(in *model.FollowActionInput) (out *model.FollowActionOutput, err error) {

	// 首先判断是否已经关注 ToUserId 用户
	exists, err := redis.Follow().IsUserFollow(in.UserId, in.ToUserId)
	if err != nil {
		return nil, err
	}
	if exists { // 已关注，则关注操作报错
		return nil, consts.ErrUserFollowedTargetUser
	}

	// 关注用户，UserId 用户的关注数 +1
	if err = redis.Follow().UserFollowAction(in.UserId, 1); err != nil {
		return nil, err
	}

	// 关注用户，ToUserId 用户的粉丝数 +1
	if err = redis.Follow().UserFollowerAction(in.ToUserId, 1); err != nil {
		return nil, err
	}

	// 添加用户与用户之间的关注关系
	if err := redis.Follow().UserFollowActionToTargetUser(in.ActionType, in.UserId, in.ToUserId); err != nil {
		return nil, err
	}
	return
}

// CancelFollow 取消关注
func CancelFollow(in *model.FollowActionInput) (out *model.FollowActionOutput, err error) {

	// 判断是否已经对该用户关注
	exists, err := redis.Follow().IsUserFollow(in.UserId, in.ToUserId)
	if err != nil {
		return nil, err
	}
	// 未关注，取消关注操作报错
	if !exists {
		return nil, consts.ErrUserNotFollowTargetUser
	}

	// 取消关注，UserId 用户的关注数 -1
	if err = redis.Follow().UserFollowAction(in.UserId, -1); err != nil {
		return nil, err
	}

	// 取消关注，关注的该用户粉丝数 -1
	if err = redis.Follow().UserFollowerAction(in.ToUserId, -1); err != nil {
		return nil, err
	}

	// 取消用户与该用户的粉丝关系
	if err := redis.Follow().UserFollowActionToTargetUser(in.ActionType, in.UserId, in.ToUserId); err != nil {
		return nil, err
	}
	return
}

// GetFollowList 获取关注列表
func (*sFollow) GetFollowList(in *model.GetFollowListInput) (out *model.GetFollowListOutput, err error) {

	// 获取到 UserId 的关注列表
	out, err = mysql.Follow().GetFollowList(in)
	if err != nil {
		return nil, err
	}
	return
}

// GetFollowerList 获取粉丝列表
func (*sFollow) GetFollowerList(in *model.GetFollowerListInput) (out *model.GetFollowerListOutput, err error) {
	// 获取到 UserId 的粉丝列表
	out, err = mysql.Follow().GetFollowerList(in)
	return
}

// GetFriendList 获取用户好友列表
func (*sFollow) GetFriendList(in *model.GetFriendListInput) (out *model.GetFriendListOutput, err error) {
	out, err = mysql.Follow().GetFriendList(in)
	return
}

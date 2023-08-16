package redis

import (
	"context"
	"fmt"
)

type dFollow struct {
}

var (
	follow *dFollow
)

func Follow() *dFollow {
	if follow == nil {
		once.Do(func() {
			follow = &dFollow{}
		})
	}
	return follow
}

// IsUserFollow 判断是否已经关注了该用户
func (*dFollow) IsUserFollow(userId int64, targetUserId int64) (bool, error) {

	userKey := fmt.Sprintf("user:%d:follow_user", userId)

	return rdb.SIsMember(context.Background(), userKey, targetUserId).Result()
}

// UserFollowAction 添加或减少用户的关注数
func (*dFollow) UserFollowAction(userId int64, increment int) error {

	userFollowerSet := "user_follow"
	userKey := fmt.Sprintf("user:%d", userId)

	return rdb.ZIncrBy(context.Background(), userFollowerSet, float64(increment), userKey).Err()
}

// UserFollowerAction 添加或减少用户的粉丝数
func (*dFollow) UserFollowerAction(targetUserId int64, increment int) error {

	userFollowerSet := "user_follower"
	userKey := fmt.Sprintf("user:%d", targetUserId)

	return rdb.ZIncrBy(context.Background(), userFollowerSet, float64(increment), userKey).Err()
}

// UserFollowActionToTargetUser 用户与目标用户建立关注关系
func (*dFollow) UserFollowActionToTargetUser(actionType string, userId int64, targetUserId int64) error {

	userFollowKey := fmt.Sprintf("user:%d:follow_user", userId)           // userId 用户的关注集合
	userFollowerKey := fmt.Sprintf("user:%d:follower_user", targetUserId) // targetUserId 用户的粉丝集合

	if actionType == "1" {
		if err := rdb.SAdd(context.Background(), userFollowKey, targetUserId).Err(); err != nil {
			return err
		}

		if err := rdb.SAdd(context.Background(), userFollowerKey, userId).Err(); err != nil {
			return err
		}
	} else {
		if err := rdb.SRem(context.Background(), userFollowKey, targetUserId).Err(); err != nil {
			return err
		}

		if err := rdb.SRem(context.Background(), userFollowerKey, userId).Err(); err != nil {
			return err
		}
	}
	return nil
}

package mysql

import (
	"go.uber.org/zap"
	"simple_tiktok_rime/internal/model"
	"simple_tiktok_rime/internal/model/entity"
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

// GetFollowList 获取关注列表
func (*dFollow) GetFollowList(in *model.GetFollowListInput) (*model.GetFollowListOutput, error) {

	out := &model.GetFollowListOutput{
		UserList: []*model.UserItem{},
	}

	// 通过用户的ID查询对应的关注的人
	followList := []*entity.Follow{}
	err := engine.Where("follower_id = ?", in.UserId).Find(&followList).Error
	if err != nil {
		return nil, err
	}

	// 根据 followId 查询用户信息
	for _, follow := range followList {
		user := &entity.User{}
		err := engine.Where("id = ?", follow.UserId).First(user).Error
		if err != nil {
			zap.L().Error("Query Follow UserInfo Failed", zap.Error(err))
			continue
		}
		out.UserList = append(out.UserList, &model.UserItem{
			Id:              user.Id,
			FollowCount:     user.FollowerCount,
			FollowerCount:   user.FollowerCount,
			TotalFavorited:  user.TotalFavorited,
			WorkCount:       user.WorkCount,
			FavoriteCount:   user.FavoriteCount,
			Name:            user.Name,
			Avatar:          user.Avatar,
			BackgroundImage: user.BackgroundImage,
			Signature:       user.Signature,
			IsFollow:        true,
		})
	}

	return out, nil
}

// GetFollowerList 获取关注列表
func (*dFollow) GetFollowerList(in *model.GetFollowerListInput) (*model.GetFollowerListOutput, error) {

	out := &model.GetFollowerListOutput{
		UserList: []*model.UserItem{},
	}

	// 先查找关注的人的ID
	followList := []*entity.Follow{}
	err := engine.Where("user_id = ?", in.UserId).Find(&followList).Error
	if err != nil {
		return nil, err
	}

	// 根据 followId 查询用户信息
	for _, follow := range followList {
		user := &entity.User{}
		err := engine.Where("id = ?", follow.FollowerId).First(user).Error
		if err != nil {
			zap.L().Error("Query Follow UserInfo Failed", zap.Error(err))
			continue
		}
		out.UserList = append(out.UserList, &model.UserItem{
			Id:              user.Id,
			FollowCount:     user.FollowerCount,
			FollowerCount:   user.FollowerCount,
			TotalFavorited:  user.TotalFavorited,
			WorkCount:       user.WorkCount,
			FavoriteCount:   user.FavoriteCount,
			Name:            user.Name,
			Avatar:          user.Avatar,
			BackgroundImage: user.BackgroundImage,
			Signature:       user.Signature,
			IsFollow:        true,
		})
	}

	return out, nil
}

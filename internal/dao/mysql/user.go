package mysql

import (
	"errors"
	"gorm.io/gorm"
	"simple_tiktok_rime/internal/consts"
	"simple_tiktok_rime/internal/model"
	"simple_tiktok_rime/internal/model/entity"
	"strings"
)

type dUser struct {
}

var (
	user *dUser
)

func User() *dUser {
	if user == nil {
		once.Do(func() {
			user = &dUser{}
		})
	}
	return user
}

// QueryUserByName 根据用户名查找用户数据
func (*dUser) QueryUserByName(Username string) (*entity.User, error) {

	user := &entity.User{Name: Username}

	result := engine.Take(user)
	if result.Error != nil {
		// 如果未找到用户数据，返回自定义的错误 ErrUserNotExists
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, consts.ErrUserNotExists
		}

		return nil, result.Error
	}
	return user, nil
}

// InsertUserInfo 插入用户信息
func (*dUser) InsertUserInfo(in *model.UserRegisterInput) (*model.UserRegisterOutput, error) {

	newUser := &entity.User{
		Id:              in.Id,
		Name:            in.Username,
		Password:        in.Password,
		Avatar:          consts.DefaultAvatar,
		BackgroundImage: consts.DefaultBackgroundImage,
		Signature:       consts.DefaultSignature,
	}

	result := engine.Create(newUser)
	if result.Error != nil {
		if strings.Contains(result.Error.Error(), "Duplicate entry") {
			return nil, consts.ErrUserExists
		}
		return nil, result.Error
	}

	return &model.UserRegisterOutput{
		Id: newUser.Id,
	}, nil
}

// QueryUserByNameAndPwd 根据用户名和密码查询用户
func (*dUser) QueryUserByNameAndPwd(in *model.UserLoginInput) (*model.UserLoginOutput, error) {

	user := &entity.User{}

	result := engine.Where("name = ? AND password = ?", in.Username, in.Password).First(user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, consts.ErrUserNotExists
		}
		return nil, result.Error
	}

	return &model.UserLoginOutput{
		Id: user.Id,
	}, nil
}

// QueryUserById 根据用户 ID 查询用户
func (*dUser) QueryUserById(in *model.GetUserInfoInput) (*model.GetUserInfoOutput, error) {

	user := &entity.User{}

	result := engine.Where("id = ?", in.UserId).First(user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, consts.ErrUserNotExists
		}
		return nil, result.Error
	}

	return &model.GetUserInfoOutput{
		UserItem: &model.UserItem{
			Id:              user.Id,
			FollowCount:     user.FollowCount,
			FollowerCount:   user.FollowerCount,
			TotalFavorited:  user.TotalFavorited,
			WorkCount:       user.WorkCount,
			FavoriteCount:   user.FavoriteCount,
			Name:            user.Name,
			Avatar:          user.Avatar,
			BackgroundImage: user.BackgroundImage,
			Signature:       user.Signature,
			IsFollow:        user.IsFollow,
		},
	}, nil
}

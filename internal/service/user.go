package service

import "simple_tiktok_rime/internal/model"

type IUser interface {
	UserRegister(in *model.UserRegisterInput) (out *model.UserRegisterOutput, err error)
	UserLogin(in *model.UserLoginInput) (out *model.UserLoginOutput, err error)
	GetUserInfo(in *model.GetUserInfoInput) (out *model.GetUserInfoOutput, err error)
}

var (
	localUser IUser
)

func User() IUser {
	if localUser == nil {
		panic("implement not found for interface IUser, forgot register?")
	}
	return localUser
}

func RegisterUser(i IUser) {
	localUser = i
}

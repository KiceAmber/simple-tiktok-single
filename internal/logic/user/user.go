package user

import (
	"errors"
	"simple_tiktok_rime/internal/consts"
	"simple_tiktok_rime/internal/dao/mysql"
	"simple_tiktok_rime/internal/model"
	"simple_tiktok_rime/internal/service"
	"simple_tiktok_rime/pkg/jwt"
	"simple_tiktok_rime/pkg/snowflake"
	"simple_tiktok_rime/pkg/toolx"
)

type sUser struct{}

func init() {
	service.RegisterUser(New())
}

func New() *sUser {
	return &sUser{}
}

func (*sUser) UserRegister(in *model.UserRegisterInput) (out *model.UserRegisterOutput, err error) {
	// 首先查询用户名是否存在
	_, err = mysql.User().QueryUserByName(in.Username)
	if err != nil {
		if !errors.Is(err, consts.ErrUserNotExists) {
			return nil, err
		}
	}

	// 密码加密
	in.Password = toolx.Md5Encrypt(in.Password)

	// 生成随机用户 ID
	userId := snowflake.GenID()
	in.Id = userId

	// 数据库入库操作
	out, err = mysql.User().InsertUserInfo(in)
	if err != nil {
		return nil, err
	}

	// 生成 Token
	token, err := jwt.GenToken(out.Id)
	if err != nil {
		return nil, err
	}

	out.Token = token
	return
}

func (*sUser) UserLogin(in *model.UserLoginInput) (out *model.UserLoginOutput, err error) {

	in.Password = toolx.Md5Encrypt(in.Password)
	out, err = mysql.User().QueryUserByNameAndPwd(in)
	if err != nil {
		return nil, err
	}

	// 生成 Token
	token, err := jwt.GenToken(out.Id)
	if err != nil {
		return nil, err
	}

	out.Token = token
	return
}

func (*sUser) GetUserInfo(in *model.GetUserInfoInput) (out *model.GetUserInfoOutput, err error) {

	out, err = mysql.User().QueryUserById(in)
	return
}

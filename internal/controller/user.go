package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	v1 "simple_tiktok_rime/api/v1"
	"simple_tiktok_rime/internal/consts"
	"simple_tiktok_rime/internal/model"
	"simple_tiktok_rime/internal/service"
	"simple_tiktok_rime/pkg/jwt"
	"strconv"
)

// UserRegister 用户注册
func UserRegister(ctx *gin.Context) {
	var req = new(v1.UserRegisterReq)

	// 绑定参数
	req.Username = ctx.Query("username")
	req.Password = ctx.Query("password")

	// 业务操作
	out, err := service.User().UserRegister(&model.UserRegisterInput{
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		zap.L().Error("service.User().UserRegister Failed", zap.Error(err))
		if errors.Is(err, consts.ErrUserExists) {
			consts.ResponseError(ctx, &v1.UserRegisterResp{
				ResponseData: consts.ResponseErrorData(consts.CodeUserExists),
			})
			return
		}
		consts.ResponseError(ctx, &v1.UserRegisterResp{
			ResponseData: consts.ResponseErrorData(consts.CodeServerBusy),
		})
		return
	}

	// 返回成功响应
	consts.ResponseSuccess(ctx, &v1.UserRegisterResp{
		ResponseData: consts.ResponseSuccessData("注册成功"),
		UserId:       out.Id,
		Token:        out.Token,
	})
}

// UserLogin 用户登录
func UserLogin(ctx *gin.Context) {
	var req = new(v1.UserLoginReq)

	// 绑定参数
	req.Username = ctx.Query("username")
	req.Password = ctx.Query("password")

	// 业务操作
	out, err := service.User().UserLogin(&model.UserLoginInput{
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		zap.L().Error("service.User().UserLogin Failed", zap.Error(err))
		if errors.Is(err, consts.ErrUserNotExists) {
			consts.ResponseError(ctx, &v1.UserLoginResp{
				ResponseData: consts.ResponseErrorData(consts.CodeLoginFailed),
			})
			return
		}
		consts.ResponseError(ctx, &v1.UserLoginResp{
			ResponseData: consts.ResponseErrorData(consts.CodeServerBusy),
		})
		return
	}

	consts.ResponseSuccess(ctx, &v1.UserLoginResp{
		ResponseData: consts.ResponseSuccessData("登录成功"),
		UserId:       out.Id,
		Token:        out.Token,
	})
}

// GetUserInfo 获取用户主页个人信息
func GetUserInfo(ctx *gin.Context) {

	var req = new(v1.GetUserInfoReq)

	// 接收参数
	userId, err := strconv.ParseInt(ctx.Query("user_id"), 10, 64)
	if err != nil {
		zap.L().Error("GetUserInfo Parse Param Failed", zap.Error(err))
		consts.ResponseError(ctx, v1.GetUserInfoResp{
			ResponseData: consts.ResponseErrorData(consts.CodeInvalidParam),
		})
		return
	}

	req.Token = ctx.Query("token")
	req.UserId = userId

	// 解析 Token，检测 Token 是否合法
	myClaims, err := jwt.ParseToken(req.Token)
	if err != nil {
		zap.L().Error("jwt.ParseToken Failed", zap.Error(err))
		consts.ResponseError(ctx, v1.GetUserInfoResp{
			ResponseData: consts.ResponseErrorData(consts.CodeInvalidToken),
		})
		return
	}
	ctx.Set(consts.CtxUserIdKey, myClaims.Id) // 将当前的 UserId 保存到请求的上下文中

	// 业务处理
	out, err := service.User().GetUserInfo(&model.GetUserInfoInput{
		UserId: req.UserId,
	})
	if err != nil {
		zap.L().Error("service.User().GetUserInfo Failed", zap.Error(err))
		consts.ResponseError(ctx, v1.GetUserInfoResp{
			ResponseData: consts.ResponseErrorData(consts.CodeServerBusy),
		})
		return
	}

	// 返回响应
	consts.ResponseSuccess(ctx, v1.GetUserInfoResp{
		ResponseData: consts.ResponseSuccessData("获取用户个人信息成功"),
		User:         out.UserItem,
	})
}

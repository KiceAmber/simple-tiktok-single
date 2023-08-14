package middleware

import (
	"github.com/gin-gonic/gin"
	"simple_tiktok_rime/internal/consts"
	"simple_tiktok_rime/pkg/jwt"
	"strings"
)

// JWTAuthMiddleware 基于 JWT认证的中间件
func JWTAuthMiddleware() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		auth := ctx.Request.Header.Get("Authorization")
		if auth == "" {
			consts.ResponseError(ctx, consts.CodeNeedLogin)
			ctx.Abort()
			return
		}

		// 按空格分割，因为格式为 Bearer xxxxx.xxx.xxx
		parts := strings.SplitN(auth, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			consts.ResponseError(ctx, consts.CodeInvalidToken)
			ctx.Abort()
			return
		}

		// parts[1]是获取的 token，对其解析
		myClaims, err := jwt.ParseToken(parts[1])
		if err != nil {
			consts.ResponseError(ctx, consts.CodeInvalidToken)
			ctx.Abort()
			return
		}

		// 将当前的 userID 信息保存到请求的上下文 ctx 上
		ctx.Set(consts.CtxUserIdKey, myClaims.Id)
		ctx.Next() // 后续处理请求的函数中，可以用过 ctx.Get(consts.CtxUserIdKey) 来获取当前用户的信息
	}
}

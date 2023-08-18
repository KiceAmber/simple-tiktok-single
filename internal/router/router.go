package router

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"simple_tiktok_rime/internal/controller"
	"simple_tiktok_rime/internal/middleware"
	"simple_tiktok_rime/logs"
	"simple_tiktok_rime/manifest/config"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Init() *gin.Engine {
	r := gin.New()
	r.Use(logs.GinLogger(), logs.GinRecovery(true), middleware.CORSMiddleware())

	douyin := r.Group("/douyin")
	// 用户模块
	{
		douyin.POST("/user/register/", controller.UserRegister) // 用户注册
		douyin.POST("/user/login/", controller.UserLogin)       // 用户登录
		douyin.GET("/user/", controller.GetUserInfo)            // 获取用户信息
	}

	// 视频模块
	{
		douyin.GET("/feed/", controller.VideoFeed)                     // 视频流
		douyin.POST("/publish/action/", controller.PublishVideo)       // 发布视频投稿
		douyin.GET("/publish/list/", controller.GetVideoPublishedList) // 视频发布列表
	}

	// 评论模块
	{
		douyin.POST("/comment/action/", controller.CommentAction) // 评论操作
		douyin.GET("/comment/list/", controller.GetCommentList)   // 获取评论列表操作
	}

	// 点赞模块
	{
		douyin.POST("/favorite/action/", controller.FavoriteAction)   // 点赞操作
		douyin.GET("/favorite/list/", controller.GetUserFavoriteList) // 点赞列表(喜欢列表)
	}

	// 关注模块
	{
		douyin.POST("/relation/action/", controller.FollowAction)          // 关注操作
		douyin.GET("/relation/follow/list/", controller.GetFollowList)     // 关注列表
		douyin.GET("/relation/follower/list/", controller.GetFollowerList) // 粉丝列表
		douyin.GET("/relation/friend/list/", controller.GetFriendList)     // 好友列表                                                            // 好友列表
	}

	// 聊天模块
	{
		douyin.POST("/message/action/", controller.MessageAction) // 发送消息操作
		douyin.GET("/message/chat/", controller.MessageList)      // 获取消息记录
	}

	// 错误路由
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"msg": "未知页面",
		})
	})

	return r
}

func Setup(r *gin.Engine) {
	// 启动服务(优雅关机)
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", config.Conf.Port),
		Handler: r,
	}

	go func() {
		// 开启一个goroutine启动服务
		fmt.Printf("\nServer Running on port%s\n", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			zap.L().Fatal("listen: ", zap.Error(err))
		}
	}()

	// 等待中断信号来优雅地关闭服务器，为关闭服务器操作设置一个5秒的超时
	quit := make(chan os.Signal, 1) // 创建一个接收信号的通道

	// kill 默认会发送 syscall.SIGTERM 信号
	// kill -2 发送 syscall.SIGINT 信号，我们常用的Ctrl+C就是触发系统SIGINT信号
	// kill -9 发送 syscall.SIGKILL 信号，但是不能被捕获，所以不需要添加它
	// signal.Notify把收到的 syscall.SIGINT或syscall.SIGTERM 信号转发给quit
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM) // 此处不会阻塞
	<-quit                                               // 阻塞在此，当接收到上述两种信号时才会往下执行
	zap.L().Info("Shutdown Server ...")
	// 创建一个5秒超时的context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// 5秒内优雅关闭服务（将未处理完的请求处理完再关闭服务），超过5秒就超时退出
	if err := srv.Shutdown(ctx); err != nil {
		zap.L().Fatal("Server Shutdown:", zap.Error(err))
	}

	zap.L().Info("Server exiting")
}

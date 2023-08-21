package websocket

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// upgrader 用于升级 http 为 websocket 协议
var upgrader = websocket.Upgrader{}

// connSet 存储用户聊天连接
var connSet = make(map[*websocket.Conn]struct{})

// New 升级 HTTP 协议并返回 WebSocket 的连接
func New(ctx *gin.Context) (*websocket.Conn, error) {
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

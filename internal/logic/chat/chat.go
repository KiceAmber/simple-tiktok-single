package chat

import (
	"simple_tiktok_single/internal/dao/mysql"
	"simple_tiktok_single/internal/model"
	"simple_tiktok_single/internal/service"
	"simple_tiktok_single/pkg/snowflake"
)

type sChat struct{}

func init() {
	service.RegisterChat(New())
}

func New() *sChat {
	return &sChat{}
}

// MessageAction 消息操作
func (*sChat) MessageAction(in *model.MessageActionInput) (out *model.MessageActionOutput, err error) {

	in.Id = snowflake.GenID()
	if err = mysql.Chat().InsertChatInfo(in); err != nil {
		return nil, err
	}
	return
}

// GetMessageList 获取消息列表
func (*sChat) GetMessageList(in *model.GetMessageListInput) (out *model.GetMessageListOutput, err error) {

	out, err = mysql.Chat().QueryMessageList(in)
	if err != nil {
		return nil, err
	}
	return
}

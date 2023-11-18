package service

import "simple_tiktok_single/internal/model"

type IChat interface {
	MessageAction(in *model.MessageActionInput) (out *model.MessageActionOutput, err error)
	GetMessageList(in *model.GetMessageListInput) (out *model.GetMessageListOutput, err error)
}

var (
	localChat IChat
)

func Chat() IChat {
	if localChat == nil {
		panic("implement not found for interface IChat, forgot register?")
	}
	return localChat
}

func RegisterChat(i IChat) {
	localChat = i
}

package mysql

import (
	"simple_tiktok_rime/internal/model"
	"simple_tiktok_rime/internal/model/entity"
)

type dChat struct {
}

var (
	chat *dChat
)

func Chat() *dChat {
	if chat == nil {
		once.Do(func() {
			chat = &dChat{}
		})
	}
	return chat
}

// InsertChatInfo 插入聊天数据
func (*dChat) InsertChatInfo(in *model.MessageActionInput) error {

	newChat := &entity.Chat{
		Id:       in.Id,
		UserId:   in.UserId,
		ToUserId: in.ToUserId,
		Content:  in.Content,
	}

	if err := engine.Create(newChat).Error; err != nil {
		return err
	}

	return nil
}

// QueryMessageList 获取聊天信息数据
func (*dChat) QueryMessageList(in *model.GetMessageListInput) (*model.GetMessageListOutput, error) {
	messageList := []*entity.Chat{}
	var out = new(model.GetMessageListOutput)

	if err := engine.Where("user_id = ? AND to_user_id = ?", in.UserId, in.ToUserId).Find(&messageList).Error; err != nil {
		return nil, err
	}

	for _, message := range messageList {
		var messageItem = &model.MessageItem{
			Id:         message.Id,
			FromUserId: message.UserId,
			ToUserId:   message.ToUserId,
			Content:    message.Content,
			CreateTime: message.CreatedAt,
		}
		out.MessageList = append(out.MessageList, messageItem)
	}
	return out, nil
}

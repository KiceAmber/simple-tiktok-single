package model

// MessageActionInput 消息操作 Input
type MessageActionInput struct {
	Id         int64  `json:"id"`
	UserId     int64  `json:"user_id"`
	ToUserId   int64  `json:"to_user_id"`
	Content    string `json:"content"`
	ActionType string `json:"action_type"`
}

// MessageActionOutput 消息操作 Output
type MessageActionOutput struct {
}

// GetMessageListInput 消息列表 Input
type GetMessageListInput struct {
	UserId   int64 `json:"user_id"`
	ToUserId int64 `json:"to_user_id"`
}

// GetMessageListOutput 消息列表 Output
type GetMessageListOutput struct {
	MessageList []*MessageItem `json:"message_list"`
}

// MessageItem 消息单项
type MessageItem struct {
	Id         int64  `json:"id"`
	FromUserId int64  `json:"from_user_id"`
	ToUserId   int64  `json:"to_user_id"`
	CreateTime int64  `json:"create_time"`
	Content    string `json:"content"`
}

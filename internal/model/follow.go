package model

// FollowActionInput 关注操作 Input
type FollowActionInput struct {
	ActionType string `json:"action_type"`
	UserId     int64  `json:"user_id"`
	ToUserId   int64  `json:"to_user_id"`
}

// FollowActionOutput 关注操作 Output
type FollowActionOutput struct {
}

// GetFollowListInput 获取关注列表 Input
type GetFollowListInput struct {
	UserId int64 `json:"user_id"`
}

// GetFollowListOutput 获取关注列表 Output
type GetFollowListOutput struct {
	UserList []*UserItem `json:"user_list"`
}

// GetFollowerListInput 获取粉丝列表 Input
type GetFollowerListInput struct {
	UserId int64 `json:"user_id"`
}

// GetFollowerListOutput 获取粉丝列表 Output
type GetFollowerListOutput struct {
	UserList []*UserItem `json:"user_list"`
}

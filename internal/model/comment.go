package model

// CommentActionInput 评论操作 Input
type CommentActionInput struct {
	ActionType string `json:"action_type"`
	CommentId  int64  `json:"comment_id"`
	Content    string `json:"content"`
	AuthorId   int64  `json:"author_id"`
	VideoId    int64  `json:"video_id"`
}

// CommentActionOutput 评论操作 Output
type CommentActionOutput struct {
	*CommentItem `json:"comment"`
}

type CommentItem struct {
	Id         int64     `json:"id"`
	UserId     int64     `json:"-"`
	Content    string    `json:"content"`
	CreateDate string    `json:"create_date"`
	User       *UserItem `json:"user"`
}

// GetCommentListInput 显示评论列表 Input
type GetCommentListInput struct {
	VideoId int64 `json:"video_id"`
}

// GetCommentListOutput 显示评论列表 Output
type GetCommentListOutput struct {
	CommentList []*CommentItem `json:"comment_list"`
}

package service

import "simple_tiktok_rime/internal/model"

type IComment interface {
	CommentAction(in *model.CommentActionInput) (out *model.CommentActionOutput, err error)
	GetCommentList(in *model.GetCommentListInput) (out *model.GetCommentListOutput, err error)
}

var (
	localComment IComment
)

func Comment() IComment {
	if localComment == nil {
		panic("implement not found for interface IComment, forgot register?")
	}
	return localComment
}

func RegisterComment(i IComment) {
	localComment = i
}

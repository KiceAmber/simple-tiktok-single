package comment

import (
	"simple_tiktok_single/internal/dao/mysql"
	"simple_tiktok_single/internal/model"
	"simple_tiktok_single/internal/service"
	"simple_tiktok_single/pkg/snowflake"
)

type sComment struct{}

func init() {
	service.RegisterComment(New())
}

func New() *sComment {
	return &sComment{}
}

// CommentAction 评论操作
func (*sComment) CommentAction(in *model.CommentActionInput) (out *model.CommentActionOutput, err error) {
	if in.ActionType == "1" {
		return AddComment(in)
	} else {
		return DeleteComment(in)
	}
}

// AddComment 添加评论
func AddComment(in *model.CommentActionInput) (out *model.CommentActionOutput, err error) {

	in.CommentId = snowflake.GenID()

	// 插入评论信息
	out, err = mysql.Comment().InsertCommentInfo(in)
	if err != nil {
		return nil, err
	}

	// 插入成功后，还需要将发送评论的用户也查询出来
	userOut, err := mysql.User().QueryUserById(&model.GetUserInfoInput{
		UserId: in.AuthorId,
	})
	if err != nil {
		return nil, err
	}
	out.User = userOut.UserItem
	return
}

// DeleteComment 删除评论
func DeleteComment(in *model.CommentActionInput) (out *model.CommentActionOutput, err error) {
	err = mysql.Comment().DeleteCommentInfo(in)
	return
}

// GetCommentList 获取视频评论列表
func (*sComment) GetCommentList(in *model.GetCommentListInput) (out *model.GetCommentListOutput, err error) {

	out, err = mysql.Comment().QueryCommentList(in)
	if err != nil {
		return nil, err
	}
	for _, comment := range out.CommentList {
		userOut, err := mysql.User().QueryUserById(&model.GetUserInfoInput{UserId: comment.UserId})
		if err != nil {
			continue
		}
		comment.User = userOut.UserItem
	}
	return
}

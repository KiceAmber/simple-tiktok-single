package mysql

import (
	"gorm.io/gorm"
	"simple_tiktok_rime/internal/model"
	"simple_tiktok_rime/internal/model/entity"
	"time"
)

type dComment struct {
}

var (
	comment *dComment
)

func Comment() *dComment {
	if comment == nil {
		once.Do(func() {
			comment = &dComment{}
		})
	}
	return comment
}

// InsertCommentInfo 插入评论数据
func (*dComment) InsertCommentInfo(in *model.CommentActionInput) (*model.CommentActionOutput, error) {
	newComment := &entity.Comment{
		Id:        in.CommentId,
		Content:   in.Content,
		AuthorId:  in.AuthorId,
		VideoId:   in.VideoId,
		CreatedAt: time.Now(),
	}

	// 使用事务来更新数据
	engine.Transaction(func(tx *gorm.DB) error {
		// 插入评论数据
		if err := tx.Create(newComment).Error; err != nil {
			return err
		}
		// 更新视频表的评论数量
		video := &entity.Video{}
		if err := tx.Where("id = ?", in.VideoId).First(video).Error; err != nil {
			return err
		}
		if err := tx.Model(&entity.Video{}).Where("id = ?", in.VideoId).Update("comment_count", video.CommentCount+1).Error; err != nil {
			return err
		}
		return nil
	})

	return &model.CommentActionOutput{
		CommentItem: &model.CommentItem{
			Id:         in.CommentId,
			Content:    in.Content,
			CreateDate: newComment.CreatedAt,
		},
	}, nil
}

// DeleteCommentInfo 删除评论数据
func (*dComment) DeleteCommentInfo(in *model.CommentActionInput) error {

	engine.Transaction(func(tx *gorm.DB) error {
		// 删除评论
		comment := &entity.Comment{Id: in.CommentId}
		if err := engine.Delete(comment).Error; err != nil {
			return err
		}

		// 同步更新视频的评论数量
		video := &entity.Video{}
		tx.Where("id = ?", in.VideoId).First(video)
		if err := tx.Where("id = ?", in.VideoId).Update("comment_count", video.CommentCount-1).Error; err != nil {
			return err
		}
		return nil
	})
	return nil
}

// QueryCommentList 查询评论列表
func (*dComment) QueryCommentList(in *model.GetCommentListInput) (*model.GetCommentListOutput, error) {

	var commentList = make([]*entity.Comment, 0)

	result := engine.Where("video_id = ?", in.VideoId).Find(&commentList)
	if result.Error != nil {
		return nil, result.Error
	}

	out := &model.GetCommentListOutput{CommentList: make([]*model.CommentItem, 0)}
	for _, comment := range commentList {
		var temp = &model.CommentItem{
			Id:         comment.Id,
			UserId:     comment.AuthorId,
			Content:    comment.Content,
			CreateDate: comment.CreatedAt,
		}
		out.CommentList = append(out.CommentList, temp)
	}

	return out, nil
}

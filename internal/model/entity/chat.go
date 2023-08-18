package entity

import "time"

type Chat struct {
	Id        int64     `json:"id"`                             // 评论ID
	UserId    int64     `json:"author_id"`                      // 评论发布者ID
	ToUserId  int64     `json:"video_id"`                       // 评论所属视频ID
	Content   string    `json:"content"`                        // 评论内容
	CreatedAt time.Time `json:"created_at"`                     // 评论创建时间
	UpdatedAt time.Time `json:"updated_at"`                     // 评论更新时间
	DeletedAt time.Time `json:"deleted_at" gorm:"default:null"` // 评论删除时间
}

func (*Chat) TableName() string {
	return "chat"
}

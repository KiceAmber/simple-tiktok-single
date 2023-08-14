package entity

import "time"

type Video struct {
	Id            int64     `json:"id"`                             // 视频ID
	AuthorId      int64     `json:"author_id"`                      // 视频作者ID
	FavoriteCount int64     `json:"favorite_count"`                 // 视频点赞数
	CommentCount  int64     `json:"comment_count"`                  // 视频评论数
	Title         string    `json:"title"`                          // 视频标题
	PlayUrl       string    `json:"play_url"`                       // 视频播放链接
	CoverUrl      string    `json:"cover_url"`                      // 视频封面链接
	IsFavorite    bool      `json:"is_favorite"`                    // 是否点赞该视频
	CreatedAt     time.Time `json:"created_at"`                     // 视频发布时间
	UpdatedAt     time.Time `json:"updated_at"`                     // 视频更新时间
	DeletedAt     time.Time `json:"deleted_at" gorm:"default:null"` // 视频删除时间
}

func (*Video) TableName() string {
	return "video"
}

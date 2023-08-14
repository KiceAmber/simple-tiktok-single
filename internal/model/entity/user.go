package entity

import "time"

type User struct {
	Id              int64     `json:"id"`                             // 用户ID
	TotalFavorited  int64     `json:"total_favorited"`                // 用户总被点赞数
	FavoriteCount   int64     `json:"favorite_count"`                 // 用户点赞的数量
	WorkCount       int64     `json:"work_count"`                     // 用户作品数量
	FollowCount     int64     `json:"follow_count"`                   // 用户关注数量
	FollowerCount   int64     `json:"follower_count"`                 // 用户粉丝数量
	Name            string    `json:"name"`                           // 用户名
	Password        string    `json:"password"`                       // 用户密码
	Avatar          string    `json:"avatar"`                         // 用户头像
	BackgroundImage string    `json:"background_image"`               // 用户个人主页顶部大图
	Signature       string    `json:"signature"`                      // 用户简介
	IsFollow        bool      `json:"is_follow"`                      // 是否关注该用户
	CreatedAt       time.Time `json:"created_at"`                     // 用户账户创建时间
	UpdatedAt       time.Time `json:"updated_at"`                     // 用户账户更新时间
	DeletedAt       time.Time `json:"deleted_at" gorm:"default:null"` // 用户账户删除时间
}

func (*User) TableName() string {
	return "user"
}

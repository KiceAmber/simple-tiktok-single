package entity

import "time"

type Follow struct {
	Id         int64     `json:"id"`
	UserId     int64     `json:"user_id"`
	FollowerId int64     `json:"follower_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	DeletedAt  time.Time `json:"deleted_at" gorm:"default:null"`
}

func (*Follow) TableName() string {
	return "follow"
}

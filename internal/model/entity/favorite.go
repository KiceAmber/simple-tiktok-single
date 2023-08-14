package entity

import "time"

type Favorite struct {
	Id        int64     `json:"id"`
	UserId    int64     `json:"user_id"`
	VideoId   int64     `json:"video_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`
}

func (*Favorite) TableName() string {
	return "favorite"
}

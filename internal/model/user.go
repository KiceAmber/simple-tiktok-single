package model

type UserRegisterInput struct {
	Id       int64  `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserRegisterOutput struct {
	Id    int64  `json:"id"`
	Token string `json:"token"`
}

type UserLoginInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserLoginOutput struct {
	Id    int64  `json:"id"`
	Token string `json:"token"`
}

type GetUserInfoInput struct {
	UserId int64 `json:"user_id"`
}

type GetUserInfoOutput struct {
	UserItem *UserItem `json:"user_item"`
}

type UserItem struct {
	Id              int64  `json:"id"`
	FollowCount     int64  `json:"follow_count"`
	FollowerCount   int64  `json:"follower_count"`
	WorkCount       int64  `json:"work_count"`
	FavoriteCount   int64  `json:"favorite_count"`
	TotalFavorited  string `json:"total_favorited"`
	Name            string `json:"name"`
	Avatar          string `json:"avatar"`
	BackgroundImage string `json:"background_image"`
	Signature       string `json:"signature"`
	IsFollow        bool   `json:"is_follow"`
}

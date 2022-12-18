package controller

import "aweme_kitex/models"

type (
	VideoRawData models.VideoRawData
	FavouriteRaw models.FavouriteRaw
	UserRawData  models.UserRawData
	RelationRaw  models.RelationRaw
	CommetRaw    models.CommentRaw
)

type Favourite struct {
	Id      string `json:"identity,omitempty"`
	UserId  string `json:"user_id,omitempty"`
	VideoId string `json:"video_id,omitempty"`
}

// video
type Video struct {
	Id             string `json:"id,omitempty"gorm:"column:video_id"` // id
	Author         User   `json:"author"gorm:"column:user_id"`        // author
	PlayUrl        string `json:"play_url,omitempty"`
	CoverUrl       string `json:"cover_url,omitempty"`
	FavouriteCount int64  `json:"favourite_count,omitempty"`
	CommentCount   int64  `json:"comment_count,omitempty"`
	IsFavourite    bool   `json:"is_favourite,omitempty"`
	Title          string `json:"title,omitempty"`
}

type Comment struct {
	Id         string `json:"id,omitempty"gorm:"column:comment_id"`
	UserId     string `json:"user_id"`
	VideoId    string `json:"video_id"`
	Content    string `json:"content,omitempty"`
	CreateDate string `json:"createDate,omitempty"`
}

type User struct {
	UserId        string `json:"identity,omitempty"`
	Name          string `json:"name,omitempty"`
	Password      string `json:"password,omitempty"`
	FollowCount   int64  `json:"follow_count,omitempty"`
	FollowerCount int64  `json:"follower_count,omitempty"`
	IsFollow      bool   `json:"is_follow,omitempty"`
}

// -------------response--------------
type Response struct {
	StatusCode int32  `json:"status_code,omitempty"`
	StatusMsg  string `json:"status_msg,omitempty"`
}

type UserLoginResponse struct {
	Response
	UserId   string `json:"user_id,omitempty"`
	UserName string `json:"user_name,omitempty"`
	Token    string `json:"token,omitempty"`
}

type UserRegisterResponse struct {
	Response
	UserId string `json:"user_id,omitempty"`
	Token  string `json:"token"`
}

type UserResponse struct {
	Response
	User User `json:"user"`
}

type UserListResponse struct {
	Response
	UserList []User `json:"user_list"`
}

type VideoListResponse struct {
	Response
	VideoList []Video `json:"video_list,omitempty"`
}

// comment
type CommentListResponse struct {
	Response
	CommentList []Comment `json:"comment_list,omitempty"`
}

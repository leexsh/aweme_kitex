package models

import (
	"time"

	"gorm.io/gorm"
)

// common struct
type Response struct {
	StatusCode int32  `json:"status_code,omitempty"`
	StatusMsg  string `json:"status_msg,omitempty"`
}

type Favourite struct {
	Identity string `json:"identity,omitempty"`
	UserId   string `json:"user_id,omitempty"`
	VideoId  string `json:"video_id,omitempty"`
}

// video
type Video struct {
	*gorm.Model    `json:",omitempty"`
	Id             string `json:"id,omitempty"gorm:"column:video_id"` // id
	Author         string `json:"author"gorm:"column:user_id"`        // author
	PlayUrl        string `json:"play_url,omitempty"`
	CoverUrl       string `json:"cover_url,omitempty"`
	FavouriteCount int64  `json:"favourite_count,omitempty"`
	CommentCount   int64  `json:"comment_count,omitempty"`
	IsFavourite    bool   `json:"is_favourite,omitempty"gorm:"-"`
}

type Comment struct {
	Id         int64  `json:"id,omitempty"`
	User       User   `json:"user"`
	Content    string `json:"content,omitempty"`
	CreateDate string `json:"create_date,omitempty"`
}

type User struct {
	*gorm.Model   `json:",omitempty"`
	UserId        string `json:"identity,omitempty"`
	Name          string `json:"name,omitempty"`
	Password      string `json:"password,omitempty"`
	FollowCount   int64  `json:"follow_count,omitempty"`
	FollowerCount int64  `json:"follower_count,omitempty"`
	IsFollow      bool   `json:"is_follow,omitempty"`
}

func (u *User) TableName() string {
	return "user"
}

type Users struct {
	Id        int64  `json:"id"`
	Name      string `json:"name"`
	Password  string `json:"password"`
	FanNum    int64  `json:"follow_count,omitempty"`
	FollowNum int64  `json:"follow_count,omitempty"`
}

type Videos struct {
	Id            int64     `json:"id,omitempty"`
	PlayUrl       string    `json:"play_url" json:"play_url,omitempty"`
	CoverUrl      string    `json:"cover_url,omitempty"`
	FavoriteCount int64     `json:"favorite_count,omitempty"`
	CommentCount  int64     `json:"comment_count,omitempty"`
	Create_at     time.Time `json:"create_at"`
	Title         string    `json:"title"`
}

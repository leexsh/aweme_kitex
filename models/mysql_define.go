package models

import (
	"time"
)

// -------------mysql define----------------
type TimeModel struct {
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
	DeletedAt time.Time `gorm:"column:deleted_at"`
}

type VideoRawData struct {
	TimeModel
	VideoId        string `gorm:"column:video_id"`
	UserId         string `gorm:"column:user_id"`
	Title          string `gorm:"column:title"`
	PlayUrl        string `gorm:"column:play_url"`
	CoverUrl       string `gorm:"column:cover_url"`
	FavouriteCount int64  `gorm:"column:favourite_count"`
	CommentCount   int64  `gorm:"column:comment_count"`
}

func (vr *VideoRawData) TableName() string {
	return "video"
}

// user
type UserRawData struct {
	TimeModel
	UserId        string `gorm:"column:user_id"json:"userId,omitempty"`
	Name          string `gorm:"column:name"json:"name,omitempty"`
	Password      string `gorm:"column:password"json:"password,omitempty"`
	Token         string `gorm:"column:token"json:"token,omitempty"`
	FollowCount   int64  `gorm:"column:follow_count"json:"followCount,omitempty"`
	FollowerCount int64  `gorm:"column:follower_count"json:"FollowerCount,omitempty"`
}

func (u2 *UserRawData) TableName() string {
	return "user"
}

// 喜欢
type FavouriteRaw struct {
	Id      string `gorm:"column:favour_id"`
	UserId  string `gorm:"column:user_id"`
	VideoId string `gorm:"column:video_id"`
}

func (f *FavouriteRaw) TableName() string {
	return "favourite"
}

// 关注
type RelationRaw struct {
	Id       string `gorm:"column:relation_id"`
	UserId   string `gorm:"column:user_id"`
	ToUserId string `gorm:"column:to_user_id"`
	Status   int64  `gorm:"column:status"`
}

func (r *RelationRaw) TableName() string {
	return "relation"
}

// 关注
type CommentRaw struct {
	TimeModel
	Id      string `gorm:"column:comment_id"`
	UserId  string `gorm:"column:user_id"`
	VideoId string `gorm:"column:video_id"`
	Content string `gorm:"column:content"`
}

func (c *CommentRaw) TableName() string {
	return "comment"
}

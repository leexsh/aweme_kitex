package controller

import (
	"aweme_kitex/utils"
	"errors"
	"sync"
	"time"

	"gorm.io/gorm"
)

// video
type VideoRawData struct {
	VideoId        string    `gorm:"column:video_id"`
	UserId         string    `gorm:"column:user_id"`
	Title          string    `gorm:"column:title"`
	PlayUrl        string    `gorm:"column:play_url"`
	CoverUrl       string    `gorm:"column:cover_url"`
	FavouriteCount int64     `gorm:"column:favourite_count"`
	CommentCount   int64     `gorm:"column:comment_count"`
	CreatedTime    time.Time `gorm:"column:created_at"`
	UpdatedTime    time.Time `gorm:"column:updated_at"`
	DeletedTime    time.Time `gorm:"column:deleted_at"`
}

func (vr *VideoRawData) TableName() string {
	return "video"
}

type VideoDao struct {
}

var (
	videoDao  *VideoDao
	videoOnce sync.Once
)

func NewVideoDaoInstance() *VideoDao {
	videoOnce.Do(
		func() {
			videoDao = &VideoDao{}
		})
	return videoDao
}

// 根据最新的时间戳获取视频信息
func (v *VideoDao) QueryVideoByLatestTime(latestTime int64) ([]*VideoRawData, error) {
	var videos []*VideoRawData
	// err := db.Table("video").Debug().Limit(20).Order("created_at desc").Where("created_at<?", time.Unix(int64(latestTime), 0)).Find(&videos).Error
	err := db.Table("video").Debug().Limit(20).Order("created_at desc").Where("created_at<?", utils.UnixToTime(latestTime)).Find(&videos).Error
	if err == gorm.ErrRecordNotFound {
		return nil, errors.New("not found videos")
	}
	if err != nil {
		return nil, errors.New("found videos error")
	}
	return videos, nil
}

// user
type UserRawData struct {
	UserId        string    `gorm:"column:user_id"`
	Name          string    `gorm:"column:name"`
	Password      string    `gorm:"column:password"`
	Token         string    `gorm:"column:token"`
	FollowCount   int64     `gorm:"column:follow_count"`
	FollowerCount int64     `gorm:"column:follower_count"`
	CreatedTime   time.Time `gorm:"column:created_at"`
	UpdatedTime   time.Time `gorm:"column:updated_at"`
	DeletedTime   time.Time `gorm:"column:deleted_at"`
}

func (u2 *UserRawData) TableName() string {
	return "user"
}

type UserDao struct {
}

// type TokenDao struct {
// }

var (
	userDao *UserDao
	// tokenDao  *TokenDao
	userOnce sync.Once
	// tokenOnce sync.Once
)

// func NewTokenDaoInstance() *TokenDao {
// 	tokenOnce.Do(
// 		func() {
// 			tokenDao = &TokenDao{}
// 		})
// 	return tokenDao
// }

func NewUserDaoInstance() *UserDao {
	userOnce.Do(
		func() {
			userDao = &UserDao{}
		})
	return userDao
}

// 根据token获取用户
// func (t *UserDao) QueryUserIdByToken(identity string) (string, error) {
// 	var user UserRawData
// 	err := db.Where("token=?", identity).Find(&user).Error
// 	if err != nil {
// 		return "", err
// 	}
// 	return user.UserId, nil
// }

func (u2 *UserDao) QueryUserByIds(uIds []string) (map[string]*UserRawData, error) {
	var users []*UserRawData
	err := db.Debug().Where("user_id in (?)", uIds).Find(&users).Error
	if err != nil {
		return nil, errors.New("query users fail")
	}
	userMap := make(map[string]*UserRawData)
	for _, user := range users {
		userMap[user.UserId] = user
	}
	return userMap, nil
}

// 喜欢
type FavouriteRaw struct {
	Id      string `gorm:"column:identity"`
	UserId  string `gorm:"column:user_id"`
	VideoId string `gorm:"column:video_id"`
}

func (f *FavouriteRaw) TableName() string {
	return "favourite"
}

type FavouriteDao struct {
}

var (
	favouriteDao  *FavouriteDao
	favouriteOnce sync.Once
)

func NewFavouriteDaoInstance() *FavouriteDao {
	favouriteOnce.Do(
		func() {
			favouriteDao = &FavouriteDao{}
		})
	return favouriteDao
}

// 根据uid 和 videoId 获取喜欢列表
func (f *FavouriteDao) QueryFavoursByIds(currentUId string, videoIds []string) (map[string]*FavouriteRaw, error) {
	var favours []*FavouriteRaw
	err := db.Table("favourite").Where("user_id=? AND video_id IN ?", currentUId, videoIds).Find(&favours).Error
	if err == gorm.ErrRecordNotFound {
		return nil, errors.New("favourite not found")
	}
	if err != nil {
		return nil, err
	}
	favoursMap := make(map[string]*FavouriteRaw)
	for _, favour := range favours {
		favoursMap[favour.VideoId] = favour
	}
	return favoursMap, nil
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

type RelationDao struct {
}

var (
	relationDao  *RelationDao
	relationOnce sync.Once
)

func NewRelationDaoInstance() *RelationDao {
	relationOnce.Do(
		func() {
			relationDao = &RelationDao{}
		})
	return relationDao
}

// 根据当前用户Id和视频作者的id获取关注信息
func (r *RelationDao) QueryRelationByIds(currentUid string, userIds []string) (map[string]*RelationRaw, error) {
	var relations []*RelationRaw
	err := db.Where("user_id=? AND to_user_id IN ? AND status IN ?", currentUid, userIds, []int64{0, -1}).
		Or("user_id IN ? AND to_user_id = ? AND status = ?", userIds, currentUid, 1).Find(&relations).Error

	if err == gorm.ErrRecordNotFound {
		return nil, errors.New("relation record not found")
	}
	if err != nil {
		return nil, errors.New("query relation record fail")
	}
	relationMap := make(map[string]*RelationRaw)
	for _, relation := range relations {
		if relation.Status == 1 {
			relationMap[relation.UserId] = relation
		} else {
			relationMap[relation.ToUserId] = relation
		}
	}
	return relationMap, nil
}

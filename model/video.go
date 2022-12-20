package model

import (
	"aweme_kitex/cfg"
	"aweme_kitex/utils"
	"errors"
	"sync"

	"gorm.io/gorm"
)

// ------------------------video---------------------

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
	err := cfg.DB.Table("video").Debug().Limit(20).Order("created_at desc").Where("created_at<?", utils.UnixToTime(latestTime)).Find(&videos).Error
	if err == gorm.ErrRecordNotFound {
		return nil, errors.New("not found videos")
	}
	if err != nil {
		return nil, errors.New("found videos error")
	}
	return videos, nil
}

func (v *VideoDao) QueryVideosByUserId(userId string) ([]*VideoRawData, error) {
	var videos []*VideoRawData
	err := cfg.DB.Table("video").Limit(20).Debug().Order("created_at desc").Where("user_id=?", userId).Find(&videos).Error
	if err == gorm.ErrRecordNotFound {
		return nil, errors.New("noe found videos")
	}
	if err != nil {
		return nil, err
	}
	return videos, nil
}

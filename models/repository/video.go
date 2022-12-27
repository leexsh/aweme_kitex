package repository

import (
	"aweme_kitex/cfg"
	"aweme_kitex/models"
	"aweme_kitex/utils"
	"errors"
	"mime/multipart"
	"sync"

	"github.com/gin-gonic/gin"
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

func (v *VideoDao) QueryVideoByLatestTime(latestTime int64) ([]*models.VideoRawData, error) {
	var videos []*models.VideoRawData
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

func (v *VideoDao) QueryVideosByUserId(userId string) ([]*models.VideoRawData, error) {
	var videos []*models.VideoRawData
	err := cfg.DB.Table("video").Limit(20).Debug().Order("created_at desc").Where("user_id=?", userId).Find(&videos).Error
	if err == gorm.ErrRecordNotFound {
		return nil, errors.New("noe found videos")
	}
	if err != nil {
		return nil, err
	}
	return videos, nil
}

// 将视频保存到本地文件夹中
func (*VideoDao) PublishVideoToPublic(video *multipart.FileHeader, path string, c *gin.Context) error {
	if err := c.SaveUploadedFile(video, path); err != nil {
		return err
	}
	return nil
}

func (*VideoDao) SaveVideoData(videoData *models.VideoRawData) error {
	err := DB.Table("video").Debug().Create(videoData).Error
	if err != nil {
		return err
	}
	return nil
}

// update favourite
func (*VideoDao) UpdateFavouriteCount(videoId, action string) error {
	var err error
	if action == "1" {
		err = DB.Table("video").Where("video_id=?", videoId).Update("favourite_count", gorm.Expr("favourite_count + ?", 1)).Error
	} else if action == "2" {
		err = DB.Table("video").Where("video_id=?", videoId).Update("favourite_count", gorm.Expr("favourite_count - ?", 1)).Error
	}
	return err
}

func (*VideoDao) QueryVideosByIs(videoId []string) ([]*models.VideoRawData, error) {
	var videos []*models.VideoRawData
	err := DB.Table("video").Where("video_id in (?)", videoId).Find(&videos).Error
	if err != nil {
		return nil, err
	}
	return videos, nil
}

// 通过视频id增加视频的评论数
func (*VideoDao) UpdateCommentCount(videoId string, action string) error {
	var err error
	if action == "1" {
		err = DB.Table("video").Where("video_id = ?", videoId).Update("comment_count", gorm.Expr("comment_count + ?", 1)).Error
	} else if action == "2" {
		err = DB.Table("video").Where("video_id = ?", videoId).Update("comment_count", gorm.Expr("comment_count - ?", 1)).Error
	}
	return err
}

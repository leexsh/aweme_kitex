package dal

import (
	"aweme_kitex/cfg"
	"aweme_kitex/models"
	"aweme_kitex/pkg/logger"
	"aweme_kitex/utils"
	"context"
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

func (v *VideoDao) QueryVideoByLatestTime(ctx context.Context, latestTime int64) ([]*models.VideoRawData, error) {
	var videos []*models.VideoRawData
	err := cfg.DB.WithContext(ctx).Table("video").Debug().Limit(20).Order("created_at desc").Where("created_at<?", utils.UnixToTime(latestTime)).Find(&videos).Error
	if err == gorm.ErrRecordNotFound {
		return nil, errors.New("not found videos")
	}
	if err != nil {
		logger.Error("query videos by latest time error : " + err.Error())
		return nil, errors.New("found videos error")
	}
	return videos, nil
}

func (v *VideoDao) QueryVideosByUserId(ctx context.Context, userId string) ([]*models.VideoRawData, error) {
	var videos []*models.VideoRawData
	err := cfg.DB.WithContext(ctx).Table("video").Limit(20).Debug().Order("created_at desc").Where("user_id=?", userId).Find(&videos).Error
	if err == gorm.ErrRecordNotFound {
		return nil, errors.New("noe found videos")
	}
	if err != nil {
		logger.Error("query videos by userId error : " + err.Error())
		return nil, err
	}
	return videos, nil
}

// 将视频保存到本地文件夹中
func (*VideoDao) PublishVideoToPublic(ctx context.Context, video *multipart.FileHeader, path string, c *gin.Context) error {
	if err := c.SaveUploadedFile(video, path); err != nil {
		logger.Error("save videos to local error : " + err.Error())
		return err
	}
	return nil
}

func (*VideoDao) SaveVideoData(ctx context.Context, videoData *models.VideoRawData) error {
	err := cfg.DB.WithContext(ctx).Table("video").Debug().Create(videoData).Error
	if err != nil {
		logger.Error("create video error : " + err.Error())
		return err
	}
	return nil
}

func (*VideoDao) QueryVideosByIs(ctx context.Context, videoId []string) ([]*models.VideoRawData, error) {
	var videos []*models.VideoRawData
	err := cfg.DB.WithContext(ctx).Table("video").Where("video_id in (?)", videoId).Find(&videos).Error
	if err != nil {
		logger.Error("query video by id error : " + err.Error())
		return nil, err
	}
	return videos, nil
}

package videoDB

import (
	"aweme_kitex/cfg"
	"aweme_kitex/pkg/logger"
	"aweme_kitex/pkg/types"
	"aweme_kitex/pkg/utils"
	"context"
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

func (v *VideoDao) QueryVideoByLatestTime(ctx context.Context, latestTime int64) ([]*types.VideoRawData, error) {
	var videos []*types.VideoRawData
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

func (v *VideoDao) QueryVideosByUserId(ctx context.Context, userId string) ([]*types.VideoRawData, error) {
	var videos []*types.VideoRawData
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

func (*VideoDao) QueryVideosByIs(ctx context.Context, videoId []string) ([]*types.VideoRawData, error) {
	var videos []*types.VideoRawData
	err := cfg.DB.WithContext(ctx).Table("video").Where("video_id in (?)", videoId).Find(&videos).Error
	if err != nil {
		logger.Error("query video by id error : " + err.Error())
		return nil, err
	}
	return videos, nil
}

func (v *VideoDao) CheckVideoId(ctx context.Context, videoId []string) ([]*types.VideoRawData, error) {
	videos, err := v.QueryVideosByIs(ctx, videoId)
	if err != nil {
		return nil, err
	}
	if len(videos) == 0 {
		return nil, errors.New("video not exist")
	}
	return videos, nil
}

// 增加评论后修改评论数目
func (v *VideoDao) IncreaseCommentCount(ctx context.Context, videoId string) error {
	err := cfg.DB.WithContext(ctx).Table("video").Where("video_id = ?", videoId).Update("comment_count", gorm.Expr("comment_count + ?", 1)).Error
	return err
}

// 删除评论后修改评论数目
func (v *VideoDao) DecreaseCommentCount(ctx context.Context, videoId string) error {
	err := cfg.DB.WithContext(ctx).Table("video").Where("video_id = ?", videoId).Update("comment_count", gorm.Expr("comment_count - ?", 1)).Error
	return err
}

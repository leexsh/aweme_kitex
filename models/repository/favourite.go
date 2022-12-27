package repository

import (
	"aweme_kitex/models"
	"errors"
	"sync"

	"gorm.io/gorm"
)

// ---------------收藏---------------
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
func (f *FavouriteDao) QueryFavoursByIds(currentUId string, videoIds []string) (map[string]*models.FavouriteRaw, error) {
	var favours []*models.FavouriteRaw
	err := DB.Table("favourite").Where("user_id=? AND video_id IN ?", currentUId, videoIds).Find(&favours).Error
	if err == gorm.ErrRecordNotFound {
		return nil, errors.New("favourite not found")
	}
	if err != nil {
		return nil, err
	}
	favoursMap := make(map[string]*models.FavouriteRaw)
	for _, favour := range favours {
		favoursMap[favour.VideoId] = favour
	}
	return favoursMap, nil
}

// 创建一条点赞
func (f *FavouriteDao) CreateFavour(favour *models.FavouriteRaw) error {
	err := DB.Debug().Table("favourite").Create(favour).Error
	if err != nil {
		return err
	}
	return nil
}

// del
func (f *FavouriteDao) DelFavour(userId, videoId string) error {
	var favour *models.FavouriteRaw
	err := DB.Table("favourite").Where("user_id = ? AND video_id = ?", userId, videoId).Delete(favour).Error
	if err != nil {
		return err
	}
	return nil
}

// quiery videos by uid
func (f *FavouriteDao) QueryFavoursVideoIdByUid(uid string) ([]string, error) {
	var favours []*models.FavouriteRaw
	err := DB.Debug().Table("favourite").Where("user_id=?", uid).Find(&favours).Error
	if err != nil {
		return nil, err
	}
	var videos []string
	for _, favour := range favours {
		videos = append(videos, favour.VideoId)
	}
	return videos, nil
}

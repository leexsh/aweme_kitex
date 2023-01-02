package repository

import (
	"aweme_kitex/cfg"
	"aweme_kitex/models"
	"aweme_kitex/utils"
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
	err := cfg.DB.Table("favourite").Where("user_id=? AND video_id IN ?", currentUId, videoIds).Find(&favours).Error
	if err == gorm.ErrRecordNotFound {
		return nil, errors.New("favourite not found")
	}
	if err != nil {
		utils.Error("query favours by id err: " + err.Error())
		return nil, err
	}
	favoursMap := make(map[string]*models.FavouriteRaw)
	for _, favour := range favours {
		favoursMap[favour.VideoId] = favour
	}
	return favoursMap, nil
}

// 创建一条点赞
func (f *FavouriteDao) CreateFavour(favour *models.FavouriteRaw, videoId string) error {
	cfg.DB.Transaction(func(tx *gorm.DB) error {
		err := tx.Table("video").Where("video_id = ?", videoId).Update("favorite_count", gorm.Expr("favorite_count + ?", 1)).Error
		if err != nil {
			utils.Error("AddFavoriteCount error " + err.Error())
			return err
		}

		err = tx.Table("favorite").Create(favour).Error
		if err != nil {
			utils.Error("create favorite record fail " + err.Error())
			return err
		}

		return nil
	})
	return nil
}

// del
func (f *FavouriteDao) DelFavour(userId, videoId string) error {
	cfg.DB.Transaction(func(tx *gorm.DB) error {
		var favorite *models.FavouriteRaw
		err := tx.Table("favorite").Where("user_id = ? AND video_id = ?", userId, videoId).Delete(&favorite).Error
		if err != nil {
			utils.Error("delete favorite record fail " + err.Error())
			return err
		}

		err = tx.Table("video").Where("video_id = ?", videoId).Update("favorite_count", gorm.Expr("favorite_count - ?", 1)).Error
		if err != nil {
			utils.Error("SubFavoriteCount error " + err.Error())
			return err
		}
		return nil
	})
	return nil
}

// quiery videos by uid
func (f *FavouriteDao) QueryFavoursVideoIdByUid(uid string) ([]string, error) {
	var favours []*models.FavouriteRaw
	err := cfg.DB.Debug().Table("favourite").Where("user_id=?", uid).Find(&favours).Error
	if err != nil {
		utils.Error("query favourite video err: " + err.Error())
		return nil, err
	}
	var videos []string
	for _, favour := range favours {
		videos = append(videos, favour.VideoId)
	}
	return videos, nil
}

package service

import (
	"aweme_kitex/model"
	"errors"
	"fmt"
	"sync"
)

// --------------------
type QueryUserVideoList struct {
	UserName string
	UserId   string

	VideoList    []model.Video
	VideoData    []*model.VideoRawData
	UserMap      map[string]*model.UserRawData
	FavouriteMap map[string]*model.FavouriteRaw
	RelationMap  map[string]*model.RelationRaw
}

func NewQueryUserVideoList(userId string) *QueryUserVideoList {
	return &QueryUserVideoList{
		UserId: userId,
	}
}

func (f *QueryUserVideoList) prepareVideoInfo() error {
	videoData, err := model.NewVideoDaoInstance().QueryVideosByUserId(f.UserId)
	if err != nil {
		return err
	}
	f.VideoData = videoData

	videoIds := make([]string, 0)
	userIds := []string{f.UserId}
	for _, video := range f.VideoData {
		videoIds = append(videoIds, video.VideoId)
	}

	userMap, err := model.NewUserDaoInstance().QueryUserByIds(userIds)
	if err != nil {
		return err
	}
	f.UserMap = userMap

	var wg sync.WaitGroup
	wg.Add(2)
	var favoriteErr, relationErr error
	// 获取点赞信息
	go func() {
		defer wg.Done()
		favoriteMap, err := model.NewFavouriteDaoInstance().QueryFavoursByIds(f.UserId, videoIds)
		if err != nil {
			favoriteErr = err
			return
		}
		f.FavouriteMap = favoriteMap
	}()
	// 获取关注信息
	go func() {
		defer wg.Done()
		relationMap, err := model.NewRelationDaoInstance().QueryRelationByIds(f.UserId, userIds)
		if err != nil {
			relationErr = err
			return
		}
		f.RelationMap = relationMap
	}()
	wg.Wait()
	if favoriteErr != nil {
		return favoriteErr
	}
	if relationErr != nil {
		return relationErr
	}
	return nil
}

func (f *QueryUserVideoList) packVideoInfo() error {
	videoList := make([]model.Video, 0)
	for _, video := range f.VideoData {
		videoUser, ok := f.UserMap[video.UserId]
		if !ok {
			return errors.New("has no video user info for " + fmt.Sprint(video.UserId))
		}

		var isFavorite bool = false
		var isFollow bool = false

		_, ok = f.FavouriteMap[video.VideoId]
		if ok {
			isFavorite = true
		}
		_, ok = f.RelationMap[video.UserId]
		if ok {
			isFollow = true
		}
		videoList = append(videoList, model.Video{
			Id: video.VideoId,
			Author: model.User{
				UserId:        videoUser.UserId,
				Name:          videoUser.Name,
				FollowCount:   videoUser.FollowCount,
				FollowerCount: videoUser.FollowerCount,
				IsFollow:      isFollow,
			},
			PlayUrl:        video.PlayUrl,
			CoverUrl:       video.CoverUrl,
			FavouriteCount: video.FavouriteCount,
			CommentCount:   video.CommentCount,
			IsFavourite:    isFavorite,
			Title:          video.Title,
		})
	}

	f.VideoList = videoList
	return nil
}

func (f *QueryUserVideoList) do() ([]model.Video, error) {
	if err := f.prepareVideoInfo(); err != nil {
		return nil, err
	}
	if err := f.packVideoInfo(); err != nil {
		return nil, err
	}
	return f.VideoList, nil
}

func QueryUserVideos(token string) ([]model.Video, error) {
	return NewQueryUserVideoList(token).do()
}

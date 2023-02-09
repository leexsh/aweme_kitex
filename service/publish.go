package service

import (
	"aweme_kitex/models"
	"aweme_kitex/models/dal"
	"context"
	"errors"
	"fmt"
	"sync"
)

// -------------------- published list
type userVideoList struct {
	UserName string
	UserId   string

	VideoList    []*models.Video
	VideoData    []*models.VideoRawData
	UserMap      map[string]*models.UserRawData
	FavouriteMap map[string]*models.FavouriteRaw
	RelationMap  map[string]*models.RelationRaw
}

func newQueryUserVideoList(userId string) *userVideoList {
	return &userVideoList{
		UserId: userId,
	}
}

func (f *userVideoList) prepareVideoInfo() error {
	videoData, err := dal.NewVideoDaoInstance().QueryVideosByUserId(context.Background(), f.UserId)
	if err != nil {
		return err
	}
	f.VideoData = videoData

	videoIds := make([]string, 0)
	userIds := []string{f.UserId}
	for _, video := range f.VideoData {
		videoIds = append(videoIds, video.VideoId)
	}

	users, err := dal.NewUserDaoInstance().QueryUserByIds(context.Background(), userIds)
	if err != nil {
		return err
	}
	userMap := make(map[string]*models.UserRawData)
	for _, user := range users {
		userMap[user.UserId] = user
	}
	f.UserMap = userMap

	var wg sync.WaitGroup
	wg.Add(2)
	var favoriteErr, relationErr error
	// 获取点赞信息
	go func() {
		defer wg.Done()
		favoriteMap, err := dal.NewFavouriteDaoInstance().QueryFavoursByIds(context.Background(), f.UserId, videoIds)
		if err != nil {
			favoriteErr = err
			return
		}
		f.FavouriteMap = favoriteMap
	}()
	// 获取关注信息
	go func() {
		defer wg.Done()
		relationMap, err := dal.NewRelationDaoInstance().QueryRelationByIds(context.Background(), f.UserId, userIds)
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

func (f *userVideoList) packVideoInfo() error {
	videoList := make([]*models.Video, 0)
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
		videoList = append(videoList, &models.Video{
			Id: video.VideoId,
			Author: &models.User{
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

func (f *userVideoList) do() ([]*models.Video, error) {
	if err := f.prepareVideoInfo(); err != nil {
		return nil, err
	}
	if err := f.packVideoInfo(); err != nil {
		return nil, err
	}
	return f.VideoList, nil
}

func QueryUserVideos(userId string) ([]*models.Video, error) {
	return newQueryUserVideoList(userId).do()
}

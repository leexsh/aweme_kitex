package service

import (
	"aweme_kitex/models"
	"aweme_kitex/models/dal"
	"context"
	"errors"
	"fmt"
	"sync"
)

// ------------service_user--------------------
// 该层负责鉴权  向repository获取视频数据和封装数据

func QueryVideoData(latestTime int64, userId string) ([]*models.Video, int64, error) {
	return newQueryVideoDataFlow(latestTime, userId).Do()
}

func newQueryVideoDataFlow(latestTime int64, userId string) *queryVideoDataFlow {
	return &queryVideoDataFlow{
		LatestTime:    latestTime,
		CurrentUserId: userId,
	}
}

// video data
type queryVideoDataFlow struct {
	CurrentUserId   string
	CurrentUserName string

	LatestTime int64
	VideoList  []*models.Video
	NextTime   int64

	VideoData   []*models.VideoRawData
	UserMap     map[string]*models.UserRawData
	FavoursMap  map[string]*models.FavouriteRaw
	RelationMap map[string]*models.RelationRaw
}

func (f *queryVideoDataFlow) Do() ([]*models.Video, int64, error) {
	if err := f.prepareVideoInfo(); err != nil {
		return nil, 0, err
	}
	if err := f.packVideoInfo(); err != nil {
		return nil, 0, err
	}
	return f.VideoList, f.NextTime, nil
}

// prepare video
func (f *queryVideoDataFlow) prepareVideoInfo() error {
	// 1.get video
	videoData, err := dal.NewVideoDaoInstance().QueryVideoByLatestTime(context.Background(), f.LatestTime)
	if err != nil {
		return err
	}
	f.VideoData = videoData

	// 2. get video_id and user_id
	videoIds := make([]string, 0)
	authorIds := make([]string, 0)
	for _, video := range f.VideoData {
		videoIds = append(videoIds, video.VideoId)
		authorIds = append(authorIds, video.UserId)
	}

	// 3. get user info
	users, err := dal.NewUserDaoInstance().QueryUserByIds(context.Background(), authorIds)
	if err != nil {
		return err
	}
	userMap := make(map[string]*models.UserRawData)
	for _, user := range users {
		userMap[user.UserId] = user
	}
	f.UserMap = userMap

	// 4. should login
	if f.CurrentUserId == "" {
		return nil
	}

	var wg sync.WaitGroup
	wg.Add(2)
	var favourErr, relationErr error
	// 5.获取点赞信息
	go func() {
		defer wg.Done()
		favoursMap, err := dal.NewFavouriteDaoInstance().QueryFavoursByIds(context.Background(), f.CurrentUserId, videoIds)
		if err != nil {
			favourErr = err
			return
		}
		f.FavoursMap = favoursMap
	}()
	// 6.获取关注信息
	go func() {
		defer wg.Done()
		relationMap, err := dal.NewRelationDaoInstance().QueryRelationByIds(context.Background(), f.CurrentUserId, authorIds)
		if err != nil {
			relationErr = err
			return
		}
		f.RelationMap = relationMap

	}()
	wg.Wait()

	if favourErr != nil {
		return favourErr
	}
	if relationErr != nil {
		return relationErr
	}
	return nil
}

func (f *queryVideoDataFlow) packVideoInfo() error {
	videoList := make([]*models.Video, 0)
	for _, video := range f.VideoData {
		videoAuthor, ok := f.UserMap[video.UserId]
		if !ok {
			return errors.New("has no video user info for " + fmt.Sprint(video.UserId))
		}
		var isFavourite bool = false
		var isFollow bool = false
		if f.CurrentUserId != "" {
			if _, ok := f.FavoursMap[video.VideoId]; ok {
				isFavourite = true
			}
			if _, ok := f.RelationMap[video.UserId]; ok {
				isFollow = true
			}
		}
		videoList = append(videoList, &models.Video{
			Id: video.VideoId,
			Author: &models.User{
				UserId:        videoAuthor.UserId,
				Name:          videoAuthor.Name,
				FollowCount:   videoAuthor.FollowCount,
				FollowerCount: videoAuthor.FollowerCount,
				IsFollow:      isFollow,
			},
			PlayUrl:        video.PlayUrl,
			CoverUrl:       video.CoverUrl,
			Title:          video.Title,
			FavouriteCount: video.FavouriteCount,
			CommentCount:   video.CommentCount,
			IsFavourite:    isFavourite,
		})
	}
	f.VideoList = videoList
	if len(f.VideoList) == 0 {
		f.NextTime = 0
	} else {
		f.NextTime = f.VideoData[len(f.VideoData)-1].CreatedAt.Unix()
	}
	return nil
}

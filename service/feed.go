package service

import (
	"aweme_kitex/model"
	"errors"
	"fmt"
	"sync"
)

// ------------service--------------------
// 该层负责鉴权  向repository获取视频数据和封装数据

func QueryVideoData(latestTime int64, userId string) ([]model.Video, int64, error) {
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
	VideoList  []model.Video
	NextTime   int64

	VideoData   []*model.VideoRawData
	UserMap     map[string]*model.UserRawData
	FavoursMap  map[string]*model.FavouriteRaw
	RelationMap map[string]*model.RelationRaw
}

func (f *queryVideoDataFlow) Do() ([]model.Video, int64, error) {
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
	videoData, err := model.NewVideoDaoInstance().QueryVideoByLatestTime(f.LatestTime)
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
	usersMap, err := model.NewUserDaoInstance().QueryUserByIds(authorIds)
	if err != nil {
		return err
	}
	f.UserMap = usersMap

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
		favoursMap, err := model.NewFavouriteDaoInstance().QueryFavoursByIds(f.CurrentUserId, videoIds)
		if err != nil {
			favourErr = err
			return
		}
		f.FavoursMap = favoursMap
	}()
	// 6.获取关注信息
	go func() {
		defer wg.Done()
		relationMap, err := model.NewRelationDaoInstance().QueryRelationByIds(f.CurrentUserId, authorIds)
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
	videoList := make([]model.Video, 0)
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
		videoList = append(videoList, model.Video{
			Id: video.VideoId,
			Author: model.User{
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

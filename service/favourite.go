package service

import (
	"aweme_kitex/model"
	"aweme_kitex/utils"
	"errors"
	"fmt"
	"sync"
)

func FavouriteActionService(user *model.UserClaim, videoId, action string) error {
	return newFavouriteActionData(user, videoId, action).do()
}

type favouriteActionDataFlow struct {
	CurrentUid  string
	CurrentName string
	videoId     string
	action      string

	FavouriteData *model.FavouriteRaw
}

func newFavouriteActionData(user *model.UserClaim, videoId, action string) *favouriteActionDataFlow {
	return &favouriteActionDataFlow{
		CurrentUid:  user.Id,
		CurrentName: user.Name,
		videoId:     videoId,
		action:      action,
	}
}

func (f *favouriteActionDataFlow) do() error {
	if f.action != "1" && f.action != "2" {
		return errors.New("invalid action type")
	}
	if f.action == "1" {
		favour := &model.FavouriteRaw{
			Id:      utils.GenerateUUID(),
			UserId:  f.CurrentUid,
			VideoId: f.videoId,
		}
		f.FavouriteData = favour
		var wg sync.WaitGroup
		wg.Add(2)
		var favourErr, videoErr error
		// 事务
		go func() {
			defer wg.Done()
			err := model.NewFavouriteDaoInstance().CreateFavour(favour)
			if err != nil {
				favourErr = err
			}
		}()
		go func() {
			wg.Done()
			err := model.NewVideoDaoInstance().UpdateFavouriteCount(f.videoId, f.action)
			if err != nil {
				videoErr = err
			}
		}()
		wg.Wait()
		if favourErr != nil {
			return favourErr
		}
		if videoErr != nil {
			return videoErr
		}

	} else if f.action == "2" {
		var wg sync.WaitGroup
		wg.Add(2)
		var favourErr, videoErr error
		// 事务
		go func() {
			defer wg.Done()
			err := model.NewFavouriteDaoInstance().DelFavour(f.CurrentUid, f.videoId)
			if err != nil {
				favourErr = err
			}
		}()
		go func() {
			wg.Done()
			err := model.NewVideoDaoInstance().UpdateFavouriteCount(f.videoId, f.action)
			if err != nil {
				videoErr = err
			}
		}()
		wg.Wait()
		if favourErr != nil {
			return favourErr
		}
		if videoErr != nil {
			return videoErr
		}

	}
	return nil
}

func FavouriteListService(userId, userName string) ([]model.Video, error) {
	return newFavouriteListDataFlow(userId, userName).do()
}

type favouriteListDataFlow struct {
	currentUId   string
	currentUName string

	favours []model.Video

	videoRawData []*model.VideoRawData
	users        map[string]*model.UserRawData
	favoursMap   map[string]*model.FavouriteRaw
	RelationMap  map[string]*model.RelationRaw
}

func newFavouriteListDataFlow(id, name string) *favouriteListDataFlow {
	return &favouriteListDataFlow{
		currentUName: name,
		currentUId:   id,
	}
}

func (f *favouriteListDataFlow) do() ([]model.Video, error) {
	if err := f.prepareVideoInfo(); err != nil {
		return nil, err
	}
	if err := f.packageVideo(); err != nil {
		return nil, err
	}
	return f.favours, nil
}

func (f *favouriteListDataFlow) prepareVideoInfo() error {
	videosIds, err := model.NewFavouriteDaoInstance().QueryFavoursVideoIdByUid(f.currentUId)
	if err != nil {
		return err
	}
	// get videos
	videoData, err := model.NewVideoDaoInstance().QueryVideosByIs(videosIds)
	if err != nil {
		return err
	}
	f.videoRawData = videoData

	uids := []string{}
	for _, video := range f.videoRawData {
		uids = append(uids, video.UserId)
	}

	// get video authors
	userMap, err := model.NewUserDaoInstance().QueryUserByIds(uids)
	if err != nil {
		return err
	}
	f.users = userMap

	var wg sync.WaitGroup
	wg.Add(2)
	var favErr, relationErr error
	go func() {
		defer wg.Done()
		favoursMap, err := model.NewFavouriteDaoInstance().QueryFavoursByIds(f.currentUId, videosIds)
		if err != nil {
			favErr = err
		}
		f.favoursMap = favoursMap
	}()
	go func() {
		defer wg.Done()
		relationMap, err := model.NewRelationDaoInstance().QueryRelationByIds(f.currentUId, videosIds)
		if err != nil {
			relationErr = err
		}
		f.RelationMap = relationMap
	}()
	wg.Wait()
	if favErr != nil {
		return favErr
	}
	if relationErr != nil {
		return relationErr
	}
	return nil

}

func (f *favouriteListDataFlow) packageVideo() error {
	videoList := make([]model.Video, 0)
	for _, video := range f.videoRawData {
		author, ok := f.users[video.UserId]
		if !ok {
			return errors.New("has no video user info for " + fmt.Sprint(video.UserId))
		}
		var isFavour bool = false
		var isFollow bool = false
		_, ok = f.favoursMap[video.VideoId]
		if ok {
			isFavour = true
		}
		_, ok = f.RelationMap[video.UserId]
		if ok {
			isFollow = true
		}
		videoList = append(videoList, model.Video{
			Id: video.VideoId,
			Author: model.User{
				UserId:        author.UserId,
				Name:          author.Name,
				FollowCount:   author.FollowCount,
				FollowerCount: author.FollowerCount,
				IsFollow:      isFollow,
			},
			PlayUrl:        video.PlayUrl,
			CoverUrl:       video.CoverUrl,
			FavouriteCount: video.FavouriteCount,
			CommentCount:   video.CommentCount,
			IsFavourite:    isFavour,
			Title:          video.Title,
		})
	}
	f.favours = videoList
	return nil
}

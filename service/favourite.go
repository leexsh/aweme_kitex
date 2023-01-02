package service

import (
	"aweme_kitex/models"
	"aweme_kitex/models/repository"
	"aweme_kitex/utils"
	"errors"
	"fmt"
	"sync"
)

func FavouriteActionService(user *models.UserClaim, videoId, action string) error {
	return newFavouriteActionData(user, videoId, action).do()
}

type favouriteActionDataFlow struct {
	CurrentUid  string
	CurrentName string
	videoId     string
	action      string

	FavouriteData *models.FavouriteRaw
}

func newFavouriteActionData(user *models.UserClaim, videoId, action string) *favouriteActionDataFlow {
	return &favouriteActionDataFlow{
		CurrentUid:  user.Id,
		CurrentName: user.Name,
		videoId:     videoId,
		action:      action,
	}
}

func (f *favouriteActionDataFlow) do() error {
	if _, err := checkVideoId([]string{f.videoId}); err != nil {
		return err
	}
	if f.action != "1" && f.action != "2" {
		return errors.New("invalid action type")
	}
	if f.action == "1" {
		favour := &models.FavouriteRaw{
			Id:      utils.GenerateUUID(),
			UserId:  f.CurrentUid,
			VideoId: f.videoId,
		}
		f.FavouriteData = favour
		err := repository.NewFavouriteDaoInstance().CreateFavour(favour, f.videoId)
		if err != nil {
			return err
		}
	} else if f.action == "2" {
		err := repository.NewFavouriteDaoInstance().DelFavour(f.CurrentUid, f.videoId)
		if err != nil {
			return err
		}
	}
	return nil
}

func FavouriteListService(userId, userName string) ([]*models.Video, error) {
	return newFavouriteListDataFlow(userId, userName).do()
}

type favouriteListDataFlow struct {
	currentUId   string
	currentUName string

	favours []*models.Video

	videoRawData []*models.VideoRawData
	users        map[string]*models.UserRawData
	favoursMap   map[string]*models.FavouriteRaw
	RelationMap  map[string]*models.RelationRaw
}

func newFavouriteListDataFlow(id, name string) *favouriteListDataFlow {
	return &favouriteListDataFlow{
		currentUName: name,
		currentUId:   id,
	}
}

func (f *favouriteListDataFlow) do() ([]*models.Video, error) {
	if err := f.prepareVideoInfo(); err != nil {
		return nil, err
	}
	if err := f.packageVideo(); err != nil {
		return nil, err
	}
	return f.favours, nil
}

func (f *favouriteListDataFlow) prepareVideoInfo() error {
	videosIds, err := repository.NewFavouriteDaoInstance().QueryFavoursVideoIdByUid(f.currentUId)
	if err != nil {
		return err
	}
	// get videos
	videoData, err := repository.NewVideoDaoInstance().QueryVideosByIs(videosIds)
	if err != nil {
		return err
	}
	f.videoRawData = videoData

	uids := []string{}
	for _, video := range f.videoRawData {
		uids = append(uids, video.UserId)
	}

	// get video authors
	users, err := repository.NewUserDaoInstance().QueryUserByIds(uids)
	if err != nil {
		return err
	}
	userMap := make(map[string]*models.UserRawData)
	for _, user := range users {
		userMap[user.UserId] = user
	}
	f.users = userMap

	var wg sync.WaitGroup
	wg.Add(2)
	var favErr, relationErr error
	go func() {
		defer wg.Done()
		favoursMap, err := repository.NewFavouriteDaoInstance().QueryFavoursByIds(f.currentUId, videosIds)
		if err != nil {
			favErr = err
		}
		f.favoursMap = favoursMap
	}()
	go func() {
		defer wg.Done()
		relationMap, err := repository.NewRelationDaoInstance().QueryRelationByIds(f.currentUId, videosIds)
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
	videoList := make([]*models.Video, 0)
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
		videoList = append(videoList, &models.Video{
			Id: video.VideoId,
			Author: &models.User{
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

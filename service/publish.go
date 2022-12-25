package service

import (
	"aweme_kitex/model"
	"aweme_kitex/utils"
	"errors"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"sync"

	"github.com/gin-gonic/gin"
)

// -------------publish
func PublishVideoService(userId, userName, title string, data *multipart.FileHeader) error {
	return newPublishVideoServiceData(userId, userName, title, data).Do()
}

func newPublishVideoServiceData(userId, userName, title string, data *multipart.FileHeader) *publishVideoServiceData {
	return &publishVideoServiceData{
		Data:            data,
		Title:           title,
		CurrentUserId:   userId,
		CurrentUserName: userName,
	}
}

type publishVideoServiceData struct {
	Data  *multipart.FileHeader
	Title string
	Gin   *gin.Context

	CurrentUserId   string
	CurrentUserName string
	Video           model.VideoRawData
}

func (f *publishVideoServiceData) Do() error {
	if err := f.publishVideo(); err != nil {
		return err
	}
	return nil
}
func (f *publishVideoServiceData) publishVideo() error {
	fileName := filepath.Base(f.Data.Filename)
	finalName := fmt.Sprintf("%s_%s", f.CurrentUserName, fileName)

	saveFile := filepath.Join("./public/", finalName)
	// 1.save public
	err := model.NewVideoDaoInstance().PublishVideoToPublic(f.Data, saveFile, f.Gin)
	if err != nil {
		return err
	}
	cosKey := fmt.Sprintf("%s/%s", f.CurrentUserName, finalName)
	var wg sync.WaitGroup
	wg.Add(2)
	var cosErr, sqlErr error
	go func() {
		// 2. save COS
		defer wg.Done()
		err := model.NewCOSDaoInstance().PublishVideoToCOS(cosKey, saveFile)
		if err != nil {
			cosErr = err
		}
	}()

	go func() {
		defer wg.Done()
		ourl := model.NewCOSDaoInstance().GetCOSVideoURL(cosKey)
		video := &model.VideoRawData{
			VideoId: utils.GenerateUUID(),
			UserId:  f.CurrentUserId,
			Title:   f.Title,
			PlayUrl: ourl.String(),
		}
		err := model.NewVideoDaoInstance().SaveVideoData(video)
		if err != nil {
			sqlErr = err
		}
	}()

	wg.Wait()
	if cosErr != nil {
		return cosErr
	}
	if sqlErr != nil {
		return sqlErr
	}
	return nil
}

// -------------------- published list
type userVideoList struct {
	UserName string
	UserId   string

	VideoList    []model.Video
	VideoData    []*model.VideoRawData
	UserMap      map[string]*model.UserRawData
	FavouriteMap map[string]*model.FavouriteRaw
	RelationMap  map[string]*model.RelationRaw
}

func newQueryUserVideoList(userId string) *userVideoList {
	return &userVideoList{
		UserId: userId,
	}
}

func (f *userVideoList) prepareVideoInfo() error {
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

func (f *userVideoList) packVideoInfo() error {
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

func (f *userVideoList) do() ([]model.Video, error) {
	if err := f.prepareVideoInfo(); err != nil {
		return nil, err
	}
	if err := f.packVideoInfo(); err != nil {
		return nil, err
	}
	return f.VideoList, nil
}

func QueryUserVideos(userId string) ([]model.Video, error) {
	return newQueryUserVideoList(userId).do()
}

package service_publish

import (
	"aweme_kitex/cmd/publish/kitex_gen/feed"
	"aweme_kitex/cmd/publish/kitex_gen/publish"
	"aweme_kitex/cmd/publish/kitex_gen/user"
	"aweme_kitex/cmd/relation/service_relation/db"
	db2 "aweme_kitex/cmd/user/service_user/db"
	"aweme_kitex/models"
	"aweme_kitex/models/dal"
	"aweme_kitex/pkg/jwt"
	"context"
	"errors"
	"fmt"
	"sync"
)

type PublishListService struct {
	ctx context.Context
}

// NewPublishService new PublishService
func NewPublishListService(ctx context.Context) *PublishListService {
	return &PublishListService{ctx: ctx}
}

func (s *PublishListService) PublishList(req *publish.PublishListRequest) ([]*feed.Video, error) {
	uc, _ := jwt.AnalyzeToken(req.Token)
	return queryUserVideos(s.ctx, uc.Id)
}

func queryUserVideos(ctx context.Context, userId string) ([]*feed.Video, error) {
	return newQueryUserVideoList(ctx, userId).do()
}

type userVideoList struct {
	ctx    context.Context
	UserId string

	VideoList    []*feed.Video
	VideoData    []*models.VideoRawData
	UserMap      map[string]*models.UserRawData
	FavouriteMap map[string]*models.FavouriteRaw
	RelationMap  map[string]*models.RelationRaw
}

func newQueryUserVideoList(ctx context.Context, userId string) *userVideoList {
	return &userVideoList{
		ctx:    ctx,
		UserId: userId,
	}
}

func (f *userVideoList) prepareVideoInfo() error {
	videoData, err := dal.NewVideoDaoInstance().QueryVideosByUserId(f.ctx, f.UserId)
	if err != nil {
		return err
	}
	f.VideoData = videoData

	videoIds := make([]string, 0)
	userIds := []string{f.UserId}
	for _, video := range f.VideoData {
		videoIds = append(videoIds, video.VideoId)
	}

	users, err := db2.NewUserDaoInstance().QueryUserByIds(f.ctx, userIds)
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
		favoriteMap, err := dal.NewFavouriteDaoInstance().QueryFavoursByIds(f.ctx, f.UserId, videoIds)
		if err != nil {
			favoriteErr = err
			return
		}
		f.FavouriteMap = favoriteMap
	}()
	// 获取关注信息
	go func() {
		defer wg.Done()
		relationMap, err := db.NewRelationDaoInstance().QueryRelationByIds(f.ctx, f.UserId, userIds)
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
	videoList := make([]*feed.Video, 0)
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
		curVideo := &feed.Video{
			VideoId: video.VideoId,
			Author: &user.User{
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
		}
		videoList = append(videoList, curVideo)
	}

	f.VideoList = videoList
	return nil
}

func (f *userVideoList) do() ([]*feed.Video, error) {
	if err := f.prepareVideoInfo(); err != nil {
		return nil, err
	}
	if err := f.packVideoInfo(); err != nil {
		return nil, err
	}
	return f.VideoList, nil
}

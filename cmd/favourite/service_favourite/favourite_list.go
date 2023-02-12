package service_favourite

import (
	"aweme_kitex/cmd/favourite/kitex_gen/favourite"
	"aweme_kitex/cmd/favourite/kitex_gen/feed"
	"aweme_kitex/cmd/favourite/kitex_gen/user"
	"aweme_kitex/models"
	"aweme_kitex/models/dal"
	"aweme_kitex/pkg/jwt"
	"context"
	"errors"
	"fmt"
	"sync"
)

type FavoriteListService struct {
	ctx context.Context
}

// NewFavoriteListService new FavoriteListService
func NewFavoriteListService(ctx context.Context) *FavoriteListService {
	return &FavoriteListService{ctx: ctx}
}

func (s *FavoriteListService) FavoriteList(req *favourite.FavouriteListRequest) ([]*feed.Video, error) {
	uc, err := jwt.AnalyzeToken(req.Token)
	if err != nil {
		return nil, err
	}
	videos, err := FavouriteListService(s.ctx, uc.Id, uc.Name)
	if err != nil {
		return nil, err
	}
	return videos, nil
}

// ---------FavouriteList Service---------------
func FavouriteListService(ctx context.Context, userId, userName string) ([]*feed.Video, error) {
	return newFavouriteListDataFlow(ctx, userId, userName).do()
}

type favouriteListDataFlow struct {
	currentUId string

	favours []*feed.Video

	videoRawData []*models.VideoRawData
	users        map[string]*models.UserRawData
	favoursMap   map[string]*models.FavouriteRaw
	RelationMap  map[string]*models.RelationRaw
	ctx          context.Context
}

func newFavouriteListDataFlow(ctx context.Context, id, name string) *favouriteListDataFlow {
	return &favouriteListDataFlow{
		currentUId: id,
		ctx:        ctx,
	}
}

func (f *favouriteListDataFlow) do() ([]*feed.Video, error) {
	if err := f.prepareVideoInfo(); err != nil {
		return nil, err
	}
	if err := f.packageVideo(); err != nil {
		return nil, err
	}
	return f.favours, nil
}

func (f *favouriteListDataFlow) prepareVideoInfo() error {
	videosIds, err := dal.NewFavouriteDaoInstance().QueryFavoursVideoIdByUid(f.ctx, f.currentUId)
	if err != nil {
		return err
	}
	// get videos
	videoData, err := dal.NewVideoDaoInstance().QueryVideosByIs(f.ctx, videosIds)
	if err != nil {
		return err
	}
	f.videoRawData = videoData

	uids := []string{}
	for _, video := range f.videoRawData {
		uids = append(uids, video.UserId)
	}

	// get video authors
	users, err := dal.NewUserDaoInstance().QueryUserByIds(f.ctx, uids)
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
		favoursMap, err := dal.NewFavouriteDaoInstance().QueryFavoursByIds(f.ctx, f.currentUId, videosIds)
		if err != nil {
			favErr = err
		}
		f.favoursMap = favoursMap
	}()
	go func() {
		defer wg.Done()
		relationMap, err := dal.NewRelationDaoInstance().QueryRelationByIds(f.ctx, f.currentUId, videosIds)
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
	videoList := make([]*feed.Video, 0)
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
		videoList = append(videoList, &feed.Video{
			VideoId: video.VideoId,
			Author: &user.User{
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

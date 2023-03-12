package service_favourite

import (
	"aweme_kitex/cmd/favourite/kitex_gen/favourite"
	"aweme_kitex/cmd/favourite/kitex_gen/feed"
	"aweme_kitex/cmd/favourite/kitex_gen/user"
	favRPC "aweme_kitex/cmd/favourite/rpc"
	favouriteDB "aweme_kitex/cmd/favourite/service_favourite/db"
	feed2 "aweme_kitex/cmd/feed/kitex_gen/feed"
	relation2 "aweme_kitex/cmd/relation/kitex_gen/relation"
	user2 "aweme_kitex/cmd/user/kitex_gen/user"
	"aweme_kitex/pkg/types"
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

func (s *FavoriteListService) FavoriteList(req *favourite.FavouriteListRequest) (map[string]*feed.Video, error) {
	videos, err := FavouriteListService(s.ctx, req.UserId)
	if err != nil {
		return nil, err
	}
	return videos, nil
}

// ---------FavouriteList Service---------------
func FavouriteListService(ctx context.Context, userId string) (map[string]*feed.Video, error) {
	return newFavouriteListDataFlow(ctx, userId).do()
}

type favouriteListDataFlow struct {
	currentUId string

	favours map[string]*feed.Video

	videoRawData []*feed2.Video
	users        map[string]*user2.User
	favoursMap   map[string]*types.FavouriteRaw
	RelationMap  map[string]bool
	ctx          context.Context
}

func newFavouriteListDataFlow(ctx context.Context, id string) *favouriteListDataFlow {
	return &favouriteListDataFlow{
		currentUId: id,
		ctx:        ctx,
	}
}

func (f *favouriteListDataFlow) do() (map[string]*feed.Video, error) {
	if err := f.prepareVideoInfo(); err != nil {
		return nil, err
	}
	if err := f.packageVideo(); err != nil {
		return nil, err
	}
	return f.favours, nil
}

func (f *favouriteListDataFlow) prepareVideoInfo() error {
	videosIds, err := favouriteDB.NewFavouriteDaoInstance().QueryFavoursVideoIdByUid(f.ctx, f.currentUId)
	if err != nil {
		return err
	}
	// 1.get videos
	videoData, err := favRPC.GetVideosById(f.ctx, &feed2.CheckVideoInvalidRequest{VideoId: videosIds})
	if err != nil {
		return err
	}
	f.videoRawData = videoData

	uids := []string{}
	for _, video := range f.videoRawData {
		uids = append(uids, video.Author.UserId)
	}

	// 2.get video authors & isfollow
	users, err := favRPC.GetUserInfo(f.ctx, &user2.SingleUserInfoRequest{UserIds: uids})
	if err != nil {
		return err
	}
	userMap := make(map[string]*user2.User)
	relationMap := make(map[string]bool, len(users))

	var wg sync.WaitGroup
	wg.Add(2)
	var favErr error
	go func() {
		defer wg.Done()
		for _, user := range users {
			userMap[user.UserId] = user
			req := &relation2.QueryRelationRequest{
				UserId:   f.currentUId,
				ToUserId: user.UserId,
				IsFollow: false,
			}
			isfollow, _ := favRPC.QueryRelation(f.ctx, req)
			relationMap[user.UserId] = isfollow
		}
		f.users = userMap
		f.RelationMap = relationMap
	}()
	// 3.get fav
	go func() {
		defer wg.Done()
		favoursMap, err := favouriteDB.NewFavouriteDaoInstance().QueryFavoursByIds(f.ctx, f.currentUId, videosIds)
		if err != nil {
			favErr = err
		}
		f.favoursMap = favoursMap
	}()
	wg.Wait()
	if favErr != nil {
		return favErr
	}
	return nil

}

func (f *favouriteListDataFlow) packageVideo() error {
	videoList := make(map[string]*feed.Video, 0)
	for _, video := range f.videoRawData {
		author, ok := f.users[video.Author.UserId]
		if !ok {
			return errors.New("has no video user info for " + fmt.Sprint(video.Author.UserId))
		}
		var isFavour bool = false
		var isFollow bool = false
		_, ok = f.favoursMap[video.VideoId]
		if ok {
			isFavour = true
		}
		_, ok = f.RelationMap[video.Author.UserId]
		if ok {
			isFollow = true
		}
		videoList[video.VideoId] = &feed.Video{
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
		}
	}
	f.favours = videoList
	return nil
}

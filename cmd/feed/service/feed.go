package service

import (
	feed "aweme_kitex/cmd/feed/kitex_gen/feed"
	"aweme_kitex/controller"
	"aweme_kitex/models"
	"aweme_kitex/models/dal"
	serviceRPC "aweme_kitex/service/rpc"
	"context"
	"sync"
)

type FeedService struct {
	ctx context.Context
}

func NewFeedService(ctx context.Context) *FeedService {
	return &FeedService{ctx: ctx}
}
func (f *FeedService) Feed(req *feed.FeedRequest) ([]*feed.Video, int64, error) {
	user, err := controller.CheckToken(req.Token)
	if err != nil {
		return nil, 0, err
	}
	videoData, err := dal.NewVideoDaoInstance().QueryVideoByLatestTime(f.ctx, req.LatestTime)
	if err != nil {
		return nil, 0, err
	}
	videoIds := make([]string, 0)
	userIds := make([]string, 0)
	for _, video := range videoData {
		videoIds = append(videoIds, video.VideoId)
		userIds = append(userIds, video.UserId)
	}
	users, err := dal.NewUserDaoInstance().QueryUserByIds(f.ctx, userIds)
	if err != nil {
		return nil, 0, err
	}
	userMap := make(map[string]*models.UserRawData)
	for _, user := range users {
		userMap[user.UserId] = user
	}

	var favMap map[string]*models.FavouriteRaw
	var relationMap map[string]*models.RelationRaw
	var favoriteErr, relationErr error
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		favMap, err = dal.NewFavouriteDaoInstance().QueryFavoursByIds(f.ctx, user.Id, videoIds)
		if err != nil {
			favoriteErr = err
			return
		}
	}()

	go func() {
		defer wg.Done()
		relationMap, err = dal.NewRelationDaoInstance().QueryRelationByIds(f.ctx, user.Id, userIds)
		if err != nil {
			relationErr = err
			return
		}
	}()
	wg.Wait()
	if favoriteErr != nil {
		return nil, 0, favoriteErr
	}
	if relationErr != nil {
		return nil, 0, relationErr
	}
	videos, nextTime := serviceRPC.PackRPCVideoInfo(user.Id, videoData, userMap, favMap, relationMap)
	return videos, nextTime, nil
}

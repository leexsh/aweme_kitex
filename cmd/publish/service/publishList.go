package service

import (
	"aweme_kitex/cmd/publish/kitex_gen/feed"
	"aweme_kitex/cmd/publish/kitex_gen/publish"
	"aweme_kitex/controller"
	"aweme_kitex/models"
	"aweme_kitex/models/dal"
	serviceRPC "aweme_kitex/service/rpc"
	"sync"

	"github.com/pkg/errors"
)

func (s *PublishService) PublishList(req *publish.PublishListRequest) ([]*feed.Video, error) {
	user, err := controller.CheckToken(req.Token)
	if err != nil {
		return nil, err
	}

	videoData, err := dal.NewVideoDaoInstance().QueryVideosByUserId(s.ctx, user.Id)
	if err != nil {
		return nil, err
	}

	videoIds := make([]string, 0)
	userIds := []string{user.Id}
	for _, video := range videoData {
		videoIds = append(videoIds, video.UserId)
	}

	users, err := dal.NewUserDaoInstance().QueryUserByIds(s.ctx, userIds)
	if err != nil {
		return nil, err
	}
	if len(users) == 0 {
		return nil, errors.New("user not exist")
	}
	userMap := make(map[string]*models.UserRawData)
	for _, user := range users {
		userMap[string(user.UserId)] = user
	}

	var favoriteMap map[string]*models.FavouriteRaw
	var relationMap map[string]*models.RelationRaw
	var wg sync.WaitGroup
	wg.Add(2)
	var favoriteErr, relationErr error
	// 获取点赞信息
	go func() {
		defer wg.Done()
		favoriteMap, err = dal.NewFavouriteDaoInstance().QueryFavoursByIds(s.ctx, user.Id, videoIds)
		if err != nil {
			favoriteErr = err
			return
		}
	}()
	// 获取关注信息
	go func() {
		defer wg.Done()
		relationMap, err = dal.NewRelationDaoInstance().QueryRelationByIds(s.ctx, user.Id, userIds)
		if err != nil {
			relationErr = err
			return
		}
	}()
	wg.Wait()
	if favoriteErr != nil {
		return nil, favoriteErr
	}
	if relationErr != nil {
		return nil, relationErr
	}
	videoList := serviceRPC.PublishInfo(user.Id, videoData, userMap, favoriteMap, relationMap)
	return videoList, nil
}

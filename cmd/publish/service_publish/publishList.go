package service_publish

import (
	"aweme_kitex/cmd/favourite/kitex_gen/favourite"
	feed2 "aweme_kitex/cmd/feed/kitex_gen/feed"
	"aweme_kitex/cmd/publish/kitex_gen/feed"
	"aweme_kitex/cmd/publish/kitex_gen/publish"
	"aweme_kitex/cmd/publish/kitex_gen/user"
	publishRPC "aweme_kitex/cmd/publish/rpc"
	relation2 "aweme_kitex/cmd/relation/kitex_gen/relation"
	user2 "aweme_kitex/cmd/user/kitex_gen/user"
	"aweme_kitex/pkg/jwt"
	"context"
	"errors"
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
	_, err := jwt.AnalyzeToken(req.Token)
	if err != nil {
		return nil, err
	}
	return queryUserVideos(s.ctx, req.UserId)
}

func queryUserVideos(ctx context.Context, userId string) ([]*feed.Video, error) {
	return newQueryUserVideoList(ctx, userId).do()
}

type userVideoList struct {
	ctx    context.Context
	UserId string

	VideoList    []*feed.Video
	VideoData    []*feed2.Video
	PublishUser  *user2.User
	FavouriteMap map[string]bool
	isFollow     bool
}

func newQueryUserVideoList(ctx context.Context, userId string) *userVideoList {
	return &userVideoList{
		ctx:    ctx,
		UserId: userId,
	}
}

func (f *userVideoList) prepareVideoInfo() error {
	videos, err := publishRPC.GetVideosByUserId(f.ctx, &feed2.GetVideoByUserIDRequest{UserId: f.UserId})
	if err != nil {
		return err
	}
	if len(videos) <= 0 {
		return errors.New("videos is nil")
	}
	f.VideoData = videos

	videoIds := make([]string, 0)
	userIds := []string{f.UserId}
	for _, video := range f.VideoData {
		videoIds = append(videoIds, video.VideoId)
	}
	publishUser, err := publishRPC.GetUserInfo(f.ctx, &user2.SingleUserInfoRequest{UserIds: userIds})
	if err != nil {
		return err
	}
	if len(publishUser) > 0 {
		f.PublishUser = publishUser[userIds[0]]
	}
	// 查看别人发布的作品， 查看自己发布的作品
	var wg sync.WaitGroup
	wg.Add(2)
	var relationErr, favourErr error
	// 3. get relation
	go func() {
		defer wg.Done()
		req := &relation2.QueryRelationRequest{
			UserId:   f.UserId,
			ToUserId: userIds[0],
			IsFollow: false,
		}
		isfollow, err := publishRPC.QueryRelation(f.ctx, req)
		if err != nil {
			relationErr = err
		}
		f.isFollow = isfollow
	}()

	// 5.获取点赞信息
	go func() {
		defer wg.Done()
		favList, err := publishRPC.QueryIsFavourite(f.ctx, &favourite.QueryVideoIsFavouriteRequest{VideosId: videoIds, UserId: f.UserId})
		if err != nil {
			favourErr = err
			return
		}
		f.FavouriteMap = favList
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

func (f *userVideoList) packVideoInfo() error {
	videoList := make([]*feed.Video, 0)
	for _, video := range f.VideoData {
		var isFavorite bool = false

		_, ok := f.FavouriteMap[video.VideoId]
		if ok {
			isFavorite = true
		}
		curVideo := &feed.Video{
			VideoId: video.VideoId,
			Author: &user.User{
				UserId:        f.PublishUser.UserId,
				Name:          f.PublishUser.Name,
				FollowCount:   f.PublishUser.FollowCount,
				FollowerCount: f.PublishUser.FollowerCount,
				IsFollow:      f.isFollow,
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

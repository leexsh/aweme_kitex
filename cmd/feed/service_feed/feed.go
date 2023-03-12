package service_feed

import (
	favourite2 "aweme_kitex/cmd/favourite/kitex_gen/favourite"
	feed "aweme_kitex/cmd/feed/kitex_gen/feed"
	"aweme_kitex/cmd/feed/kitex_gen/user"
	feedRPC "aweme_kitex/cmd/feed/rpc"
	videoDB "aweme_kitex/cmd/feed/service_feed/db"
	relation2 "aweme_kitex/cmd/relation/kitex_gen/relation"
	user2 "aweme_kitex/cmd/user/kitex_gen/user"
	"aweme_kitex/pkg/types"
	"context"
	"errors"
	"fmt"
	"sync"
)

type FeedService struct {
	ctx context.Context
}

func NewFeedService(ctx context.Context) *FeedService {
	return &FeedService{ctx: ctx}
}
func (f *FeedService) Feed(req *feed.FeedRequest) ([]*feed.Video, int64, error) {
	videoList, nextTime, err := queryVideoData(f.ctx, req.LatestTime, req.UserId)
	if err != nil {
		return nil, 0, err
	}

	return videoList, nextTime, nil
}

func queryVideoData(ctx context.Context, latestTime int64, userId string) ([]*feed.Video, int64, error) {
	return newQueryVideoDataFlow(ctx, latestTime, userId).Do()
}

func newQueryVideoDataFlow(ctx context.Context, latestTime int64, userId string) *queryVideoDataFlow {
	return &queryVideoDataFlow{
		LatestTime:    latestTime,
		CurrentUserId: userId,
		ctx:           ctx,
	}
}

// video data
type queryVideoDataFlow struct {
	ctx           context.Context
	CurrentUserId string

	LatestTime int64
	VideoList  []*feed.Video
	NextTime   int64

	VideoData   []*types.VideoRawData
	UserMap     map[string]*user2.User
	FavoursMap  map[string]bool
	RelationMap map[string]bool
}

func (f *queryVideoDataFlow) Do() ([]*feed.Video, int64, error) {
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
	videoData, err := videoDB.NewVideoDaoInstance().QueryVideoByLatestTime(f.ctx, f.LatestTime)
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
	var wg sync.WaitGroup
	wg.Add(2)
	var favourErr, relationErr, userErr error
	userMap := make(map[string]*user2.User)
	relationMap := make(map[string]bool, len(authorIds))
	// 3. get user info
	go func() {
		defer wg.Done()
		users, err := feedRPC.GetUserInfo(f.ctx, &user2.SingleUserInfoRequest{UserIds: authorIds})
		if err != nil {
			userErr = err
		}
		for _, user := range users {
			userMap[user.UserId] = user
			req := &relation2.QueryRelationRequest{
				UserId:   f.CurrentUserId,
				ToUserId: user.UserId,
				IsFollow: false,
			}
			isfollow, err := feedRPC.QueryRelation(f.ctx, req)
			if err != nil {
				relationErr = err
			}
			relationMap[user.UserId] = isfollow
		}
		f.UserMap = userMap
		f.RelationMap = relationMap
	}()

	// 5.获取点赞信息
	go func() {
		defer wg.Done()
		favList, err := feedRPC.QueryIsFavourite(f.ctx, &favourite2.QueryVideoIsFavouriteRequest{VideosId: videoIds, UserId: f.CurrentUserId})
		if err != nil {
			favourErr = err
			return
		}
		f.FavoursMap = favList
	}()
	wg.Wait()

	if favourErr != nil {
		return favourErr
	}
	if relationErr != nil {
		return relationErr
	}
	if userErr != nil {
		return userErr
	}
	return nil
}

func (f *queryVideoDataFlow) packVideoInfo() error {
	videoList := make([]*feed.Video, 0)
	for _, video := range f.VideoData {
		videoAuthor, ok := f.UserMap[video.UserId]
		if !ok {
			return errors.New("has no video user info for " + fmt.Sprint(video.UserId))
		}
		var isFavourite = false
		var isFollow = false
		if f.CurrentUserId != "" {
			if _, ok := f.FavoursMap[video.VideoId]; ok {
				isFavourite = true
			}
			if _, ok := f.RelationMap[video.UserId]; ok {
				isFollow = true
			}
		}
		video := &feed.Video{
			VideoId: video.VideoId,
			Author: &user.User{
				UserId:        videoAuthor.UserId,
				Name:          videoAuthor.Name,
				FollowCount:   videoAuthor.FollowCount,
				FollowerCount: videoAuthor.FollowerCount,
				IsFollow:      isFollow,
			},
			PlayUrl:        video.PlayUrl,
			CoverUrl:       video.CoverUrl,
			FavouriteCount: video.FavouriteCount,
			CommentCount:   video.CommentCount,
			IsFavourite:    isFavourite,
			Title:          video.Title,
		}
		videoList = append(videoList, video)
	}
	f.VideoList = videoList
	if len(f.VideoList) == 0 {
		f.NextTime = 0
	} else {
		f.NextTime = f.VideoData[len(f.VideoData)-1].CreatedAt.Unix()
	}
	return nil
}

package service_feed

import (
	feed "aweme_kitex/cmd/feed/kitex_gen/feed"
	"aweme_kitex/cmd/feed/kitex_gen/user"
	db3 "aweme_kitex/cmd/feed/service_feed/db"
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

type FeedService struct {
	ctx context.Context
}

func NewFeedService(ctx context.Context) *FeedService {
	return &FeedService{ctx: ctx}
}
func (f *FeedService) Feed(req *feed.FeedRequest) ([]*feed.Video, int64, error) {
	uc, _ := jwt.AnalyzeToken(req.Token)
	videoList, nextTime, err := queryVideoData(f.ctx, req.LatestTime, uc.Id)
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

	VideoData   []*models.VideoRawData
	UserMap     map[string]*models.UserRawData
	FavoursMap  map[string]*models.FavouriteRaw
	RelationMap map[string]*models.RelationRaw
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
	videoData, err := db3.NewVideoDaoInstance().QueryVideoByLatestTime(f.ctx, f.LatestTime)
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

	// 3. get user info
	users, err := db2.NewUserDaoInstance().QueryUserByIds(f.ctx, authorIds)
	if err != nil {
		return err
	}
	userMap := make(map[string]*models.UserRawData)
	for _, user := range users {
		userMap[user.UserId] = user
	}
	f.UserMap = userMap

	// 4. should login
	if f.CurrentUserId == "" {
		return nil
	}

	var wg sync.WaitGroup
	wg.Add(2)
	var favourErr, relationErr error
	// 5.获取点赞信息
	go func() {
		defer wg.Done()
		favoursMap, err := dal.NewFavouriteDaoInstance().QueryFavoursByIds(f.ctx, f.CurrentUserId, videoIds)
		if err != nil {
			favourErr = err
			return
		}
		f.FavoursMap = favoursMap
	}()
	// 6.获取关注信息
	go func() {
		defer wg.Done()
		relationMap, err := db.NewRelationDaoInstance().QueryRelationByIds(f.ctx, f.CurrentUserId, authorIds)
		if err != nil {
			relationErr = err
			return
		}
		f.RelationMap = relationMap

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

func (f *queryVideoDataFlow) packVideoInfo() error {
	videoList := make([]*feed.Video, 0)
	for _, video := range f.VideoData {
		videoAuthor, ok := f.UserMap[video.UserId]
		if !ok {
			return errors.New("has no video user info for " + fmt.Sprint(video.UserId))
		}
		var isFavourite bool = false
		var isFollow bool = false
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

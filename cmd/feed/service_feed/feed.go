package service_feed

import (
	feed "aweme_kitex/cmd/feed/kitex_gen/feed"
	"aweme_kitex/cmd/feed/kitex_gen/user"
	"aweme_kitex/models"
	"aweme_kitex/pkg/jwt"
	"aweme_kitex/service"
	"context"
)

type FeedService struct {
	ctx context.Context
}

func NewFeedService(ctx context.Context) *FeedService {
	return &FeedService{ctx: ctx}
}
func (f *FeedService) Feed(req *feed.FeedRequest) ([]*feed.Video, int64, error) {
	uc, _ := jwt.AnalyzeToken(req.Token)
	videoList, nextTime, err := service.QueryVideoData(req.LatestTime, uc.Id)
	if err != nil {
		return nil, 0, err
	}

	return f.packVideoInfo(videoList), nextTime, nil
}
func (f *FeedService) packVideoInfo(videos []*models.Video) []*feed.Video {
	videoList := make([]*feed.Video, 0)
	for _, v := range videos {
		video := &feed.Video{
			VideoId: v.Id,
			Author: &user.User{
				UserId:        v.Author.UserId,
				Name:          v.Author.Name,
				FollowCount:   v.Author.FollowCount,
				FollowerCount: v.Author.FollowerCount,
				IsFollow:      v.Author.IsFollow,
			},
			PlayUrl:        v.PlayUrl,
			CoverUrl:       v.CoverUrl,
			FavouriteCount: v.FavouriteCount,
			CommentCount:   v.CommentCount,
			IsFavourite:    v.IsFavourite,
			Title:          v.Title,
		}
		videoList = append(videoList, video)
	}
	return videoList
}

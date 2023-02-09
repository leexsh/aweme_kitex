package service_publish

import (
	"aweme_kitex/cmd/publish/kitex_gen/feed"
	"aweme_kitex/cmd/publish/kitex_gen/publish"
	"aweme_kitex/cmd/publish/kitex_gen/user"
	"aweme_kitex/models"
	"aweme_kitex/pkg/jwt"
	"aweme_kitex/service"
	"context"
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
	videos, err := service.QueryUserVideos(uc.Id)
	if err != nil {
		return nil, err
	}
	return s.packVideoInfo(videos), nil
}

func (s *PublishListService) packVideoInfo(videos []*models.Video) []*feed.Video {
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

package service_favourite

import (
	"aweme_kitex/cmd/favourite/kitex_gen/favourite"
	"aweme_kitex/cmd/favourite/kitex_gen/feed"
	"aweme_kitex/cmd/favourite/kitex_gen/user"
	"aweme_kitex/pkg/jwt"
	"aweme_kitex/service"
	"context"
)

type FavoriteListService struct {
	ctx context.Context
}

// NewFavoriteListService new FavoriteListService
func NewFavoriteListService(ctx context.Context) *FavoriteListService {
	return &FavoriteListService{ctx: ctx}
}

func (s *FavoriteListService) FavoriteList(req *favourite.FavouriteListRequest) ([]*feed.Video, error) {
	uc, err := s.CheckToken(req.Token)
	if err != nil {
		return nil, err
	}
	videos, err := service.FavouriteListService(uc.Id, uc.Name)
	if err != nil {
		return nil, err
	}
	videoList := make([]*feed.Video, 0)
	for _, v := range videos {
		fvideo := &feed.Video{
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
		videoList = append(videoList, fvideo)
	}
	return videoList, nil
}

func (s *FavoriteListService) CheckToken(token string) (*jwt.UserClaim, error) {
	uc, err := jwt.AnalyzeToken(token)
	if err != nil {
		return nil, err
	}
	return uc, nil
}

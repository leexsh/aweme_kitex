package service_favourite

import (
	"aweme_kitex/cmd/favourite/kitex_gen/favourite"
	"aweme_kitex/pkg/jwt"
	"aweme_kitex/service"
	"context"
)

type FavoriteActionService struct {
	ctx context.Context
}

// NewFavoriteActionService new FavoriteActionService
func NewFavoriteActionService(ctx context.Context) *FavoriteActionService {
	return &FavoriteActionService{ctx: ctx}
}

func (s *FavoriteActionService) FavoriteAction(req *favourite.FavouriteActionRequest) error {
	uc, err := s.CheckToken(req.Token)
	if err != nil {
		return err
	}
	err = service.FavouriteActionService(uc, req.VideoId, req.ActionType)
	if err != nil {
		return err
	}
	return nil
}

func (s *FavoriteActionService) CheckToken(token string) (*jwt.UserClaim, error) {
	uc, err := jwt.AnalyzeToken(token)
	if err != nil {
		return nil, err
	}
	return uc, nil
}

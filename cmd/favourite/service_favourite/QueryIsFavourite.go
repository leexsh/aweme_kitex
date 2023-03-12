package service_favourite

import (
	"aweme_kitex/cmd/favourite/kitex_gen/favourite"
	favouriteDB "aweme_kitex/cmd/favourite/service_favourite/db"
	"context"
)

type QueryIsFavouriteService struct {
	ctx context.Context
}

// NewFavoriteListService new FavoriteListService
func NewQueryIsFavouriteService(ctx context.Context) *QueryIsFavouriteService {
	return &QueryIsFavouriteService{ctx: ctx}
}

func (s *QueryIsFavouriteService) QueryIsFavours(req *favourite.QueryVideoIsFavouriteRequest) (map[string]bool, error) {
	isFavours, err := favouriteDB.NewFavouriteDaoInstance().QueryIsFavours(s.ctx, req.UserId, req.VideosId)
	if err != nil {
		return nil, err
	}
	return isFavours, nil
}

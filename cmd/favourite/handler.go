package main

import (
	favourite "aweme_kitex/cmd/favourite/kitex_gen/favourite"
	"aweme_kitex/cmd/favourite/service_favourite"
	favPack "aweme_kitex/cmd/favourite/service_favourite/pack"
	"context"
)

// FavouriteServiceImpl implements the last service_user interface defined in the IDL.
type FavouriteServiceImpl struct{}

// FavouriteAction implements the FavouriteServiceImpl interface.
func (s *FavouriteServiceImpl) FavouriteAction(ctx context.Context, req *favourite.FavouriteActionRequest) (resp *favourite.FavouriteActionResponse, err error) {

	if len(req.VideoId) == 0 || len(req.ActionType) == 0 {
		return favPack.FavouriteActionResponse(-1, "favourite action params error"), nil
	}

	err = service_favourite.NewFavoriteActionService(ctx).FavoriteAction(req)
	if err != nil {
		return favPack.FavouriteActionResponse(-1, err.Error()), nil
	}
	return favPack.FavouriteActionResponse(-1, "favourite action success"), nil
}

// FavouriteList implements the FavouriteServiceImpl interface.
func (s *FavouriteServiceImpl) FavouriteList(ctx context.Context, req *favourite.FavouriteListRequest) (resp *favourite.FavouriteListResponse, err error) {
	videoList, err := service_favourite.NewFavoriteListService(ctx).FavoriteList(req)
	if err != nil {
		return favPack.FavouriteListResponse(-1, err.Error(), nil), nil
	}

	return favPack.FavouriteListResponse(0, "favourite action success", videoList), nil
}

// QueryVideoIsFavourite implements the FavouriteServiceImpl interface.
func (s *FavouriteServiceImpl) QueryVideoIsFavourite(ctx context.Context, req *favourite.QueryVideoIsFavouriteRequest) (resp *favourite.QueryVideoIsFavouriteResponse, err error) {
	isFavours, err := service_favourite.NewQueryIsFavouriteService(ctx).QueryIsFavours(req)
	if err != nil {
		return favPack.QueryFavoursResponse(-1, err.Error(), nil), nil
	}

	return favPack.QueryFavoursResponse(0, "favourite action success", isFavours), nil
	return
}

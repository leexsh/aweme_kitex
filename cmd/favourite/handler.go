package main

import (
	favourite "aweme_kitex/cmd/favourite/kitex_gen/favourite"
	"aweme_kitex/cmd/favourite/service_favourite"
	"context"
	"time"
)

// FavouriteServiceImpl implements the last service_user interface defined in the IDL.
type FavouriteServiceImpl struct{}

// FavouriteAction implements the FavouriteServiceImpl interface.
func (s *FavouriteServiceImpl) FavouriteAction(ctx context.Context, req *favourite.FavouriteActionRequest) (resp *favourite.FavouriteActionResponse, err error) {
	// TODO: Your code here...
	resp = new(favourite.FavouriteActionResponse)

	if len(req.Token) == 0 || len(req.VideoId) == 0 || len(req.ActionType) == 0 {
		resp.BaseResp.ServiceTime = time.Now().Unix()
		resp.BaseResp.StatusCode = -1
		resp.BaseResp.StatusMsg = "favourite action params error"
		return resp, nil
	}

	err = service_favourite.NewFavoriteActionService(ctx).FavoriteAction(req)
	if err != nil {
		resp.BaseResp.ServiceTime = time.Now().Unix()
		resp.BaseResp.StatusCode = -1
		resp.BaseResp.StatusMsg = err.Error()
		return resp, nil
	}
	resp.BaseResp.ServiceTime = time.Now().Unix()
	resp.BaseResp.StatusMsg = "favourite action success"
	resp.BaseResp.StatusCode = 0

	return resp, nil
	return
}

// FavouriteList implements the FavouriteServiceImpl interface.
func (s *FavouriteServiceImpl) FavouriteList(ctx context.Context, req *favourite.FavouriteListRequest) (resp *favourite.FavouriteListResponse, err error) {
	resp = new(favourite.FavouriteListResponse)

	if len(req.Token) == 0 {
		resp.BaseResp.ServiceTime = time.Now().Unix()
		resp.BaseResp.StatusCode = -1
		resp.BaseResp.StatusMsg = "favourite list params error"
		return resp, nil
	}

	videoList, err := service_favourite.NewFavoriteListService(ctx).FavoriteList(req)
	if err != nil {
		resp.BaseResp.ServiceTime = time.Now().Unix()
		resp.BaseResp.StatusCode = -1
		resp.BaseResp.StatusMsg = err.Error()
		return resp, nil
	}

	resp.BaseResp.ServiceTime = time.Now().Unix()
	resp.BaseResp.StatusMsg = "favourite action success"
	resp.BaseResp.StatusCode = 0
	resp.VideoList = videoList
	return resp, nil
}

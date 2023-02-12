package service_favourite

import (
	"aweme_kitex/cmd/favourite/kitex_gen/favourite"
	"aweme_kitex/models"
	"aweme_kitex/models/dal"
	"aweme_kitex/pkg/jwt"
	"aweme_kitex/utils"
	"context"
	"errors"
)

type FavoriteActionService struct {
	ctx context.Context
}

// NewFavoriteActionService new FavoriteActionService
func NewFavoriteActionService(ctx context.Context) *FavoriteActionService {
	return &FavoriteActionService{ctx: ctx}
}

func (s *FavoriteActionService) FavoriteAction(req *favourite.FavouriteActionRequest) error {
	uc, err := jwt.AnalyzeToken(req.Token)
	if err != nil {
		return err
	}
	err = favouriteActionHandler(s.ctx, uc.Id, req.VideoId, req.ActionType)
	if err != nil {
		return err
	}
	return nil
}

// ---------FavouriteAction Service---------------
func favouriteActionHandler(ctx context.Context, userId, videoId, action string) error {
	return newFavouriteActionData(ctx, userId, videoId, action).do()
}

type favouriteActionDataFlow struct {
	CurrentUid string
	videoId    string
	action     string

	FavouriteData *models.FavouriteRaw

	ctx context.Context
}

func newFavouriteActionData(ctx context.Context, userId, videoId, action string) *favouriteActionDataFlow {
	return &favouriteActionDataFlow{
		CurrentUid: userId,
		videoId:    videoId,
		action:     action,
		ctx:        ctx,
	}
}

func (f *favouriteActionDataFlow) do() error {
	if _, err := dal.NewVideoDaoInstance().CheckVideoId([]string{f.videoId}); err != nil {
		return err
	}
	if f.action != "1" && f.action != "2" {
		return errors.New("invalid action type")
	}
	if f.action == "1" {
		favour := &models.FavouriteRaw{
			Id:      utils.GenerateUUID(),
			UserId:  f.CurrentUid,
			VideoId: f.videoId,
		}
		f.FavouriteData = favour
		err := dal.NewFavouriteDaoInstance().CreateFavour(f.ctx, favour, f.videoId)
		if err != nil {
			return err
		}
	} else if f.action == "2" {
		err := dal.NewFavouriteDaoInstance().DelFavour(f.ctx, f.CurrentUid, f.videoId)
		if err != nil {
			return err
		}
	}
	return nil
}

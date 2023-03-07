package service_favourite

import (
	"aweme_kitex/cmd/favourite/kitex_gen/favourite"
	favRPC "aweme_kitex/cmd/favourite/rpc"
	favKafka "aweme_kitex/cmd/favourite/service_favourite/kafka"
	"aweme_kitex/cmd/feed/kitex_gen/feed"
	constants "aweme_kitex/pkg/constant"
	"aweme_kitex/pkg/jwt"
	"aweme_kitex/pkg/utils"
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
	err := favRPC.CheckVideoInvalid(f.ctx, &feed.CheckVideoInvalidRequest{VideoId: []string{f.videoId}})
	if err != nil {
		return err
	}
	if f.action != constants.Like && f.action != constants.Unlike {
		return errors.New("invalid action type")
	}
	if f.action == constants.Like {
		msg := utils.GenerateUUID() + "&" + f.CurrentUid + "&" + f.videoId
		favKafka.ProduceAddRelation(constants.KafKaFavouriteAddTopic, msg)
	} else if f.action == constants.Unlike {
		msg := f.CurrentUid + "&" + f.videoId
		favKafka.ProduceAddRelation(constants.KafKaFavouriteDelTopic, msg)
	}
	return nil
}

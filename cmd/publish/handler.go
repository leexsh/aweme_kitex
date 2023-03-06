package main

import (
	publish "aweme_kitex/cmd/publish/kitex_gen/publish"
	"aweme_kitex/cmd/publish/service_publish"
	publishPack "aweme_kitex/cmd/publish/service_publish/pack"
	"context"
)

// PublishServiceImpl implements the last service_user interface defined in the IDL.
type PublishServiceImpl struct{}

// PublishAction implements the PublishServiceImpl interface.
func (s *PublishServiceImpl) PublishAction(ctx context.Context, req *publish.PublishActionRequest) (resp *publish.PublishActionResponse, err error) {
	if len(req.Title) <= 0 || req.Data == nil {
		return publishPack.PackPublishAction(-1, "param error"), nil
	}
	err = service_publish.NewPublishService(ctx).Publish(req)
	if err != nil {
		return publishPack.PackPublishAction(-1, err.Error()), nil
	}
	return publishPack.PackPublishAction(0, "publish success"), nil
}

// PublishList implements the PublishServiceImpl interface.
func (s *PublishServiceImpl) PublishList(ctx context.Context, req *publish.PublishListRequest) (resp *publish.PublishListResponse, err error) {
	videoList, err := service_publish.NewPublishListService(ctx).PublishList(req)
	if err != nil {
		return publishPack.PackPublishList(-1, err.Error(), nil), nil
	}
	return publishPack.PackPublishList(0, "get publish List success", videoList), nil
}

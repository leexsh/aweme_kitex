package main

import (
	"aweme_kitex/cmd/publish/kitex_gen/base"
	publish "aweme_kitex/cmd/publish/kitex_gen/publish"
	"aweme_kitex/cmd/publish/service_publish"
	"context"
	"time"
)

// PublishServiceImpl implements the last service_user interface defined in the IDL.
type PublishServiceImpl struct{}

// PublishAction implements the PublishServiceImpl interface.
func (s *PublishServiceImpl) PublishAction(ctx context.Context, req *publish.PublishActionRequest) (resp *publish.PublishActionResponse, err error) {
	resp = new(publish.PublishActionResponse)
	if len(req.Title) <= 0 || req.Data == nil {
		resp.BaseResp = &base.BaseResp{
			StatusCode:  -1,
			StatusMsg:   "param error",
			ServiceTime: time.Now().Unix(),
		}
		return resp, nil
	}
	err = service_publish.NewPublishService(ctx).Publish(req)
	if err != nil {
		resp.BaseResp.StatusCode = -1
		resp.BaseResp.StatusMsg = err.Error()
		return resp, err
	}
	resp.BaseResp.StatusCode = 0
	resp.BaseResp.StatusMsg = "success"
	return resp, nil
}

// PublishList implements the PublishServiceImpl interface.
func (s *PublishServiceImpl) PublishList(ctx context.Context, req *publish.PublishListRequest) (resp *publish.PublishListResponse, err error) {
	resp = new(publish.PublishListResponse)
	videoList, err := service_publish.NewPublishListService(ctx).PublishList(req)
	if err != nil {
		resp.BaseResp.StatusCode = -1
		resp.BaseResp.StatusMsg = err.Error()
		resp.BaseResp.ServiceTime = time.Now().Unix()
		return resp, err
	}
	resp.BaseResp.StatusCode = 0
	resp.BaseResp.StatusMsg = "get publish List success"
	resp.BaseResp.ServiceTime = time.Now().Unix()
	resp.VideoList = videoList
	return resp, nil
}

package main

import (
	feed "aweme_kitex/cmd/feed/kitex_gen/feed"
	"aweme_kitex/cmd/feed/service_feed"
	"aweme_kitex/models"
	"aweme_kitex/pkg/errno"
	"context"
)

// FeedServiceImpl implements the last service_user interface defined in the IDL.
type FeedServiceImpl struct{}

// Feed implements the FeedServiceImpl interface.
func (s *FeedServiceImpl) Feed(ctx context.Context, req *feed.FeedRequest) (resp *feed.FeedResponse, err error) {
	resp = new(feed.FeedResponse)
	if req.LatestTime <= 0 {
		resp.BaseResp = models.BuildBaseResp(errno.ParamErr)
		return resp, nil
	}
	videos, nextTime, err := service_feed.NewFeedService(ctx).Feed(req)
	if err != nil {
		resp.BaseResp = models.BuildBaseResp(err)
		return resp, nil
	}
	resp.BaseResp = models.BuildBaseResp(errno.Success)
	resp.VideoList = videos
	resp.NextTime = nextTime
	return
}

package main

import (
	"aweme_kitex/cmd/feed/kitex_gen/base"
	feed "aweme_kitex/cmd/feed/kitex_gen/feed"
	"aweme_kitex/cmd/feed/service_feed"
	"aweme_kitex/pkg/errno"
	"context"
	"errors"
	"time"
)

// FeedServiceImpl implements the last service_user interface defined in the IDL.
type FeedServiceImpl struct{}

func BuildBaseResp(err error) *base.BaseResp {
	if err == nil {
		return baseResp(errno.Success)
	}
	e := errno.ErrMsg{}
	if errors.As(err, &e) {
		return baseResp(e)
	}
	s := errno.ServiceErr.WithMessage(err.Error())
	return baseResp(s)
}

func baseResp(err errno.ErrMsg) *base.BaseResp {
	return &base.BaseResp{
		StatusCode:  err.ErrCode,
		StatusMsg:   err.ErrMsg,
		ServiceTime: time.Now().Unix(),
	}
}

// Feed implements the FeedServiceImpl interface.
func (s *FeedServiceImpl) Feed(ctx context.Context, req *feed.FeedRequest) (resp *feed.FeedResponse, err error) {
	resp = new(feed.FeedResponse)
	if req.LatestTime <= 0 {
		resp.BaseResp = BuildBaseResp(errno.ParamErr)
		return resp, nil
	}
	videos, nextTime, err := service_feed.NewFeedService(ctx).Feed(req)
	if err != nil {
		resp.BaseResp = BuildBaseResp(err)
		return resp, nil
	}
	resp.BaseResp = BuildBaseResp(errno.Success)
	resp.VideoList = videos
	resp.NextTime = nextTime
	return
}

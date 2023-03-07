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

// ChangeCommentCnt implements the FeedServiceImpl interface.
func (s *FeedServiceImpl) ChangeCommentCnt(ctx context.Context, req *feed.ChangeCommentCountRequest) (resp *feed.ChangeCommentCountResponse, err error) {
	resp = new(feed.ChangeCommentCountResponse)
	err = service_feed.NewChangeCommentCountService(ctx).ChangeCommentCount(req.VideoId, req.Action)
	if err != nil {
		resp.BaseResp.StatusCode = -1
	} else {
		resp.BaseResp.StatusCode = 0
	}
	resp.BaseResp.StatusMsg = err.Error()
	resp.BaseResp.ServiceTime = time.Now().Unix()
	return
}

// CheckVideoInvalid implements the FeedServiceImpl interface.
func (s *FeedServiceImpl) CheckVideoInvalid(ctx context.Context, req *feed.CheckVideoInvalidRequest) (resp *feed.CheckVideoInvalidResponse, err error) {
	resp = new(feed.CheckVideoInvalidResponse)
	err = service_feed.NewCheckVideoService(ctx).CheckVideoInvalid(req.VideoId)
	if err != nil {
		resp.BaseResp.StatusCode = -1
	} else {
		resp.BaseResp.StatusCode = 0
	}
	resp.BaseResp.ServiceTime = time.Now().Unix()
	return
}

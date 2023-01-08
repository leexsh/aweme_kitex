package rpc

import (
	"aweme_kitex/cmd/feed/kitex_gen/feed"
	"aweme_kitex/cmd/feed/kitex_gen/feed/feedservice"
	constants "aweme_kitex/pkg/constant"
	"aweme_kitex/pkg/errno"
	"aweme_kitex/pkg/logger"
	"aweme_kitex/pkg/middleware"
	"context"
	"time"

	client2 "github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/retry"
	etcd "github.com/kitex-contrib/registry-etcd"
	trace "github.com/kitex-contrib/tracer-opentracing"
)

var feedClient feedservice.Client

func initFeedRpc() {
	r, err := etcd.NewEtcdResolver([]string{constants.EtcdAddress})
	if err != nil {
		logger.Error("feed service register error")
		return
	}
	client, err := feedservice.NewClient(
		constants.FeedServiceName,
		client2.WithMiddleware(middleware.CommonMiddleware),
		client2.WithInstanceMW(middleware.ClientMiddleware),
		client2.WithMuxConnection(1),
		client2.WithRPCTimeout(3*time.Second),
		client2.WithConnectTimeout(50*time.Millisecond),
		client2.WithFailureRetry(retry.NewFailurePolicy()),
		client2.WithSuite(trace.NewDefaultClientSuite()),
		client2.WithResolver(r),
	)
	if err != nil {
		logger.Error("feed service register error")
	}
	feedClient = client
}

func Feed(ctx context.Context, req *feed.FeedRequest) ([]*feed.Video, int64, error) {
	resp, err := feedClient.Feed(ctx, req)
	if err != nil {
		return nil, 0, err
	}
	if resp.BaseResp.StatusCode != 0 {
		return nil, 0, errno.NewErr(resp.BaseResp.StatusCode, resp.BaseResp.StatusMsg)
	}
	return resp.VideoList, resp.NextTime, nil
}

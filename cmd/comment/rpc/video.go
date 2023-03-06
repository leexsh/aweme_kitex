package commentRPC

import (
	"aweme_kitex/cmd/feed/kitex_gen/feed"
	"aweme_kitex/cmd/feed/kitex_gen/feed/feedservice"
	constants "aweme_kitex/pkg/constant"
	"aweme_kitex/pkg/logger"
	"aweme_kitex/pkg/middleware"
	"context"
	"errors"
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
		logger.Error("feed service_user register error")
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
		logger.Error("feed service_user register error")
	}
	feedClient = client
}

func ChangeCommentCount(ctx context.Context, req *feed.ChangeCommentCountRequest) error {
	resp, err := feedClient.ChangeCommentCnt(ctx, req)
	if err != nil {
		return err
	}
	if resp.BaseResp.StatusCode != 0 {
		return errors.New(resp.BaseResp.StatusMsg)
	}
	return nil
}

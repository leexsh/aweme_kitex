package rpc

import (
	relation2 "aweme_kitex/cmd/relation/kitex_gen/relation"
	relation "aweme_kitex/cmd/relation/kitex_gen/relation/relationservice"
	"aweme_kitex/cmd/relation/kitex_gen/user"
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

var relationClient relation.Client

func initRelationRpc() {
	r, err := etcd.NewEtcdResolver([]string{constants.EtcdAddress})
	if err != nil {
		logger.Error("feed service register error")
		return
	}
	client, err := relation.NewClient(
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
	relationClient = client
}

// RelationAction implement follow and unfollow actions
func RelationAction(ctx context.Context, req *relation2.RelationActionRequest) error {
	resp, err := relationClient.RelationAction(ctx, req)
	if err != nil {
		return err
	}
	if resp.BaseResp.StatusCode != 0 {
		return errno.NewErr(resp.BaseResp.StatusCode, resp.BaseResp.StatusMsg)
	}
	return nil
}

// FollowList get user follow list info
func FollowList(ctx context.Context, req *relation2.FollowListRequest) ([]*user.User, error) {
	resp, err := relationClient.FollowList(ctx, req)
	if err != nil {
		return nil, err
	}
	if resp.BaseResp.StatusCode != 0 {
		return nil, errno.NewErr(resp.BaseResp.StatusCode, resp.BaseResp.StatusMsg)
	}
	return resp.UserList, nil
}

// FollowerList get user follower list info
func FollowerList(ctx context.Context, req *relation2.FollowerListRequest) ([]*user.User, error) {
	resp, err := relationClient.FollowerList(ctx, req)
	if err != nil {
		return nil, err
	}
	if resp.BaseResp.StatusCode != 0 {
		return nil, errno.NewErr(resp.BaseResp.StatusCode, resp.BaseResp.StatusMsg)
	}
	return resp.UserList, nil
}

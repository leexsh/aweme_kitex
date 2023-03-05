package userRPC

import (
	relation2 "aweme_kitex/cmd/relation/kitex_gen/relation"
	relation "aweme_kitex/cmd/relation/kitex_gen/relation/relationservice"
	constants "aweme_kitex/pkg/constant"
	"aweme_kitex/pkg/errno"
	"aweme_kitex/pkg/middleware"
	"context"
	"time"

	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/retry"
	etcd "github.com/kitex-contrib/registry-etcd"
	trace "github.com/kitex-contrib/tracer-opentracing"
)

func Init() {
	initRelationRpc()
}

var relationClient relation.Client

func initRelationRpc() {
	r, err := etcd.NewEtcdResolver([]string{constants.EtcdAddress})
	if err != nil {
		panic(err)
	}

	c, err := relation.NewClient(
		constants.UserServiceName,
		client.WithMiddleware(middleware.CommonMiddleware),
		client.WithInstanceMW(middleware.ClientMiddleware),
		client.WithMuxConnection(1),
		client.WithRPCTimeout(3*time.Second),
		client.WithConnectTimeout(50*time.Millisecond),
		client.WithFailureRetry(retry.NewFailurePolicy()),
		client.WithSuite(trace.NewDefaultClientSuite()),
		client.WithResolver(r),
	)
	if err != nil {
		panic(err)
	}
	relationClient = c
}

func QueryRelation(ctx context.Context, req *relation2.QueryRelationRequest) (bool, error) {
	resp, err := relationClient.QueryRelation(ctx, req)
	if err != nil {
		return false, err
	}
	if resp.BaseResp.StatusCode != 0 {
		return false, errno.NewErr(resp.BaseResp.StatusCode, resp.BaseResp.StatusMsg)
	}
	return resp.IsFollow, nil
}

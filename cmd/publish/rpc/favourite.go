package publishRPC

import (
	favourite2 "aweme_kitex/cmd/favourite/kitex_gen/favourite"
	favourite "aweme_kitex/cmd/favourite/kitex_gen/favourite/favouriteservice"
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

var favouriteClient favourite.Client

func initFavouriteRpc() {
	r, err := etcd.NewEtcdResolver([]string{constants.EtcdAddress})
	if err != nil {
		logger.Error("feed service_user register error")
		return
	}
	client, err := favourite.NewClient(
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
	favouriteClient = client
}

func QueryIsFavourite(ctx context.Context, request *favourite2.QueryVideoIsFavouriteRequest) (map[string]bool, error) {
	resp, err := favouriteClient.QueryVideoIsFavourite(ctx, request)
	if err != nil {
		return nil, err
	}
	if resp.BaseResp.StatusCode != 0 {
		return nil, errors.New("error")
	}
	return resp.IsFavourites, nil
}

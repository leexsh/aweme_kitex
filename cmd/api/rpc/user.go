package rpc

import (
	"aweme_kitex/cmd/user/kitex_gen/user"
	"aweme_kitex/cmd/user/kitex_gen/user/userservice"
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

var userClient userservice.Client

func initUserRpc() {
	r, err := etcd.NewEtcdResolver([]string{constants.EtcdAddress})
	if err != nil {
		panic(err)
	}

	c, err := userservice.NewClient(
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
	userClient = c
}

func RegisterUser(ctx context.Context, req *user.UserRegisterRequest) (string, string, error) {
	resp, err := userClient.Register(ctx, req)
	if err != nil {
		return "", "", err
	}
	if resp.BaseResp.StatusCode != 0 {
		return "", "", errno.NewErr(resp.BaseResp.StatusCode, resp.BaseResp.StatusMsg)
	}
	return resp.UserId, resp.Token, nil
}

func LoginUser(ctx context.Context, req *user.UserLoginRequest) (string, string, error) {
	resp, err := userClient.Login(ctx, req)
	if err != nil {
		return "", "", err
	}
	if resp.BaseResp.StatusCode != 0 {
		return "", "", errno.NewErr(resp.BaseResp.StatusCode, resp.BaseResp.StatusMsg)
	}

	return resp.UserId, resp.Token, nil
}

// UserInfo get user info
func UserInfo(ctx context.Context, req *user.UserInfoRequest) (*user.User, error) {
	resp, err := userClient.UserInfo(ctx, req)
	if err != nil {
		return nil, err
	}
	if resp.BaseResp.StatusCode != 0 {
		return nil, errno.NewErr(resp.BaseResp.StatusCode, resp.BaseResp.StatusMsg)
	}
	return resp.User[0], nil
}

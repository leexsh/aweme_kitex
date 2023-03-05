package relationRPC

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

func Init() {
	initUserRpc()
}

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

// 关注或者取消关注-修改user的数据库信息
func ChangeFollowCount(ctx context.Context, req *user.ChangeFollowStatusRequest) error {
	return userClient.ChangeFollowStatus(ctx, req)
}

// 获取用户信息
func GetUserInfo(ctx context.Context, req *user.SingleUserInfoRequest) ([]*user.User, error) {
	resp, err := userClient.GetUserInfoByUserId(ctx, req)
	if err != nil {
		return nil, err
	}
	if resp.BaseResp.StatusCode != 0 {
		return nil, errno.NewErr(resp.BaseResp.StatusCode, resp.BaseResp.StatusMsg)
	}
	return resp.Users, nil
}
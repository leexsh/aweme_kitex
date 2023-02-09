package rpc

import (
	comment2 "aweme_kitex/cmd/comment/kitex_gen/comment"
	comment "aweme_kitex/cmd/comment/kitex_gen/comment/commentservice"
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

var commentClient comment.Client

func initCommentRpc() {
	r, err := etcd.NewEtcdResolver([]string{constants.EtcdAddress})
	if err != nil {
		logger.Error("feed service_user register error")
		return
	}
	client, err := comment.NewClient(
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
	commentClient = client
}

// CreateComment add comment
func CreateComment(ctx context.Context, req *comment2.CommentActionRequest) (*comment2.Comment, error) {
	resp, err := commentClient.CreateComment(ctx, req)
	if err != nil {
		return nil, err
	}
	if resp.BaseResp.StatusCode != 0 {
		return nil, errno.NewErr(resp.BaseResp.StatusCode, resp.BaseResp.StatusMsg)
	}
	return resp.CommentList[0], nil
}

// DeleteComment delete comment
func DeleteComment(ctx context.Context, req *comment2.CommentActionRequest) (*comment2.Comment, error) {
	resp, err := commentClient.DelComment(ctx, req)
	if err != nil {
		return nil, err
	}
	if resp.BaseResp.StatusCode != 0 {
		return nil, errno.NewErr(resp.BaseResp.StatusCode, resp.BaseResp.StatusMsg)
	}
	return resp.CommentList[0], nil
}

// CommentList get comment list info
func CommentList(ctx context.Context, req *comment2.CommentListRequest) ([]*comment2.Comment, error) {
	resp, err := commentClient.CommentList(ctx, req)
	if err != nil {
		return nil, err
	}
	if resp.BaseResp.StatusCode != 0 {
		return nil, errno.NewErr(resp.BaseResp.StatusCode, resp.BaseResp.StatusMsg)
	}
	return resp.CommentList, nil
}

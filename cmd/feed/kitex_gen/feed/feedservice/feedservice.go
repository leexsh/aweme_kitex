// Code generated by Kitex v0.4.4. DO NOT EDIT.

package feedservice

import (
	feed "aweme_kitex/cmd/feed/kitex_gen/feed"
	"context"
	client "github.com/cloudwego/kitex/client"
	kitex "github.com/cloudwego/kitex/pkg/serviceinfo"
)

func serviceInfo() *kitex.ServiceInfo {
	return feedServiceServiceInfo
}

var feedServiceServiceInfo = NewServiceInfo()

func NewServiceInfo() *kitex.ServiceInfo {
	serviceName := "FeedService"
	handlerType := (*feed.FeedService)(nil)
	methods := map[string]kitex.MethodInfo{
		"Feed":              kitex.NewMethodInfo(feedHandler, newFeedServiceFeedArgs, newFeedServiceFeedResult, false),
		"ChangeCommentCnt":  kitex.NewMethodInfo(changeCommentCntHandler, newFeedServiceChangeCommentCntArgs, newFeedServiceChangeCommentCntResult, false),
		"CheckVideoInvalid": kitex.NewMethodInfo(checkVideoInvalidHandler, newFeedServiceCheckVideoInvalidArgs, newFeedServiceCheckVideoInvalidResult, false),
	}
	extra := map[string]interface{}{
		"PackageName": "feed",
	}
	svcInfo := &kitex.ServiceInfo{
		ServiceName:     serviceName,
		HandlerType:     handlerType,
		Methods:         methods,
		PayloadCodec:    kitex.Thrift,
		KiteXGenVersion: "v0.4.4",
		Extra:           extra,
	}
	return svcInfo
}

func feedHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*feed.FeedServiceFeedArgs)
	realResult := result.(*feed.FeedServiceFeedResult)
	success, err := handler.(feed.FeedService).Feed(ctx, realArg.Req)
	if err != nil {
		return err
	}
	realResult.Success = success
	return nil
}
func newFeedServiceFeedArgs() interface{} {
	return feed.NewFeedServiceFeedArgs()
}

func newFeedServiceFeedResult() interface{} {
	return feed.NewFeedServiceFeedResult()
}

func changeCommentCntHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*feed.FeedServiceChangeCommentCntArgs)
	realResult := result.(*feed.FeedServiceChangeCommentCntResult)
	success, err := handler.(feed.FeedService).ChangeCommentCnt(ctx, realArg.Req)
	if err != nil {
		return err
	}
	realResult.Success = success
	return nil
}
func newFeedServiceChangeCommentCntArgs() interface{} {
	return feed.NewFeedServiceChangeCommentCntArgs()
}

func newFeedServiceChangeCommentCntResult() interface{} {
	return feed.NewFeedServiceChangeCommentCntResult()
}

func checkVideoInvalidHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*feed.FeedServiceCheckVideoInvalidArgs)
	realResult := result.(*feed.FeedServiceCheckVideoInvalidResult)
	success, err := handler.(feed.FeedService).CheckVideoInvalid(ctx, realArg.Req)
	if err != nil {
		return err
	}
	realResult.Success = success
	return nil
}
func newFeedServiceCheckVideoInvalidArgs() interface{} {
	return feed.NewFeedServiceCheckVideoInvalidArgs()
}

func newFeedServiceCheckVideoInvalidResult() interface{} {
	return feed.NewFeedServiceCheckVideoInvalidResult()
}

type kClient struct {
	c client.Client
}

func newServiceClient(c client.Client) *kClient {
	return &kClient{
		c: c,
	}
}

func (p *kClient) Feed(ctx context.Context, req *feed.FeedRequest) (r *feed.FeedResponse, err error) {
	var _args feed.FeedServiceFeedArgs
	_args.Req = req
	var _result feed.FeedServiceFeedResult
	if err = p.c.Call(ctx, "Feed", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}

func (p *kClient) ChangeCommentCnt(ctx context.Context, req *feed.ChangeCommentCountRequest) (r *feed.ChangeCommentCountResponse, err error) {
	var _args feed.FeedServiceChangeCommentCntArgs
	_args.Req = req
	var _result feed.FeedServiceChangeCommentCntResult
	if err = p.c.Call(ctx, "ChangeCommentCnt", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}

func (p *kClient) CheckVideoInvalid(ctx context.Context, req *feed.CheckVideoInvalidRequest) (r *feed.CheckVideoInvalidResponse, err error) {
	var _args feed.FeedServiceCheckVideoInvalidArgs
	_args.Req = req
	var _result feed.FeedServiceCheckVideoInvalidResult
	if err = p.c.Call(ctx, "CheckVideoInvalid", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}

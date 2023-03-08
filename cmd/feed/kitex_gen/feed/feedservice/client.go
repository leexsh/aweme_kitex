// Code generated by Kitex v0.4.4. DO NOT EDIT.

package feedservice

import (
	feed "aweme_kitex/cmd/feed/kitex_gen/feed"
	"context"
	client "github.com/cloudwego/kitex/client"
	callopt "github.com/cloudwego/kitex/client/callopt"
)

// Client is designed to provide IDL-compatible methods with call-option parameter for kitex framework.
type Client interface {
	Feed(ctx context.Context, req *feed.FeedRequest, callOptions ...callopt.Option) (r *feed.FeedResponse, err error)
	ChangeCommentCnt(ctx context.Context, req *feed.ChangeCommentCountRequest, callOptions ...callopt.Option) (r *feed.ChangeCommentCountResponse, err error)
	CheckVideoInvalid(ctx context.Context, req *feed.CheckVideoInvalidRequest, callOptions ...callopt.Option) (r *feed.CheckVideoInvalidResponse, err error)
	GetVideosById(ctx context.Context, req *feed.CheckVideoInvalidRequest, callOptions ...callopt.Option) (r *feed.GetVideosResponse, err error)
}

// NewClient creates a client for the service defined in IDL.
func NewClient(destService string, opts ...client.Option) (Client, error) {
	var options []client.Option
	options = append(options, client.WithDestService(destService))

	options = append(options, opts...)

	kc, err := client.NewClient(serviceInfo(), options...)
	if err != nil {
		return nil, err
	}
	return &kFeedServiceClient{
		kClient: newServiceClient(kc),
	}, nil
}

// MustNewClient creates a client for the service defined in IDL. It panics if any error occurs.
func MustNewClient(destService string, opts ...client.Option) Client {
	kc, err := NewClient(destService, opts...)
	if err != nil {
		panic(err)
	}
	return kc
}

type kFeedServiceClient struct {
	*kClient
}

func (p *kFeedServiceClient) Feed(ctx context.Context, req *feed.FeedRequest, callOptions ...callopt.Option) (r *feed.FeedResponse, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.Feed(ctx, req)
}

func (p *kFeedServiceClient) ChangeCommentCnt(ctx context.Context, req *feed.ChangeCommentCountRequest, callOptions ...callopt.Option) (r *feed.ChangeCommentCountResponse, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.ChangeCommentCnt(ctx, req)
}

func (p *kFeedServiceClient) CheckVideoInvalid(ctx context.Context, req *feed.CheckVideoInvalidRequest, callOptions ...callopt.Option) (r *feed.CheckVideoInvalidResponse, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.CheckVideoInvalid(ctx, req)
}

func (p *kFeedServiceClient) GetVideosById(ctx context.Context, req *feed.CheckVideoInvalidRequest, callOptions ...callopt.Option) (r *feed.GetVideosResponse, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.GetVideosById(ctx, req)
}

package service

import (
	"aweme_kitex/cmd/publish/kitex_gen/publish"
	"context"
)

type PublishService struct {
	ctx context.Context
}

// NewPublishService new PublishService
func NewPublishService(ctx context.Context) *PublishService {
	return &PublishService{ctx: ctx}
}

// Publish upload video info
func (s *PublishService) Publish(req *publish.PublishActionRequest) error {
	// video := req.Data
	// title := req.Title
	// currentId := req.UserId

	// 1.将视频保存到本地文件夹
	// 2.上传oss
	// 3.获取播放链接
	return nil
}

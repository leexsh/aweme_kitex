package service_feed

import (
	db2 "aweme_kitex/cmd/feed/service_feed/db"
	"context"
)

type CheckVideoService struct {
	ctx context.Context
}

// NewCheckVideoService new CheckVideoService
func NewCheckVideoService(ctx context.Context) *CheckVideoService {
	return &CheckVideoService{
		ctx: ctx,
	}
}

func (c *CheckVideoService) CheckVideoInvalid(vids []string) error {
	_, err := db2.NewVideoDaoInstance().QueryVideosByIs(c.ctx, vids)
	return err
}

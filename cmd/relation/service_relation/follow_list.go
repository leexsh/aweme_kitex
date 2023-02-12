package service_relation

import (
	"aweme_kitex/cmd/relation/kitex_gen/relation"
	"aweme_kitex/cmd/relation/kitex_gen/user"
	"aweme_kitex/pkg/jwt"
	"context"
)

type FollowListService struct {
	ctx context.Context
}

// NewFollowListService new FollowListService
func NewFollowListService(ctx context.Context) *FollowListService {
	return &FollowListService{ctx: ctx}
}

// FollowList get user follow list info
func (s *FollowListService) FollowList(req *relation.FollowListRequest) ([]*user.User, error) {
	uc, err := jwt.AnalyzeToken(req.Token)
	if err != nil {
		return nil, err
	}
	return getFollowList(s.ctx, uc.Id)
}

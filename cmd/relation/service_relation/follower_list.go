package service_relation

import (
	"aweme_kitex/cmd/relation/kitex_gen/relation"
	"aweme_kitex/cmd/relation/kitex_gen/user"
	"aweme_kitex/pkg/jwt"
	"context"
)

type FollowerListService struct {
	ctx context.Context
}

// NewFollowerListService new FollowerListService
func NewFollowerListService(ctx context.Context) *FollowerListService {
	return &FollowerListService{ctx: ctx}
}

// FollowerList get user follower list info
func (s *FollowerListService) FollowerList(req *relation.FollowerListRequest) ([]*user.User, error) {
	uc, err := jwt.AnalyzeToken(req.Token)
	if err != nil {
		return nil, err
	}
	return getFollowerList(s.ctx, uc.Id)
}

package service_relation

import (
	"aweme_kitex/cmd/relation/kitex_gen/relation"
	"aweme_kitex/cmd/relation/kitex_gen/user"
	"aweme_kitex/models"
	"aweme_kitex/pkg/jwt"
	"aweme_kitex/service"
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
	userList, err := service.GetFollowList(uc.Id)
	return s.packUserInfo(userList), nil
}

func (s *FollowListService) packUserInfo(modelUsers []*models.User) []*user.User {
	userList := make([]*user.User, 0)
	for _, mUser := range modelUsers {
		uer := &user.User{
			mUser.UserId,
			mUser.Name,
			mUser.FollowCount,
			mUser.FollowerCount,
			mUser.IsFollow,
		}
		userList = append(userList, uer)
	}
	return userList
}

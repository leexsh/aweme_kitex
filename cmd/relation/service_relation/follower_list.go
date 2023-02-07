package service_relation

import (
	"aweme_kitex/cmd/relation/kitex_gen/relation"
	"aweme_kitex/cmd/relation/kitex_gen/user"
	"aweme_kitex/models"
	"aweme_kitex/pkg/jwt"
	"aweme_kitex/service"
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
	userList, err := service.GetFollowerList(uc.Id)
	return s.packUserInfo(userList), nil
}

func (s *FollowerListService) packUserInfo(modelUsers []*models.User) []*user.User {
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

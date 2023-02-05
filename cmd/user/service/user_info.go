package service_user

import (
	user2 "aweme_kitex/cmd/publish/kitex_gen/user"
	"aweme_kitex/cmd/user/kitex_gen/user"
	"aweme_kitex/models/dal"
	"aweme_kitex/pkg/jwt"
	"context"
	"errors"
)

type UserInfoService struct {
	ctx context.Context
}

func NewUserInfoService(ctx context.Context) *UserInfoService {
	return &UserInfoService{
		ctx: ctx,
	}
}

func (s *UserInfoService) UserInfo(req *user.UserInfoRequest) (user *user2.User, err error) {
	u, err := dal.NewUserDaoInstance().QueryUserByPassword(s.ctx, req.UserName, req.Password)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("not found user, please register")
	}
	user = new(user2.User)
	user.UserId = u.UserId
	user.Name = u.Name
	user.FollowerCount = u.FollowerCount
	user.FollowCount = u.FollowCount
	user.IsFollow = false
	return user, nil
}

func (s *UserInfoService) CheckToken(token string) (*jwt.UserClaim, error) {
	uc, err := jwt.AnalyzeToken(token)
	if err != nil {
		return nil, err
	}
	return uc, nil
}

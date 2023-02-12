package service_user

import (
	user2 "aweme_kitex/cmd/publish/kitex_gen/user"
	"aweme_kitex/cmd/user/kitex_gen/user"
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
	uc, err := jwt.AnalyzeToken(req.Token)
	res, err := QueryUserInfo(uc, req.UserId)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("not found user, please register")
	}
	user = new(user2.User)
	user.UserId = res.UserId
	user.Name = res.Name
	user.FollowerCount = res.FollowerCount
	user.FollowCount = res.FollowCount
	user.IsFollow = res.IsFollow
	return user, nil
}

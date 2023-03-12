package service_user

import (
	"aweme_kitex/cmd/user/kitex_gen/user"
	userDB "aweme_kitex/cmd/user/service_user/db"
	"aweme_kitex/pkg/utils"
	"context"
)

type LoginUserService struct {
	ctx context.Context
}

// NewRegisterUserService new RegisterUserService
func NewLoginUserService(ctx context.Context) *LoginUserService {
	return &LoginUserService{
		ctx: ctx,
	}
}

// RegisterUser register user info
func (s *LoginUserService) LoginUser(req *user.UserLoginRequest) (string, string, error) {
	userId, token, err := s.do(req.UserName, req.Password)
	if err != nil {
		return "", "", err
	}
	return userId, token, err
}

func (s *LoginUserService) do(name, password string) (uid string, token string, err error) {
	usr, err := userDB.NewUserDaoInstance().QueryUserByPassword(s.ctx, name, utils.Md5(password))
	if err != nil {
		return "", "", err
	}
	return usr.UserId, usr.Token, err
}

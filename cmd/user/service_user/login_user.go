package service_user

import (
	"aweme_kitex/cmd/user/kitex_gen/user"
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
	userId, token, err := LoginUser(req.UserName, req.Password)
	if err != nil {
		return "", "", err
	}
	return userId, token, err
}

package service_user

import (
	"aweme_kitex/cmd/user/kitex_gen/user"
	userDB "aweme_kitex/cmd/user/service_user/db"
	"aweme_kitex/pkg/jwt"
	"aweme_kitex/pkg/logger"
	"aweme_kitex/pkg/types"
	"aweme_kitex/pkg/utils"
	"context"
)

type RegisterUserService struct {
	ctx context.Context
}

// NewRegisterUserService new RegisterUserService
func NewRegisterUserService(ctx context.Context) *RegisterUserService {
	return &RegisterUserService{
		ctx: ctx,
	}
}

// RegisterUser register user info
func (s *RegisterUserService) RegisterUser(req *user.UserRegisterRequest) (string, string, error) {
	userId, token, err := s.do(req.UserName, req.Password)
	if err != nil {
		return "", "", err
	}
	return userId, token, err
}

func (s *RegisterUserService) do(name, password string) (string, string, error) {
	// insert to data
	userId := utils.GenerateUUID()
	token, _ := jwt.GenerateToken(userId, name)
	newUser := &types.UserRawData{
		UserId:        userId,
		Name:          name,
		Password:      utils.Md5(password),
		Token:         token,
		FollowCount:   0,
		FollowerCount: 0,
	}
	logger.Info("Register success userName is %s", name)
	err := userDB.NewUserDaoInstance().UploadUserData(context.Background(), newUser)
	return userId, token, err
}

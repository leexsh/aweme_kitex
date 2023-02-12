package service_user

import (
	"aweme_kitex/cmd/user/kitex_gen/user"
	"aweme_kitex/models"
	"aweme_kitex/models/dal"
	"aweme_kitex/pkg/errno"
	"aweme_kitex/pkg/jwt"
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
	user, err := dal.NewUserDaoInstance().QueryUserByUserId(s.ctx, req.UserName)
	if err != nil {
		return "", "", err
	}
	if len(user.UserId) != 0 {
		return "", "", errno.UserAlreadyExistErr
	}

	userId := utils.GenerateUUID()
	token, _ := jwt.GenerateToken(userId, req.UserName)
	newUser := &models.UserRawData{
		UserId:        userId,
		Name:          req.UserName,
		Password:      utils.Md5(req.Password),
		Token:         token,
		FollowCount:   0,
		FollowerCount: 0,
	}
	err = dal.NewUserDaoInstance().UploadUserData(context.Background(), newUser)

	return userId, token, nil
}

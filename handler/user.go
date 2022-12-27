package handler

import (
	"aweme_kitex/models"
	"aweme_kitex/service"
	"errors"
	"fmt"
)

// handle user logic

// register
func UserRegisterHandle(userName, password string) *models.UserLogRstResponse {
	if len := len(userName); len <= 0 || len > 32 {
		return userLogRstErrRes(errors.New("register user name error"))
	}
	if len := len(password); len <= 0 || len > 32 {
		return userLogRstErrRes(errors.New("register password out of range"))
	}
	// id name token
	id, token, err := service.RegisterUser(userName, password)
	if err != nil {
		return userLogRstErrRes(err)
	}

	return &models.UserLogRstResponse{
		models.Response{
			0,
			fmt.Sprintf("%s register success!", userName),
		},
		id,
		userName,
		token,
	}
}

// login
func UserLoginHandle(userName, password string) *models.UserLogRstResponse {
	if len := len(userName); len <= 0 || len > 32 {
		return userLogRstErrRes(errors.New("login user name error"))
	}
	if len := len(password); len <= 0 || len > 32 {
		return userLogRstErrRes(errors.New("login password out of range"))
	}
	uid, token, err := service.LoginUser(userName, password)
	if err != nil {
		return userLogRstErrRes(err)
	}
	return &models.UserLogRstResponse{
		models.Response{
			0,
			fmt.Sprintf("%s register success!", userName),
		},
		uid,
		userName,
		token,
	}
}

// userinfo
func UserInfoHandle(user *models.UserClaim, remoteUid string) *models.UserResponse {
	res, err := service.QueryUserInfo(user, remoteUid)
	if err != nil {
		return &models.UserResponse{
			Response: models.Response{
				-1,
				err.Error(),
			},
		}
	}
	return &models.UserResponse{
		models.Response{
			0,
			"get user info success",
		},
		*res,
	}
}

func userLogRstErrRes(err error) *models.UserLogRstResponse {
	return &models.UserLogRstResponse{
		Response: models.Response{
			-1,
			err.Error(),
		},
	}
}

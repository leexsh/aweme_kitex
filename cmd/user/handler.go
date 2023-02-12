package main

import (
	user "aweme_kitex/cmd/user/kitex_gen/user"
	"aweme_kitex/cmd/user/service_user"
	"context"
	"time"
)

// UserServiceImpl implements the last service interface defined in the IDL.
type UserServiceImpl struct{}

// Register implements the UserServiceImpl interface.
func (s *UserServiceImpl) Register(ctx context.Context, req *user.UserRegisterRequest) (resp *user.UserRegisterResponse, err error) {
	// TODO: Your code here...
	resp = new(user.UserRegisterResponse)
	if len(req.UserName) == 0 || len(req.Password) == 0 {
		resp.BaseResp.ServiceTime = time.Now().Unix()
		resp.BaseResp.StatusCode = -1
		resp.BaseResp.StatusMsg = "register error"
		return resp, nil
	}
	userId, token, err := service_user.NewRegisterUserService(ctx).RegisterUser(req)
	if err != nil {
		resp.BaseResp.ServiceTime = time.Now().Unix()
		resp.BaseResp.StatusCode = -1
		resp.BaseResp.StatusMsg = "register error"
		return resp, nil
	}
	resp.BaseResp.StatusCode = 0
	resp.BaseResp.StatusMsg = "register success"
	resp.UserId = userId
	resp.Token = token
	return resp, nil
}

// Login implements the UserServiceImpl interface.
func (s *UserServiceImpl) Login(ctx context.Context, req *user.UserLoginRequest) (resp *user.UserLoginResponse, err error) {
	// TODO: Your code here...
	resp = new(user.UserLoginResponse)
	if len(req.UserName) == 0 || len(req.Password) == 0 {
		resp.BaseResp.ServiceTime = time.Now().Unix()
		resp.BaseResp.StatusCode = -1
		resp.BaseResp.StatusMsg = "login error"
		return resp, nil
	}
	userId, token, err := service_user.NewLoginUserService(ctx).LoginUser(req)
	if err != nil {
		resp.BaseResp.ServiceTime = time.Now().Unix()
		resp.BaseResp.StatusCode = -1
		resp.BaseResp.StatusMsg = "login error"
		return resp, nil
	}
	resp.BaseResp.StatusCode = 0
	resp.BaseResp.StatusMsg = "login success"
	resp.UserId = userId
	resp.Token = token
	return
}

// UserInfo implements the UserServiceImpl interface.
func (s *UserServiceImpl) UserInfo(ctx context.Context, req *user.UserInfoRequest) (resp *user.UserInfoResponse, err error) {
	resp = new(user.UserInfoResponse)
	user, err := service_user.NewUserInfoService(ctx).UserInfo(req)
	if err != nil {
		resp.BaseResp.ServiceTime = time.Now().Unix()
		resp.BaseResp.StatusCode = -1
		resp.BaseResp.StatusMsg = "get user info error"
		return resp, nil
	}
	resp.BaseResp.StatusCode = 0
	resp.BaseResp.StatusMsg = "getuserInfo success"
	resp.User[0].UserId = user.UserId
	resp.User[0].Name = user.Name
	resp.User[0].IsFollow = user.IsFollow
	resp.User[0].FollowCount = user.FollowerCount
	resp.User[0].FollowerCount = user.FollowCount
	return resp, nil
	return
}

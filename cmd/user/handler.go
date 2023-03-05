package main

import (
	user "aweme_kitex/cmd/user/kitex_gen/user"
	"aweme_kitex/cmd/user/service_user"
	userPack "aweme_kitex/cmd/user/service_user/pack"
	"context"
	"time"
)

// UserServiceImpl implements the last service interface defined in the IDL.
type UserServiceImpl struct{}

// Register implements the UserServiceImpl interface.
func (s *UserServiceImpl) Register(ctx context.Context, req *user.UserRegisterRequest) (resp *user.UserRegisterResponse, err error) {
	if len(req.UserName) == 0 || len(req.Password) == 0 {
		return userPack.RegisterResponse(-1, "register error"), nil
	}
	userId, token, err := service_user.NewRegisterUserService(ctx).RegisterUser(req)
	if err != nil {
		return userPack.RegisterResponse(-1, "register error"), nil
	}
	resp = userPack.RegisterResponse(0, "register success")
	resp.UserId = userId
	resp.Token = token
	return resp, nil
}

// Login implements the UserServiceImpl interface.
func (s *UserServiceImpl) Login(ctx context.Context, req *user.UserLoginRequest) (resp *user.UserLoginResponse, err error) {
	resp = new(user.UserLoginResponse)
	if len(req.UserName) == 0 || len(req.Password) == 0 {
		return userPack.LoginResponse(-1, "Login error"), nil
	}
	userId, token, err := service_user.NewLoginUserService(ctx).LoginUser(req)
	if err != nil {
		return userPack.LoginResponse(-1, "Login error"), nil
	}
	resp = userPack.LoginResponse(0, "login success")
	resp.UserId = userId
	resp.Token = token
	return
}

// UserInfo implements the UserServiceImpl interface.
func (s *UserServiceImpl) UserInfo(ctx context.Context, req *user.UserInfoRequest) (resp *user.UserInfoResponse, err error) {
	resp = new(user.UserInfoResponse)
	usr, err := service_user.NewUserInfoService(ctx).UserInfo(req)
	if err != nil {
		return userPack.UserInfoResponse(-1, "get user info error", nil), nil
	}
	return userPack.UserInfoResponse(0, "getUserInfo success", []*user.User{usr}), nil
}

// GetUserInfoByUserId implements the UserServiceImpl interface.
// for other rpc service to get single user info
func (s *UserServiceImpl) GetUserInfoByUserId(ctx context.Context, req *user.SingleUserInfoRequest) (resp *user.SingleUserInfoResponse, err error) {
	resp = new(user.SingleUserInfoResponse)
	for _, uid := range req.UserIds {
		us, err := service_user.NewUserInfoService(ctx).SingleUserInfo(uid)
		if err != nil {
			resp.BaseResp.ServiceTime = time.Now().Unix()
			resp.BaseResp.StatusCode = -1
			resp.BaseResp.StatusMsg = "get signle user info error"
			return resp, err
		}
		resp.Users = append(resp.Users, us)
	}
	resp.BaseResp.StatusCode = 0
	resp.BaseResp.StatusMsg = "getuserInfo success"
	return resp, nil
}

// ChangeFollowStatus implements the UserServiceImpl interface.
func (s *UserServiceImpl) ChangeFollowStatus(ctx context.Context, req *user.ChangeFollowStatusRequest) (err error) {
	err = service_user.NewChangeFollowService(ctx).ChangeStatus(req.UserId, req.ToUserId, req.IsFollow)
	return err
}

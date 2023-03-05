// Code generated by Kitex v0.4.4. DO NOT EDIT.

package userservice

import (
	user "aweme_kitex/cmd/relation/kitex_gen/user"
	"context"
	client "github.com/cloudwego/kitex/client"
	kitex "github.com/cloudwego/kitex/pkg/serviceinfo"
)

func serviceInfo() *kitex.ServiceInfo {
	return userServiceServiceInfo
}

var userServiceServiceInfo = NewServiceInfo()

func NewServiceInfo() *kitex.ServiceInfo {
	serviceName := "UserService"
	handlerType := (*user.UserService)(nil)
	methods := map[string]kitex.MethodInfo{
		"Register":            kitex.NewMethodInfo(registerHandler, newUserServiceRegisterArgs, newUserServiceRegisterResult, false),
		"Login":               kitex.NewMethodInfo(loginHandler, newUserServiceLoginArgs, newUserServiceLoginResult, false),
		"UserInfo":            kitex.NewMethodInfo(userInfoHandler, newUserServiceUserInfoArgs, newUserServiceUserInfoResult, false),
		"GetUserInfoByUserId": kitex.NewMethodInfo(getUserInfoByUserIdHandler, newUserServiceGetUserInfoByUserIdArgs, newUserServiceGetUserInfoByUserIdResult, false),
		"ChangeFollowStatus":  kitex.NewMethodInfo(changeFollowStatusHandler, newUserServiceChangeFollowStatusArgs, newUserServiceChangeFollowStatusResult, false),
	}
	extra := map[string]interface{}{
		"PackageName": "user",
	}
	svcInfo := &kitex.ServiceInfo{
		ServiceName:     serviceName,
		HandlerType:     handlerType,
		Methods:         methods,
		PayloadCodec:    kitex.Thrift,
		KiteXGenVersion: "v0.4.4",
		Extra:           extra,
	}
	return svcInfo
}

func registerHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*user.UserServiceRegisterArgs)
	realResult := result.(*user.UserServiceRegisterResult)
	success, err := handler.(user.UserService).Register(ctx, realArg.Req)
	if err != nil {
		return err
	}
	realResult.Success = success
	return nil
}
func newUserServiceRegisterArgs() interface{} {
	return user.NewUserServiceRegisterArgs()
}

func newUserServiceRegisterResult() interface{} {
	return user.NewUserServiceRegisterResult()
}

func loginHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*user.UserServiceLoginArgs)
	realResult := result.(*user.UserServiceLoginResult)
	success, err := handler.(user.UserService).Login(ctx, realArg.Req)
	if err != nil {
		return err
	}
	realResult.Success = success
	return nil
}
func newUserServiceLoginArgs() interface{} {
	return user.NewUserServiceLoginArgs()
}

func newUserServiceLoginResult() interface{} {
	return user.NewUserServiceLoginResult()
}

func userInfoHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*user.UserServiceUserInfoArgs)
	realResult := result.(*user.UserServiceUserInfoResult)
	success, err := handler.(user.UserService).UserInfo(ctx, realArg.Req)
	if err != nil {
		return err
	}
	realResult.Success = success
	return nil
}
func newUserServiceUserInfoArgs() interface{} {
	return user.NewUserServiceUserInfoArgs()
}

func newUserServiceUserInfoResult() interface{} {
	return user.NewUserServiceUserInfoResult()
}

func getUserInfoByUserIdHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*user.UserServiceGetUserInfoByUserIdArgs)
	realResult := result.(*user.UserServiceGetUserInfoByUserIdResult)
	success, err := handler.(user.UserService).GetUserInfoByUserId(ctx, realArg.Req)
	if err != nil {
		return err
	}
	realResult.Success = success
	return nil
}
func newUserServiceGetUserInfoByUserIdArgs() interface{} {
	return user.NewUserServiceGetUserInfoByUserIdArgs()
}

func newUserServiceGetUserInfoByUserIdResult() interface{} {
	return user.NewUserServiceGetUserInfoByUserIdResult()
}

func changeFollowStatusHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*user.UserServiceChangeFollowStatusArgs)

	err := handler.(user.UserService).ChangeFollowStatus(ctx, realArg.Req)
	if err != nil {
		return err
	}

	return nil
}
func newUserServiceChangeFollowStatusArgs() interface{} {
	return user.NewUserServiceChangeFollowStatusArgs()
}

func newUserServiceChangeFollowStatusResult() interface{} {
	return user.NewUserServiceChangeFollowStatusResult()
}

type kClient struct {
	c client.Client
}

func newServiceClient(c client.Client) *kClient {
	return &kClient{
		c: c,
	}
}

func (p *kClient) Register(ctx context.Context, req *user.UserRegisterRequest) (r *user.UserRegisterResponse, err error) {
	var _args user.UserServiceRegisterArgs
	_args.Req = req
	var _result user.UserServiceRegisterResult
	if err = p.c.Call(ctx, "Register", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}

func (p *kClient) Login(ctx context.Context, req *user.UserLoginRequest) (r *user.UserLoginResponse, err error) {
	var _args user.UserServiceLoginArgs
	_args.Req = req
	var _result user.UserServiceLoginResult
	if err = p.c.Call(ctx, "Login", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}

func (p *kClient) UserInfo(ctx context.Context, req *user.UserInfoRequest) (r *user.UserInfoResponse, err error) {
	var _args user.UserServiceUserInfoArgs
	_args.Req = req
	var _result user.UserServiceUserInfoResult
	if err = p.c.Call(ctx, "UserInfo", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}

func (p *kClient) GetUserInfoByUserId(ctx context.Context, req *user.SingleUserInfoRequest) (r *user.SingleUserInfoResponse, err error) {
	var _args user.UserServiceGetUserInfoByUserIdArgs
	_args.Req = req
	var _result user.UserServiceGetUserInfoByUserIdResult
	if err = p.c.Call(ctx, "GetUserInfoByUserId", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}

func (p *kClient) ChangeFollowStatus(ctx context.Context, req *user.ChangeFollowStatusRequest) (err error) {
	var _args user.UserServiceChangeFollowStatusArgs
	_args.Req = req
	var _result user.UserServiceChangeFollowStatusResult
	if err = p.c.Call(ctx, "ChangeFollowStatus", &_args, &_result); err != nil {
		return
	}
	return nil
}

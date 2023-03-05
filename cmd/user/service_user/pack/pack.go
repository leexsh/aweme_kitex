package userPack

import (
	"aweme_kitex/cmd/user/kitex_gen/user"
	"time"
)

func RegisterResponse(code int64, msg string) (resp *user.UserRegisterResponse) {
	resp = new(user.UserRegisterResponse)
	resp.BaseResp.ServiceTime = time.Now().Unix()
	resp.BaseResp.StatusCode = code
	resp.BaseResp.StatusMsg = msg
	return
}

func LoginResponse(code int64, msg string) (resp *user.UserLoginResponse) {
	resp = new(user.UserLoginResponse)
	resp.BaseResp.ServiceTime = time.Now().Unix()
	resp.BaseResp.StatusCode = code
	resp.BaseResp.StatusMsg = msg
	return
}

func UserInfoResponse(code int64, msg string, users []*user.User) (resp *user.UserInfoResponse) {
	resp = new(user.UserInfoResponse)
	resp.BaseResp.ServiceTime = time.Now().Unix()
	resp.BaseResp.StatusCode = code
	resp.BaseResp.StatusMsg = msg
	if users != nil {
		resp.User = users
	}
	return
}

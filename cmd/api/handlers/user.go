package handlers

import (
	"aweme_kitex/cmd/api/rpc"
	"aweme_kitex/cmd/user/kitex_gen/user"
	"aweme_kitex/pkg/errno"
	"aweme_kitex/pkg/jwt"
	"context"

	"github.com/gin-gonic/gin"
)

func Register(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")
	if len := len(username); len <= 0 || len > 32 {
		SendResponse(c, errno.ParamErr, nil)
		return
	}
	if len := len(password); len <= 0 || len > 32 {
		SendResponse(c, errno.ParamErr, nil)
		return
	}
	userId, token, err := rpc.RegisterUser(context.Background(), &user.UserRegisterRequest{
		UserName: username,
		Password: password,
	})
	if err != nil {
		SendResponse(c, errno.ConvertErr(err), nil)
	}
	SendResponse(c, errno.Success, map[string]interface{}{"token": token, "userId": userId})
}

func Login(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")
	if len := len(username); len <= 0 || len > 32 {
		SendResponse(c, errno.ParamErr, nil)
		return
	}
	if len := len(password); len <= 0 || len > 32 {
		SendResponse(c, errno.ParamErr, nil)
		return
	}
	userId, token, err := rpc.LoginUser(context.Background(), &user.UserLoginRequest{username, password})
	if err != nil {
		SendResponse(c, errno.ConvertErr(err), nil)
	}
	SendResponse(c, errno.Success, map[string]interface{}{"token": token, "userId": userId})

}

// UserInfo get user info
func UserInfo(c *gin.Context) {
	token := c.DefaultQuery("token", "")
	_, err := jwt.AnalyzeToken(token)
	if err != nil {
		SendResponse(c, errno.TokenInvalidErr, nil)
		return
	}
	userid := c.Query("userid")
	if len := len(userid); len <= 0 || len > 32 {
		SendResponse(c, errno.ParamErr, nil)
		return
	}
	req := &user.UserInfoRequest{
		Token:  token,
		UserId: userid,
	}
	user, err := rpc.UserInfo(context.Background(), req)
	if err != nil {
		SendResponse(c, errno.ConvertErr(err), nil)
	}
	SendResponse(c, errno.Success, map[string]interface{}{"user": user})

}

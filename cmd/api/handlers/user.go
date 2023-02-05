package handlers

import (
	"aweme_kitex/cmd/api/rpc"
	"aweme_kitex/cmd/user/kitex_gen/user"
	"aweme_kitex/pkg/errno"
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
)

func Register(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

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
	userId, token, err := rpc.LoginUser(context.Background(), &user.UserLoginRequest{username, password})
	if err != nil {
		SendResponse(c, errno.ConvertErr(err), nil)
	}
	SendResponse(c, errno.Success, map[string]interface{}{"token": token, "userId": userId})

}

// UserInfo get user info
func UserInfo(c *gin.Context) {
	token := c.DefaultQuery("token", "")
	fmt.Println(token)

	// user, err := rpc.UserInfo(context.Background(), &user.UserInfoRequest{
	// 	UserId: userId,
	// 	Token:  token,
	// })
	// if err != nil {
	// 	SendResponse(c, errno.ConvertErr(err), nil)
	// }
	// SendResponse(c, errno.Success, map[string]interface{}{constants.User: user})

}

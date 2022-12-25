package controller

import (
	"aweme_kitex/handler"
	"net/http"

	"github.com/gin-gonic/gin"
)

/*
用户
*/

func Register(c *gin.Context) {
	userName := c.Query("username")
	password := c.Query("password")
	c.JSON(http.StatusOK, handler.UserRegisterHandle(userName, password))
}

func Login(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")
	c.JSON(200, handler.UserLoginHandle(username, password))
}

func UserInfo(c *gin.Context) {
	token := c.Query("token")
	userid := c.Query("userid")
	user, err := CheckToken(token)
	if err != nil {
		TokenErrorRes(c, err)
	}
	c.JSON(200, handler.UserInfoHandle(user, userid))

}

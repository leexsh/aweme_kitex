package controller

import (
	"net/http"
	"sync/atomic"

	"github.com/gin-gonic/gin"
)

/*
用户
*/

var (
	usersLoginInfo = map[string]User{
		"caiXuKun": {
			Id:            1,
			Name:          "caiXuKun",
			FollowerCount: 5,
			FollowCount:   20,
			IsFollow:      true,
		},
	}

	userIdSequeue = int64(1)
)

type UserLoginResponse struct {
	Response
	UserId int64  `json:"user_id,omitempty"`
	Token  string `json:"token"`
}

type UserResponse struct {
	Response
	User User `json:"user"`
}

func Register(c *gin.Context) {
	userName := c.Query("username")
	password := c.Query("password")
	token := userName + password

	if _, exist := usersLoginInfo[userName]; exist {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{
				StatusCode: 0,
				StatusMsg:  "User already exist",
			},
		})
	} else {
		atomic.AddInt64(&userIdSequeue, 1)
		newUser := User{
			Id:   userIdSequeue,
			Name: userName,
		}
		usersLoginInfo[userName] = newUser
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{
				StatusCode: 0,
				StatusMsg:  "Success",
			},
			UserId: userIdSequeue,
			Token:  token,
		})
	}
}

func Login(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	token := username + password
	if user, exist := usersLoginInfo[username]; exist {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 0, StatusMsg: "success"},
			UserId:   user.Id,
			Token:    token,
		})
	} else {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User doesn't exist"},
		})
	}
}

func UserInfo(c *gin.Context) {
	name := c.Query("username")
	if user, exist := usersLoginInfo[name]; exist {
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{
				StatusCode: 0,
				StatusMsg:  "success",
			},
			User: user,
		})
	} else {
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{
				StatusCode: -1,
				StatusMsg:  "User doesn't exist",
			},
		})
	}
}

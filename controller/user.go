package controller

import (
	"aweme_kitex/models"
	"aweme_kitex/utils"
	"fmt"
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
			Identity:      "asdd",
			Name:          "caiXuKun",
			FollowerCount: 5,
			FollowCount:   20,
			IsFollow:      true,
		},
	}

	userIdSequeue = int64(1)

	u    = models.User{}
	last = models.User{}
)

type UserLoginResponse struct {
	Response
	UserId   string `json:"user_id,omitempty"`
	UserName string `json:"user_name,omitempty"`
	Token    string `json:"token,omitempty"`
}

type UserRegisterResponse struct {
	Response
	UserId string `json:"user_id,omitempty"`
}

type UserResponse struct {
	Response
	User User `json:"user"`
}

func Register(c *gin.Context) {
	userName := c.Query("username")
	password := c.Query("password")
	identity := c.Query("identity")

	res := models.DB.Where("name=?", userName).First(&u)
	fmt.Println(res.RowsAffected)
	if _, exist := usersLoginInfo[identity]; exist || u.Name == userName {
		c.JSON(http.StatusOK, UserRegisterResponse{
			Response: Response{
				StatusCode: 0,
				StatusMsg:  "User already exist",
			},
		})
	} else {
		atomic.AddInt64(&userIdSequeue, 1)
		newUser := User{
			Id:       userIdSequeue,
			Identity: utils.GenerateUUID(),
			Name:     userName,
		}
		usersLoginInfo[newUser.Identity] = newUser

		// insert to data
		newUserData := models.User{
			Identity:      utils.GenerateUUID(),
			Name:          userName,
			Password:      utils.Md5(password),
			FollowCount:   0,
			FollowerCount: 0,
			IsFollow:      false,
		}
		// models.DB.Select("identity", "name", "password", "follow_count", "follower_count", "is_follow").Create(&newUserData)
		res := models.DB.Create(&newUserData)
		if res.Error != nil {
			c.JSON(http.StatusOK, UserRegisterResponse{
				Response: Response{
					StatusCode: 0,
					StatusMsg:  fmt.Sprintf("write fail, err: %s", res.Error.Error()),
				},
				UserId: newUserData.Identity,
			})
			return
		}
		c.JSON(http.StatusOK, UserRegisterResponse{
			Response: Response{
				StatusCode: 0,
				StatusMsg:  "write success",
			},
			UserId: newUserData.Identity,
		})
	}
}

func Login(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")
	identity := c.Query("identity")

	models.DB.Where("name=? AND password=?", username, utils.Md5(password)).First(&u)
	token, err := utils.GenerateToken(u.Id, u.Name, utils.TokenExpire)
	if err != nil {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{
				StatusCode: 0,
				StatusMsg:  err.Error()}})
	}
	if user, exist := usersLoginInfo[identity]; exist {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 0, StatusMsg: "login success"},
			UserId:   user.Identity,
			UserName: user.Name,
			Token:    token,
		})
	} else if u.Name == username {
		// 写入缓存
		usersLoginInfo[u.Identity] = User{
			Id:            u.Id,
			Name:          u.Name,
			FollowerCount: u.FollowCount,
			FollowCount:   u.FollowerCount,
			IsFollow:      false,
		}
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 0, StatusMsg: "login success"},
			UserId:   u.Identity,
			UserName: u.Name,
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
				StatusMsg:  "get user info success",
			},
			User: user,
		})
		return
	} else {
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{
				StatusCode: -1,
				StatusMsg:  "User doesn't exist",
			},
		})
	}
}

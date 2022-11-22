package controller

import "github.com/gin-gonic/gin"

type UserListResponse struct {
	Response
	UserList []User `json:"user_list"`
}

/*
好友关系
*/

// check
func RelationAction(c *gin.Context) {
	name := c.Query("username")
	if _, ok := usersLoginInfo[name]; ok {
		c.JSON(200, Response{
			StatusCode: 200,
			StatusMsg:  "Success",
		})
	} else {
		c.JSON(200, Response{
			200,
			"User doesn't exist",
		})
	}
}

// 关注
func FollowList(c *gin.Context) {
	c.JSON(200, UserListResponse{
		Response{
			StatusCode: 200,
			StatusMsg:  "Success",
		},
		[]User{},
	})
}

// 粉丝
func FollowerList(c *gin.Context) {
	c.JSON(200, UserListResponse{
		Response{
			StatusCode: 0,
			StatusMsg:  "success",
		},
		[]User{},
	})
}

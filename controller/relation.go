package controller

import (
	"aweme_kitex/handler"

	"github.com/gin-gonic/gin"
)

/*
好友关系
*/

// check
func RelationAction(c *gin.Context) {
	toUserId := c.Query("to_user_id")
	action := c.Query("action")
	token := c.Query("token")
	user, err := CheckToken(token)
	if err != nil {
		TokenErrorRes(c, err)
	}
	c.JSON(200, handler.RelationActionHandle(user, toUserId, action))
}

// 获取关注
func FollowList(c *gin.Context) {
	token := c.Query("token")
	user, err := CheckToken(token)
	if err != nil {
		TokenErrorRes(c, err)
	}
	c.JSON(200, handler.ShowFollowListHandle(user))
}

// 获取粉丝
func FollowerList(c *gin.Context) {
	token := c.Query("token")
	user, err := CheckToken(token)
	if err != nil {
		TokenErrorRes(c, err)
	}
	c.JSON(200, handler.ShowFollowerListHandle(user))
}

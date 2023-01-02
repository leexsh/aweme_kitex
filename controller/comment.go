package controller

import (
	"aweme_kitex/handler"
	"aweme_kitex/models"

	"github.com/gin-gonic/gin"
)

/*
评论
*/

func CommentAction(c *gin.Context) {
	token := c.Query("token")
	user, err := CheckToken(token)
	if err != nil {
		TokenErrorRes(c, err)
	}
	videoId := c.Query("videoId")
	actionType := c.Query("actionType")
	if actionType == "1" {
		commentText := c.Query("content")
		c.JSON(200, handler.CreateCommentHandle(user, videoId, commentText))
	} else if actionType == "2" {
		commentId := c.Query("commentId")
		c.JSON(200, handler.DelCommentHandle(user, commentId))
	} else {
		c.JSON(200, models.Response{
			StatusCode: 0,
			StatusMsg:  "action type error",
		})
	}
}

func CommentList(c *gin.Context) {
	token := c.Query("token")
	user, err := CheckToken(token)
	if err != nil {
		TokenErrorRes(c, err)
	}

	videoId := c.Query("videoId")
	c.JSON(200, handler.CommentListHandle(user, videoId))
}

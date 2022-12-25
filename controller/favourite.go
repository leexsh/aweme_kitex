package controller

import (
	"aweme_kitex/handler"

	"github.com/gin-gonic/gin"
)

func FavouriteAction(c *gin.Context) {
	token := c.Query("token")
	user, err := CheckToken(token)
	if err != nil {
		TokenErrorRes(c, err)
	}

	videoId := c.Query("videoId")
	action := c.Query("actionType")

	c.JSON(200, handler.FavouriteActionHandle(user, videoId, action))
}

func FavouriteList(c *gin.Context) {
	token := c.Query("token")
	user, err := CheckToken(token)
	if err != nil {
		TokenErrorRes(c, err)
	}
	c.JSON(200, handler.FavouriteListHandle(user))
}

package handlers

import (
	"aweme_kitex/cmd/api/rpc"
	"aweme_kitex/cmd/publish/kitex_gen/publish"
	"aweme_kitex/controller"
	"aweme_kitex/handler"
	"aweme_kitex/pkg/errno"
	"context"

	"github.com/gin-gonic/gin"
)

func Publish(c *gin.Context) {
	token := c.Query("token")
	title := c.Query("title")
	user, err := controller.CheckToken(token)
	data, err := c.FormFile("data")
	if err != nil {
		SendResponse(c, errno.ConvertErr(err), nil)
		return
	}
	c.JSON(200, handler.PublishVideoHandle(user.Id, user.Name, title, data, c))
}

func PublishList(c *gin.Context) {
	token := c.Query("token")
	user, err := controller.CheckToken(token)
	if err != nil {
		SendResponse(c, errno.ConvertErr(err), nil)
		return
	}
	videoList, err := rpc.PublishVideoList(context.Background(), &publish.PublishListRequest{Token: token})
	if err != nil {
		SendResponse(c, errno.ConvertErr(err), nil)
		return
	}
	SendResponse(c, errno.Success, map[string]interface{}{"video": videoList})
	c.JSON(200, handler.QueryVideoPublishedHandle(user.Id))
}

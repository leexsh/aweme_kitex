package controller

import (
	"aweme_kitex/handler"

	"github.com/gin-gonic/gin"
)

/*
发布作品
	1.视频上传本地./public
	2.视频上传COS
	3.视频信息写入mysql
*/
func Publish(c *gin.Context) {
	token := c.Query("token")
	title := c.Query("title")
	user, err := CheckToken(token)
	data, err := c.FormFile("data")
	if err != nil {
		TokenErrorRes(c, err)
	}

	publishRes := handler.PublishVideoHandler(user.Id, user.Name, title, data)
	c.JSON(200, publishRes)
}

func PublishList(c *gin.Context) {
	token := c.Query("token")

	user, err := CheckToken(token)
	if err != nil {
		TokenErrorRes(c, err)
	}
	videoRes := handler.QueryVideoPublishedHandler(user.Id)
	c.JSON(200, videoRes)
}

package controller

import "github.com/gin-gonic/gin"

func FavouriteAction(c *gin.Context) {
	name := c.Query("username")

	if _, ok := usersLoginInfo[name]; ok {
		c.JSON(200, Response{
			StatusCode: 0,
			StatusMsg:  "success",
		})
	} else {
		c.JSON(200, Response{
			-1, "User doesn't exist",
		})
	}
}

func FavouriteList(c *gin.Context) {
	c.JSON(200, VideoListResponse{
		Response: Response{
			StatusCode: 0,
			StatusMsg:  "success",
		},
		VideoList: DemoVideos,
	})
}

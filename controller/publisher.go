package controller

import (
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

/*
发布作品
*/

func Publish(c *gin.Context) {
	token := c.Query("token")
	user, err := CheckToken(token)

	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"message": "User doesn't exist",
		})
		return
	}

	data, err := c.FormFile("data")
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusMsg:  err.Error(),
			StatusCode: -1,
		})
		return
	}
	fileName := filepath.Base(data.Filename)
	finalName := fmt.Sprintf("%s_%s", user.Name, fileName)
	saveFile := filepath.Join("./public/", finalName)
	if err := c.SaveUploadedFile(data, saveFile); err != nil {
		c.JSON(200, Response{
			-1,
			err.Error(),
		})
		return
	}
	c.JSON(200, Response{
		0,
		finalName + "uploaded successfully",
	})
}

func PublishList(c *gin.Context) {
	c.JSON(200, VideoListResponse{
		Response{
			StatusCode: 0,
			StatusMsg:  "success",
		},
		DemoVideos,
	})
}

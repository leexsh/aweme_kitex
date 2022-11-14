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

type VideoListResponse struct {
	Response
	VideoList []Video `json:"video_list,omitempty"`
}

func Publish(c *gin.Context) {
	name := c.Query("username")

	if _, ok := usersLoginInfo[name]; !ok {
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
	user := usersLoginInfo[name]
	finalName := fmt.Sprintf("%d_%s", user.Id, fileName)
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

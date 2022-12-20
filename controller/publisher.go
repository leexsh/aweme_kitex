package controller

import (
	"aweme_kitex/cfg"
	"aweme_kitex/handler"
	"aweme_kitex/model"
	"aweme_kitex/utils"
	"context"
	"fmt"
	"net/http"
	"path/filepath"

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
	user, err := CheckToken(token)
	if err != nil {
		TokenErrorRes(c, err)
	}

	data, err := c.FormFile("data")
	if err != nil {
		c.JSON(http.StatusOK, model.Response{
			StatusMsg:  err.Error(),
			StatusCode: -1,
		})
		return
	}
	fileName := filepath.Base(data.Filename)
	finalName := fmt.Sprintf("%s_%s", user.Name, fileName)

	saveFile := filepath.Join("./public/", finalName)
	// 1.save local
	if err := c.SaveUploadedFile(data, saveFile); err != nil {
		c.JSON(200, model.Response{
			-1,
			err.Error(),
		})
		return
	}

	// 2.upload COS
	cosKey := fmt.Sprintf("%s/%s", user.Name, finalName)
	_, _, err = cfg.COSClient.Object.Upload(
		context.Background(), cosKey, saveFile, nil,
	)
	if err != nil {
		c.JSON(200, model.Response{
			-1,
			err.Error(),
		})
		return
	}
	// 3.wirte to mysql
	ourl := cos.Object.GetObjectURL(cosKey)

	title := c.Query("title")
	video := model.VideoRawData{
		VideoId: utils.GenerateUUID(),
		UserId:  user.Id,
		Title:   title,
		PlayUrl: ourl.String(),
	}
	if err := db.Table("video").Debug().Create(&video).Error; err != nil {
		c.JSON(200, model.Response{
			-1,
			err.Error(),
		})
		return
	}
	c.JSON(200, model.Response{
		0,
		fmt.Sprintf("%s uploaded successfully", title),
	})
}

func PublishList(c *gin.Context) {
	token := c.Query("token")

	user, err := CheckToken(token)
	if err != nil {
		TokenErrorRes(c, err)
	}
	videoRes := handler.QueryVodeoPublishLishHandler(user.Id)
	c.JSON(200, videoRes)
}

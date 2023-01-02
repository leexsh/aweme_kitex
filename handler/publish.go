package handler

import (
	"aweme_kitex/models"
	"aweme_kitex/service"
	"mime/multipart"

	"github.com/gin-gonic/gin"
)

// --------- handler ---------------
// 该层功能包括处理传入参数，向service层请求
func QueryVideoPublishedHandle(userId string) *models.VideoListResponse {
	videos, err := service.QueryUserVideos(userId)
	if err != nil {
		return &models.VideoListResponse{
			models.Response{
				-1,
				err.Error(),
			},
			nil,
		}
	}
	return &models.VideoListResponse{
		models.Response{
			0, "success",
		},
		videos,
	}
}

func PublishVideoHandle(userId, userName, title string, data *multipart.FileHeader, c *gin.Context) *models.Response {
	if data == nil {
		return &models.Response{
			StatusCode: -1, StatusMsg: "publish info is null",
		}
	}
	if len(title) > 128 {
		return &models.Response{
			StatusCode: -1, StatusMsg: "title out of range",
		}
	}
	err := service.PublishVideoService(userId, userName, title, data, c)
	if err != nil {
		return &models.Response{
			-1,
			err.Error(),
		}
	}
	return &models.Response{
		StatusCode: 0, StatusMsg: "publish video success",
	}
}

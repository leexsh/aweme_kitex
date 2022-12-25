package handler

import (
	"aweme_kitex/model"
	"aweme_kitex/service"
	"mime/multipart"
)

// --------- handler ---------------
// 该层功能包括处理传入参数，向service层请求
func QueryVideoPublishedHandle(userId string) model.VideoListResponse {
	videos, err := service.QueryUserVideos(userId)
	if err != nil {
		return model.VideoListResponse{
			model.Response{
				-1,
				err.Error(),
			},
			nil,
		}
	}
	return model.VideoListResponse{
		model.Response{
			0, "success",
		},
		videos,
	}
}

func PublishVideoHandle(userId, userName, title string, data *multipart.FileHeader) model.Response {
	if data == nil {
		return model.Response{
			StatusCode: -1, StatusMsg: "publish info is null",
		}
	}
	if len(title) > 128 {
		return model.Response{
			StatusCode: -1, StatusMsg: "title out of range",
		}
	}
	err := service.PublishVideoService(userId, userName, title, data)
	if err != nil {
		return model.Response{
			-1,
			err.Error(),
		}
	}
	return model.Response{
		StatusCode: 0, StatusMsg: "publish video success",
	}
}

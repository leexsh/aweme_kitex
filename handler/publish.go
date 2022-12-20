package handler

import (
	"aweme_kitex/model"
	"aweme_kitex/service"
)

// --------- handler ---------------
// 该层功能包括处理传入参数，向service层请求
func QueryVodeoPublishLishHandler(userId string) model.VideoListResponse {
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

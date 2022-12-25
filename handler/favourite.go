package handler

import (
	"aweme_kitex/model"
	"aweme_kitex/service"
)

// --- 关注和收藏列表----

func FavouriteActionHandle(user *model.UserClaim, videoId, action string) *model.Response {
	err := service.FavouriteActionService(user, videoId, action)
	if err != nil {
		return &model.Response{
			-1,
			err.Error(),
		}
	}
	return &model.Response{
		0,
		"favourite action success",
	}
}

func FavouriteListHandle(user *model.UserClaim) *model.VideoListResponse {
	videos, err := service.FavouriteListService(user.Id, user.Name)
	if err != nil {
		return &model.VideoListResponse{
			Response: model.Response{
				-1, err.Error(),
			},
		}
	}
	return &model.VideoListResponse{
		model.Response{
			0, "get favourite list success",
		},
		videos,
	}
}

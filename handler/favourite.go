package handler

import (
	"aweme_kitex/models"
	"aweme_kitex/service"
)

// --- 关注和收藏列表----

func FavouriteActionHandle(user *models.UserClaim, videoId, action string) *models.Response {
	err := service.FavouriteActionService(user, videoId, action)
	if err != nil {
		return &models.Response{
			-1,
			err.Error(),
		}
	}
	return &models.Response{
		0,
		"favourite action success",
	}
}

func FavouriteListHandle(user *models.UserClaim) *models.VideoListResponse {
	videos, err := service.FavouriteListService(user.Id, user.Name)
	if err != nil {
		return &models.VideoListResponse{
			Response: models.Response{
				-1, err.Error(),
			},
		}
	}
	return &models.VideoListResponse{
		models.Response{
			0, "get favourite list success",
		},
		videos,
	}
}

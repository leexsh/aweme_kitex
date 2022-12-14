package handler

import (
	"aweme_kitex/models"
	"aweme_kitex/pkg/jwt"
	"aweme_kitex/service"
)

func RelationActionHandle(user *jwt.UserClaim, toUser, action string) *models.Response {
	if action != "1" && action != "2" {
		return &models.Response{
			-1, "action type error",
		}
	}
	err := service.RelationAction(user.Id, toUser, action)
	if err != nil {
		return &models.Response{
			-1, err.Error(),
		}
	}
	return &models.Response{
		0,
		"relation action success",
	}
}

func ShowFollowListHandle(u *jwt.UserClaim) *models.RelationListhResponse {
	userList, err := service.GetFollowList(u.Id)
	if err != nil {
		return &models.RelationListhResponse{
			Response: models.Response{
				-1, err.Error(),
			},
		}
	}
	return &models.RelationListhResponse{
		models.Response{
			0,
			"get follow list succes",
		},
		userList,
	}
}

func ShowFollowerListHandle(u *jwt.UserClaim) *models.RelationListhResponse {
	userList, err := service.GetFollowerList(u.Id)
	if err != nil {
		return &models.RelationListhResponse{
			Response: models.Response{
				-1, err.Error(),
			},
		}
	}
	return &models.RelationListhResponse{
		models.Response{
			0,
			"get follower list succes",
		},
		userList,
	}
}

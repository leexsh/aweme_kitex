package handler

import (
	"aweme_kitex/models"
	"aweme_kitex/pkg/jwt"
	"aweme_kitex/service"
)

func CreateCommentHandle(user *jwt.UserClaim, videoId, content string) *models.CommentActionResponse {
	if len(content) > 512 {
		return &models.CommentActionResponse{
			Response: models.Response{
				-1, "content too large",
			},
		}
	}
	commet, err := service.CreateComment(user.Id, videoId, content, "")
	if err != nil {
		return &models.CommentActionResponse{
			Response: models.Response{
				-1, err.Error(),
			},
		}
	}
	return &models.CommentActionResponse{
		models.Response{0,
			"create comment success",
		},
		*commet,
	}
}

func DelCommentHandle(user *jwt.UserClaim, commentId string) *models.CommentActionResponse {
	commet, err := service.DelComment(user.Id, commentId)
	if err != nil {
		return &models.CommentActionResponse{
			Response: models.Response{
				-1, err.Error(),
			},
		}
	}
	return &models.CommentActionResponse{
		models.Response{0,
			"delete comment success",
		},
		*commet,
	}
}

func CommentListHandle(user *jwt.UserClaim, videoId string) *models.CommentListResponse {
	res, err := service.ShowCommentList(user.Id, videoId)
	if err != nil {
		return &models.CommentListResponse{
			Response: models.Response{
				-1, err.Error(),
			},
		}
	}
	return &models.CommentListResponse{
		models.Response{0,
			"get comment list success",
		},
		res,
	}
}

package handler

import (
	"aweme_kitex/model"
	"aweme_kitex/service"
)

func CreateCommentHandle(user *model.UserClaim, videoId, content string) *model.CommentActionResponse {
	if len(content) > 512 {
		return &model.CommentActionResponse{
			Response: model.Response{
				-1, "content too large",
			},
		}
	}
	commet, err := service.CreateComment(user.Id, videoId, content, "")
	if err != nil {
		return &model.CommentActionResponse{
			Response: model.Response{
				-1, err.Error(),
			},
		}
	}
	return &model.CommentActionResponse{
		model.Response{0,
			"create comment success",
		},
		*commet,
	}
}

func DelCommentHandle(user *model.UserClaim, videoId, commentId string) *model.CommentActionResponse {
	commet, err := service.DelComment(user.Id, commentId)
	if err != nil {
		return &model.CommentActionResponse{
			Response: model.Response{
				-1, err.Error(),
			},
		}
	}
	return &model.CommentActionResponse{
		model.Response{0,
			"delete comment success",
		},
		*commet,
	}
}

func CommentListHandle(user *model.UserClaim, videoId string) *model.CommentListResponse {
	res, err := service.ShowCommentList(user.Id, videoId)
	if err != nil {
		return &model.CommentListResponse{
			Response: model.Response{
				-1, err.Error(),
			},
		}
	}
	return &model.CommentListResponse{
		model.Response{0,
			"delete comment success",
		},
		res,
	}
}

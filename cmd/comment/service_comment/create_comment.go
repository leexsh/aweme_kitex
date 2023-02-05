package service_comment

import (
	"aweme_kitex/cmd/comment/kitex_gen/comment"
	"aweme_kitex/cmd/comment/kitex_gen/user"
	"aweme_kitex/pkg/jwt"
	"aweme_kitex/service"
	"context"
)

type CreateCommentService struct {
	ctx context.Context
}

// NewCreateCommentService new CreateCommentService
func NewCreateCommentService(ctx context.Context) *CreateCommentService {
	return &CreateCommentService{ctx: ctx}
}

// CreateComment add comment
func (s *CreateCommentService) CreateComment(req *comment.CommentActionRequest) (*comment.Comment, error) {
	uc, err := s.CheckToken(req.Token)
	if err != nil {
		return nil, err
	}
	commet, err := service.CreateComment(uc.Id, req.VideoId, *req.CommentContent, "")
	curComment := &comment.Comment{
		CommentId: commet.Id,
		User: &user.User{
			UserId:        commet.User.UserId,
			Name:          commet.User.Name,
			FollowCount:   commet.User.FollowCount,
			FollowerCount: commet.User.FollowerCount,
			IsFollow:      commet.User.IsFollow,
		},
		Content:    commet.Content,
		CreateTime: commet.CreateDate,
	}
	return curComment, nil

}

func (s *CreateCommentService) CheckToken(token string) (*jwt.UserClaim, error) {
	uc, err := jwt.AnalyzeToken(token)
	if err != nil {
		return nil, err
	}
	return uc, nil
}

package service_comment

import (
	"aweme_kitex/cmd/comment/kitex_gen/comment"
	"aweme_kitex/pkg/jwt"
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
	uc, err := jwt.AnalyzeToken(req.Token)
	if err != nil {
		return nil, err
	}
	return createComment(s.ctx, uc.Id, req.VideoId, *req.CommentContent, "")
}

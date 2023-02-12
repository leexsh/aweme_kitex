package service_comment

import (
	"aweme_kitex/cmd/comment/kitex_gen/comment"
	"aweme_kitex/pkg/jwt"
	"context"
)

type DeleteCommentService struct {
	ctx context.Context
}

// NewDeleteCommentService new DeleteCommentService
func NewDeleteCommentService(ctx context.Context) *DeleteCommentService {
	return &DeleteCommentService{ctx: ctx}
}

func (s *DeleteCommentService) DelComment(req *comment.CommentActionRequest) (*comment.Comment, error) {
	uc, err := jwt.AnalyzeToken(req.Token)
	if err != nil {
		return nil, err
	}
	return delComment(s.ctx, uc.Id, req.CommentId)
}

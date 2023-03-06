package main

import (
	comment "aweme_kitex/cmd/comment/kitex_gen/comment"
	"aweme_kitex/cmd/comment/service_comment"
	commentPack "aweme_kitex/cmd/comment/service_comment/pack"
	"context"
	"unicode/utf8"
)

// CommentServiceImpl implements the last service_user interface defined in the IDL.
type CommentServiceImpl struct{}

// CreateComment implements the CommentServiceImpl interface.
func (s *CommentServiceImpl) CreateComment(ctx context.Context, req *comment.CommentActionRequest) (resp *comment.CommentActionResponse, err error) {
	if utf8.RuneCountInString(*req.CommentContent) > 20 {
		return commentPack.PackCommentAction(-1, "params error", nil), nil
	}

	comm, err := service_comment.NewCreateCommentService(ctx).CreateComment(req)
	if err != nil {
		return commentPack.PackCommentAction(-1, err.Error(), nil), nil
	}
	return commentPack.PackCommentAction(0, "create comment success", []*comment.Comment{comm}), nil
}

// DelComment implements the CommentServiceImpl interface.
func (s *CommentServiceImpl) DelComment(ctx context.Context, req *comment.CommentActionRequest) (resp *comment.CommentActionResponse, err error) {
	if utf8.RuneCountInString(*req.CommentContent) > 20 {
		return commentPack.PackCommentAction(-1, "params error", nil), nil
	}

	comm, err := service_comment.NewDeleteCommentService(ctx).DelComment(req)
	if err != nil {
		return commentPack.PackCommentAction(-1, err.Error(), nil), nil
	}
	return commentPack.PackCommentAction(0, "create comment success", []*comment.Comment{comm}), nil
}

// CommentList implements the CommentServiceImpl interface.
func (s *CommentServiceImpl) CommentList(ctx context.Context, req *comment.CommentListRequest) (resp *comment.CommentListResponse, err error) {
	if len(req.VideoId) == 0 {
		return commentPack.PackCommentList(-1, "params error", nil), nil
	}

	commentList, err := service_comment.NewCommentListService(ctx).CommentList(req)
	if err != nil {
		return commentPack.PackCommentList(-1, err.Error(), nil), nil
	}
	return commentPack.PackCommentList(0, "create comment success", commentList), nil
}

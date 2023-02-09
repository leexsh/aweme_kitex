package main

import (
	comment "aweme_kitex/cmd/comment/kitex_gen/comment"
	"aweme_kitex/cmd/comment/service_comment"
	"context"
	"time"
	"unicode/utf8"
)

// CommentServiceImpl implements the last service_user interface defined in the IDL.
type CommentServiceImpl struct{}

// CreateComment implements the CommentServiceImpl interface.
func (s *CommentServiceImpl) CreateComment(ctx context.Context, req *comment.CommentActionRequest) (resp *comment.CommentActionResponse, err error) {
	// TODO: Your code here...
	resp = new(comment.CommentActionResponse)
	if len(req.Token) == 0 || utf8.RuneCountInString(*req.CommentContent) > 20 {
		resp.BaseResp.ServiceTime = time.Now().Unix()
		resp.BaseResp.StatusCode = -1
		resp.BaseResp.StatusMsg = "params error"
		return resp, nil
	}

	comment, err := service_comment.NewCreateCommentService(ctx).CreateComment(req)
	if err != nil {
		resp.BaseResp.ServiceTime = time.Now().Unix()
		resp.BaseResp.StatusCode = -1
		resp.BaseResp.StatusMsg = err.Error()
		return resp, nil
	}
	resp.BaseResp.ServiceTime = time.Now().Unix()
	resp.BaseResp.StatusCode = 0
	resp.BaseResp.StatusMsg = "create comment success"
	resp.CommentList[0] = comment
	return resp, nil
}

// DelComment implements the CommentServiceImpl interface.
func (s *CommentServiceImpl) DelComment(ctx context.Context, req *comment.CommentActionRequest) (resp *comment.CommentActionResponse, err error) {
	// TODO: Your code here...
	resp = new(comment.CommentActionResponse)
	if len(req.Token) == 0 || utf8.RuneCountInString(*req.CommentContent) > 20 {
		resp.BaseResp.ServiceTime = time.Now().Unix()
		resp.BaseResp.StatusCode = -1
		resp.BaseResp.StatusMsg = "params error"
		return resp, nil
	}

	commet, err := service_comment.NewDeleteCommentService(ctx).DelComment(req)
	if err != nil {
		resp.BaseResp.ServiceTime = time.Now().Unix()
		resp.BaseResp.StatusCode = -1
		resp.BaseResp.StatusMsg = err.Error()
		return resp, nil
	}
	resp.BaseResp.ServiceTime = time.Now().Unix()
	resp.BaseResp.StatusCode = 0
	resp.BaseResp.StatusMsg = "create comment success"
	resp.CommentList[0] = commet
	return resp, nil
}

// CommentList implements the CommentServiceImpl interface.
func (s *CommentServiceImpl) CommentList(ctx context.Context, req *comment.CommentListRequest) (resp *comment.CommentListResponse, err error) {
	// TODO: Your code here...
	resp = new(comment.CommentListResponse)

	if len(req.Token) == 0 || len(req.VideoId) == 0 {
		resp.BaseResp.ServiceTime = time.Now().Unix()
		resp.BaseResp.StatusCode = -1
		resp.BaseResp.StatusMsg = "params error"
		return resp, nil
	}

	commentList, err := service_comment.NewCommentListService(ctx).CommentList(req)
	if err != nil {
		resp.BaseResp.ServiceTime = time.Now().Unix()
		resp.BaseResp.StatusCode = -1
		resp.BaseResp.StatusMsg = err.Error()
		return resp, nil
	}
	resp.BaseResp.ServiceTime = time.Now().Unix()
	resp.BaseResp.StatusCode = 0
	resp.BaseResp.StatusMsg = "get comment List success"
	resp.CommentList = commentList
	return
}

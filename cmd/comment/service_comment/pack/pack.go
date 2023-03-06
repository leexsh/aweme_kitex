package commentPack

import (
	"aweme_kitex/cmd/comment/kitex_gen/comment"
	"time"
)

func PackCommentAction(code int64, msg string, commentList []*comment.Comment) (resp *comment.CommentActionResponse) {
	resp = new(comment.CommentActionResponse)
	resp.BaseResp.StatusCode = code
	resp.BaseResp.StatusMsg = msg
	resp.BaseResp.ServiceTime = time.Now().Unix()
	if commentList != nil {
		resp.CommentList = commentList
	}
	return
}

func PackCommentList(code int64, msg string, commentList []*comment.Comment) (resp *comment.CommentListResponse) {
	resp = new(comment.CommentListResponse)
	resp.BaseResp.StatusCode = code
	resp.BaseResp.StatusMsg = msg
	resp.BaseResp.ServiceTime = time.Now().Unix()
	if commentList != nil {
		resp.CommentList = commentList
	}
	return
}

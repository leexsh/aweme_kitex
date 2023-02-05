package handlers

import (
	"aweme_kitex/cmd/api/rpc"
	"aweme_kitex/cmd/comment/kitex_gen/comment"
	"aweme_kitex/pkg/errno"
	"context"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
)

// CommentAction implement adding and deleting comments
func CommentAction(c *gin.Context) {
	token := c.Query("token")
	videoId := c.Query("video_id")
	actionType := c.Query("action_type")

	if actionType == "1" {
		commentText := c.Query("comment_text")

		if len := utf8.RuneCountInString(commentText); len > 20 {
			// SendResponse(c, errno., nil)
			return
		}
		req := &comment.CommentActionRequest{
			Token:          token,
			VideoId:        videoId,
			ActionType:     actionType,
			CommentContent: &commentText,
		}
		comment, err := rpc.CreateComment(context.Background(), req)

		if err != nil {
			SendResponse(c, errno.ConvertErr(err), nil)
			return
		}
		SendResponse(c, errno.Success, map[string]interface{}{"comment": comment})

	} else if actionType == "2" {
		commentIdStr := c.Query("comment_id")

		req := &comment.CommentActionRequest{
			Token:      token,
			VideoId:    videoId,
			ActionType: actionType,
			CommentId:  commentIdStr,
		}
		comment, err := rpc.DeleteComment(context.Background(), req)
		if err != nil {
			SendResponse(c, errno.ConvertErr(err), nil)
			return
		}
		SendResponse(c, errno.Success, map[string]interface{}{"comment": comment})

	} else {
		SendResponse(c, errno.ParamErr, nil)
	}
}

// CommentList get comment list info
func CommentList(c *gin.Context) {
	token := c.Query("token")
	videoId := c.Query("video_id")

	req := &comment.CommentListRequest{Token: token, VideoId: videoId}

	commentList, err := rpc.CommentList(context.Background(), req)
	if err != nil {
		SendResponse(c, errno.ConvertErr(err), nil)
		return
	}

	SendResponse(c, errno.Success, map[string]interface{}{"commentList": commentList})
}

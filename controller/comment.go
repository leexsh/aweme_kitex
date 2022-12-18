package controller

import (
	"aweme_kitex/utils"
	"fmt"

	"github.com/gin-gonic/gin"
)

/*
评论
*/

func CommentAction(c *gin.Context) {
	token := c.Query("token")
	user, err := CheckToken(token)
	if err != nil {
		c.JSON(200, Response{
			-1,
			fmt.Sprintf("occur err:%s", err.Error()),
		})
	}
	actionType := c.Query("actionType")
	commentText := c.Query("content")
	videoId := c.Query("videoId")
	commentId := c.Query("commentId")
	_ = c.Query("commetCount")

	if actionType == "1" {
		// add comment
		newComment := &CommetRaw{
			Id:      utils.GenerateUUID(),
			UserId:  user.Id,
			VideoId: videoId,
			Content: commentText,
		}
		// 1. video comment +1
		video := VideoRawData{}
		db.Table("video").Debug().Where("video_id=?", videoId).Find(&video)
		db.Table("video").Where("video_id=?", videoId).Update("comment_count", video.CommentCount+1)

		// 2.comment table +1
		db.Table("comment").Debug().Create(newComment)
		c.JSON(200, newComment)

	} else if actionType == "2" {
		// delete comment
		video := VideoRawData{}
		db.Table("video").Debug().Where("video_id=?", videoId).Find(&video)
		db.Table("video").Where("video_id=?", videoId).Update("comment_count", video.CommentCount-1)
		db.Table("comment").Where("comment_id=?", commentId).Delete(&Comment{})
		c.JSON(200, Response{
			StatusCode: 0,
			StatusMsg:  "delete comment success",
		})
	}

}

func CommentList(c *gin.Context) {
	token := c.Query("token")
	_, err := utils.AnalyzeToke(token)
	if err != nil {
		c.JSON(200, Response{
			-1,
			fmt.Sprintf("occur err:%s", err.Error()),
		})
	}
	videoId := c.Query("videoId")
	commentRawList := make([]CommetRaw, 0)
	db.Table("comment").Debug().Where("video_id=?", videoId).Find(&commentRawList)
	commentList := make([]Comment, len(commentRawList))
	for i, comment := range commentRawList {
		comm := Comment{
			Id:      comment.Id,
			UserId:  comment.UserId,
			VideoId: comment.VideoId,
			Content: comment.Content,
		}
		commentList[i] = comm
	}
	c.JSON(200, CommentListResponse{
		Response: Response{
			StatusCode: 0,
			StatusMsg:  "get comments success",
		},
		CommentList: commentList,
	})
}

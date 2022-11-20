package controller

import (
	"aweme_kitex/models"
	"aweme_kitex/utils"

	"github.com/gin-gonic/gin"
)

/*
评论
*/

type CommentListResponse struct {
	Response
	CommentList []Comment `json:"comment_list,omitempty"`
}

func CommentAction(c *gin.Context) {
	userId := c.Query("userId")
	actionType := c.Query("actionType")
	commentText := c.Query("content")
	videoId := c.Query("videoId")
	commentId := c.Query("commentId")
	// commentCount, _ := strconv.Atoi(c.Query("commentCount"))
	// todo: use token
	_ = c.Query("token")
	if actionType == "1" {
		// add comment
		newComment := &models.Comment{
			Id:      utils.GenerateUUID(),
			UserId:  userId,
			VideoId: videoId,
			Content: commentText,
		}
		// 1. video comment +1
		video := models.Video{}
		models.DB.Table("video").Debug().Where("video_id=?", videoId).Find(&video)
		models.DB.Table("video").Where("video_id=?", videoId).Update("comment_count", video.CommentCount+1)

		// 2.comment table +1
		models.DB.Table("comment").Debug().Create(newComment)
		c.JSON(200, newComment)

	} else if actionType == "2" {
		// delete comment
		video := models.Video{}
		models.DB.Table("video").Debug().Where("video_id=?", videoId).Find(&video)
		models.DB.Table("video").Where("video_id=?", videoId).Update("comment_count", video.CommentCount-1)
		db.Table("comment").Where("comment_id=?", commentId).Delete(&models.Comment{})
		c.JSON(200, Response{
			StatusCode: 0,
			StatusMsg:  "delete comment success",
		})
	}

}

func CommentList(c *gin.Context) {
	videoId := c.Query("videoId")
	commentList := make([]Comment, 0)
	db.Table("comment").Debug().Where("video_id=?", videoId).Find(&commentList)
	c.JSON(200, CommentListResponse{
		Response: Response{
			StatusCode: 0,
			StatusMsg:  "get comments success",
		},
		CommentList: commentList,
	})
}

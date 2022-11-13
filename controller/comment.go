package controller

import "github.com/gin-gonic/gin"

/*
评论
*/

type CommentListResponse struct {
	Response
	CommentList []Comment `json:"comment_list,omitempty"`
}

func CommentAction(c *gin.Context) {
	name := c.Query("name")
	if _, exist := usersLoginInfo[name]; exist {
		c.JSON(200, Response{StatusCode: 0})
	} else {
		c.JSON(200, Response{StatusCode: -1, StatusMsg: "User doestn't exist"})
	}
}

func CommentList(c *gin.Context) {
	c.JSON(200, CommentListResponse{
		Response: Response{
			StatusCode: 0,
			StatusMsg:  "success",
		},
		CommentList: DemoComments,
	})
}

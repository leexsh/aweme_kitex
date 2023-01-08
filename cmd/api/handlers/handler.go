package handlers

import (
	"aweme_kitex/pkg/errno"

	"github.com/gin-gonic/gin"
)

func SendResponse(c *gin.Context, err error, data interface{}) {
	Err := errno.ConvertErr(err)
	c.JSON(200, Response{
		Err.ErrCode,
		Err.ErrMsg,
		data,
	})
}

type Response struct {
	Code    int64       `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type FeedRequest struct {
	Token      string `json:"token" form:"token"`
	LatestTime int64  `json:"latest_time" form:"latest_time"`
}

type UserLoginParam struct {
	UserName string `form:"username" json:"userName", binding:"required"`
	Password string `form:"Password" json:"Password", binding:"required"`
}

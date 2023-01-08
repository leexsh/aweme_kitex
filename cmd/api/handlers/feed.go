package handlers

import (
	"aweme_kitex/cmd/api/rpc"
	"aweme_kitex/cmd/feed/kitex_gen/feed"
	"aweme_kitex/models"
	constants "aweme_kitex/pkg/constant"
	"aweme_kitex/pkg/errno"
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func Feed(c *gin.Context) {
	var FeedVar FeedRequest
	token := c.DefaultQuery("token", "")
	defaultTime := time.Now().UnixMilli()
	FeedVar.Token = token
	FeedVar.LatestTime = defaultTime
	req := &feed.FeedRequest{
		LatestTime: defaultTime,
		Token:      token,
	}
	video, nextTime, err := rpc.Feed(context.Background(), req)
	if err != nil {
		SendResponse(c, err, nil)
		return
	}
	SendResponse(c, errno.Success, map[string]interface{}{constants.VideoList: video, constants.NextTime: nextTime})
	c.JSON(http.StatusOK, &models.FeedResponse{})
}

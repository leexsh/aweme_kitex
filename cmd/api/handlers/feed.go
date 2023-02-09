package handlers

import (
	"aweme_kitex/cmd/api/rpc"
	"aweme_kitex/cmd/feed/kitex_gen/feed"
	"aweme_kitex/models"
	constants "aweme_kitex/pkg/constant"
	"aweme_kitex/pkg/errno"
	"aweme_kitex/pkg/jwt"
	"aweme_kitex/utils"
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func Feed(c *gin.Context) {
	var FeedVar FeedRequest
	token := c.DefaultQuery("token", "")
	_, err := jwt.AnalyzeToken(token)
	if err != nil {
		SendResponse(c, errno.TokenInvalidErr, nil)
		return
	}
	defaultTimeStr := strconv.Itoa(int(utils.GetUnix()))
	latestTime, err := strconv.ParseInt(defaultTimeStr, 10, 64)
	if err != nil {
		SendResponse(c, errno.ParamErr, nil)
		return
	}
	FeedVar.Token = token
	FeedVar.LatestTime = latestTime
	req := &feed.FeedRequest{
		LatestTime: latestTime,
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

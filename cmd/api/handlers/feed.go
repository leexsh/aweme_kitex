package handlers

import (
	"aweme_kitex/cmd/api/rpc"
	"aweme_kitex/cmd/feed/kitex_gen/feed"
	"aweme_kitex/pkg/errno"
	"aweme_kitex/pkg/jwt"
	"aweme_kitex/pkg/utils"
	"context"
	"strconv"

	"github.com/gin-gonic/gin"
)

func Feed(c *gin.Context) {
	token := c.DefaultQuery("token", "")
	_, err := jwt.AnalyzeToken(token)
	if err != nil {
		SendResponse(c, errno.TokenInvalidErr)
		return
	}
	defaultTimeStr := strconv.Itoa(int(utils.GetUnix()))
	latestTime, err := strconv.ParseInt(defaultTimeStr, 10, 64)
	if err != nil {
		SendResponse(c, errno.ParamErr)
		return
	}
	req := &feed.FeedRequest{
		LatestTime: latestTime,
		Token:      token,
	}
	video, nextTime, err := rpc.Feed(context.Background(), req)
	if err != nil {
		SendResponse(c, err)
		return
	}
	SendFeedResponse(c, errno.Success, video, nextTime)
}

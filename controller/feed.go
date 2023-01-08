package controller

import (
	"aweme_kitex/handler"
	"aweme_kitex/pkg/logger"
	"aweme_kitex/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ----- controller------
// Feed same demo video list for every request
// 该层功能包括获取传入参数，向handler获取视频信息，返回响应信息
func Feed(c *gin.Context) {
	token := c.DefaultQuery("token", defaultToken)
	defaultTimeStr := strconv.Itoa(int(utils.GetUnix()))
	user, err := CheckToken(token)
	if err != nil {
		TokenErrorRes(c, err)
	}
	feedRes := handler.QueryVideoFeedHandler(user.Id, defaultTimeStr)
	logger.Info(&feedRes)
	c.JSON(http.StatusOK, feedRes)
}

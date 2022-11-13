package controller

import (
	"aweme_kitex/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

/* feed 流
 */

type FeedResponse struct {
	Response
	VideoList []Video `json:"video_list,omitempty"`
	NextTime  int64   `json:"next_time,omitempty"`
}

func Feed(c *gin.Context) {
	c.JSON(http.StatusOK, FeedResponse{
		Response: Response{StatusCode: 0,
			StatusMsg: "success"},
		VideoList: DemoVideos,
		NextTime:  utils.GetUnix(),
	})
}

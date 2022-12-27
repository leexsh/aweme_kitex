package handler

import (
	"aweme_kitex/models"
	"strconv"

	"aweme_kitex/service"
)

// --------- handler ---------------
// 该层功能包括处理传入参数，向service层获取视频信息，封装成响应信息

func QueryVideoFeedHandler(userId string, latestTimeStr string) models.FeedResponse {
	// 1.处理传人参数
	latestTime, err := strconv.ParseInt(latestTimeStr, 10, 64)
	if err != nil {
		return models.FeedResponse{
			Response: models.Response{
				StatusCode: -1,
				StatusMsg:  err.Error(),
			},
		}
	}

	// 2.ge video
	videoList, nextTime, err := service.QueryVideoData(latestTime, userId)
	if err != nil {
		return models.FeedResponse{
			Response: models.Response{
				StatusCode: -1,
				StatusMsg:  err.Error(),
			},
		}
	}
	return models.FeedResponse{
		models.Response{
			StatusCode: 0,
			StatusMsg:  "获取video成功",
		},
		videoList,
		nextTime,
	}
}

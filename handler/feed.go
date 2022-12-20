package handler

import (
	"strconv"

	"aweme_kitex/model"
	"aweme_kitex/service"
)

// --------- handler ---------------
// 该层功能包括处理传入参数，向service层获取视频信息，封装成响应信息

func QueryVideoFeedHandler(userId string, latestTimeStr string) model.FeedResponse {
	// 1.处理传人参数
	latestTime, err := strconv.ParseInt(latestTimeStr, 10, 64)
	if err != nil {
		return model.FeedResponse{
			Response: model.Response{
				StatusCode: -1,
				StatusMsg:  err.Error(),
			},
		}
	}

	// 2.ge video
	videoList, nextTime, err := service.QueryVideoData(latestTime, userId)
	if err != nil {
		return model.FeedResponse{
			Response: model.Response{
				StatusCode: -1,
				StatusMsg:  err.Error(),
			},
		}
	}
	return model.FeedResponse{
		model.Response{
			StatusCode: 0,
			StatusMsg:  "获取video成功",
		},
		videoList,
		nextTime,
	}
}

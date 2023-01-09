package serviceRPC

import (
	"aweme_kitex/cmd/feed/kitex_gen/feed"
	"aweme_kitex/cmd/feed/kitex_gen/user"
	"aweme_kitex/models"
	"time"
)

// VideoInfo pack video list info for rpc
func PackRPCVideoInfo(currentId string, videoData []*models.VideoRawData, userMap map[string]*models.UserRawData,
	favoriteMap map[string]*models.FavouriteRaw, relationMap map[string]*models.RelationRaw) ([]*feed.Video, int64) {
	videoList := make([]*feed.Video, 0)
	var nextTime int64
	for _, video := range videoData {
		videoUser, ok := userMap[video.UserId]
		if !ok {
			videoUser = &models.UserRawData{
				Name:          "未知用户",
				FollowCount:   0,
				FollowerCount: 0,
			}
			videoUser.UserId = ""
		}

		var isFavorite bool = false
		var isFollow bool = false

		videoList = append(videoList, &feed.Video{
			VideoId: video.VideoId,
			Author: &user.User{
				UserId:        videoUser.UserId,
				Name:          videoUser.Name,
				FollowCount:   videoUser.FollowCount,
				FollowerCount: videoUser.FollowerCount,
				IsFollow:      isFollow,
			},
			PlayUrl:        video.PlayUrl,
			CoverUrl:       video.PlayUrl,
			FavouriteCount: video.FavouriteCount,
			CommentCount:   video.CommentCount,
			IsFavourite:    isFavorite,
			Title:          video.Title,
		})
	}

	if len(videoData) == 0 {
		nextTime = time.Now().UnixMilli()
	} else {
		nextTime = videoData[len(videoData)-1].UpdatedAt.UnixMilli()
	}

	return videoList, nextTime
}

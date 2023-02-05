package serviceRPC

import (
	"aweme_kitex/cmd/publish/kitex_gen/feed"
	user2 "aweme_kitex/cmd/publish/kitex_gen/user"
	"aweme_kitex/models"
)

func PublishInfo(currentId string, videoData []*models.VideoRawData, userMap map[string]*models.UserRawData,
	favoriteMap map[string]*models.FavouriteRaw, relationMap map[string]*models.RelationRaw) []*feed.Video {
	videoList := make([]*feed.Video, 0)
	for _, video := range videoData {
		videoUser, ok := userMap[video.UserId]
		if !ok {
			videoUser = &models.UserRawData{
				Name:          "unknow",
				FollowCount:   0,
				FollowerCount: 0,
			}
			videoUser.UserId = "0"
		}

		var isFavorite bool = false
		var isFollow bool = false

		if ok {
			isFavorite = true
		}
		_, ok = relationMap[video.UserId]
		if ok {
			isFollow = true
		}
		videoList = append(videoList, &feed.Video{
			VideoId: video.VideoId,
			Author: &user2.User{
				UserId:        videoUser.UserId,
				Name:          videoUser.Name,
				FollowCount:   videoUser.FollowCount,
				FollowerCount: videoUser.FollowerCount,
				IsFollow:      isFollow,
			},
			PlayUrl:        video.PlayUrl,
			CoverUrl:       video.CoverUrl,
			FavouriteCount: video.FavouriteCount,
			CommentCount:   video.CommentCount,
			IsFavourite:    isFavorite,
			Title:          video.Title,
		})
	}

	return videoList
}

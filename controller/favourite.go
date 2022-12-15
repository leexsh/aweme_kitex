package controller

import (
	"aweme_kitex/models"
	"aweme_kitex/utils"
	"fmt"

	"github.com/gin-gonic/gin"
)

func FavouriteAction(c *gin.Context) {
	token := c.Query("token")
	user, err := utils.AnalyzeToke(token)
	if err != nil {
		c.JSON(200, Response{
			-1,
			fmt.Sprintf("occur err:%s", err.Error()),
		})
	}
	videoId := c.Query("videoId")
	action := c.Query("actionType")
	favour := Favourite{}
	video := VideoRawData{}
	db.Table("favourite").Debug().Where("user_id=? && video_id=?", user.Id, videoId).First(&favour)
	if favour.Id != "" {
		if action == "1" {
			// 设置faveourite
			err := db.Table("video").Debug().Where("video_id=? AND user_id=?", videoId, user.Id).Find(&video).Error
			if err != nil {
				c.JSON(200, Response{
					-1,
					fmt.Sprintf("occur err:%s", err),
				})
				return
			}
			favourNum := video.FavouriteCount + 1
			db.Table("video").Debug().Where("video_id=?", videoId).Update("favourite_count", favourNum)
			favour = Favourite{
				Id:      utils.GenerateUUID(),
				UserId:  user.Id,
				VideoId: videoId,
			}
			db.Table("favourite").Debug().Create(&favour)
			c.JSON(200, Response{
				0,
				"收藏成功",
			})
		} else if action == "2" {
			// 取消faveourite
			res := db.Table("video").Debug().Where("video_id=?", videoId).Find(&video)
			if res.Error != nil {
				c.JSON(200, Response{
					-1,
					fmt.Sprintf("occur err:%s", res.Error.Error()),
				})
				return
			}
			num := video.FavouriteCount - 1
			models.DB.Table("video").Where("video_id=?", videoId).Update("favourite_count", num)
			models.DB.Table("favourite").Where("user_id=? AND video_id=?", user.Id, videoId).Delete(&favour)
			c.JSON(200, Response{
				0,
				"取消收藏成功",
			})
		}
	} else {
		c.JSON(200, Response{
			-1,
			"没有收藏信息",
		})
	}
}

func FavouriteList(c *gin.Context) {
	token := c.Query("token")
	user, err := utils.AnalyzeToke(token)
	if err != nil {
		c.JSON(200, Response{
			-1,
			fmt.Sprintf("occur err:%s", err.Error()),
		})
	}

	var videoIdList = make([]string, 10)
	db.Select("video_id").Debug().Table("favourite").Where("user_id=?", user.Id).Find(&videoIdList)
	videos := make([]VideoRawData, len(videoIdList))
	for i := 0; i < len(videoIdList); i++ {
		models.DB.Table("video").Debug().Select("video_id", "user_id", "title", "play_url", "cover_url", "favourite_count",
			"comment_count").Where("video_id=?", videoIdList[i]).Find(&videos[i])
	}
	var authorIdList = make([]string, 10)
	models.DB.Table("video").Debug().Select("user_id").Find(&authorIdList, videoIdList)
	var videoList = make([]Video, len(videos))
	for i := 0; i < len(videos); i++ {
		videoList[i].Id = videos[i].VideoId
		videoList[i].FavouriteCount = videos[i].FavouriteCount
		videoList[i].PlayUrl = videos[i].PlayUrl
		videoList[i].CoverUrl = videos[i].CoverUrl
		videoList[i].FavouriteCount = videos[i].FavouriteCount
		videoList[i].CommentCount = videos[i].CommentCount
		models.DB.Table("user").Select("user_id", "name", "follow_count", "follower_count").Find(&videoList[i].Author, authorIdList)
	}
	c.JSON(200, VideoListResponse{
		Response: Response{
			StatusCode: 0,
			StatusMsg:  "get favourites success",
		},
		VideoList: videoList,
	})
}

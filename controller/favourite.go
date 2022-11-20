package controller

import (
	"aweme_kitex/models"
	"aweme_kitex/utils"
	"fmt"

	"github.com/gin-gonic/gin"
)

func FavouriteAction(c *gin.Context) {
	userIdentity := c.Query("identity")
	// todo: use midware to auth token
	token := c.Query("token")
	_, _ = utils.AnalyzeToke(token)

	videoId := c.Query("videoId")
	userId := userIdentity
	action := c.Query("actionType")
	favour := models.Favourite{}
	video := models.Video{}
	models.DB.Table("favourite").Debug().Where("user_id=? && video_id=?", userId, videoId).First(&favour)
	if favour.Identity != "" {
		if action == "1" {
			// 设置faveourite
			res := models.DB.Table("video").Debug().Where("video_id=? AND user_id=?", videoId, userId).Find(&video)
			if res.Error != nil {
				c.JSON(200, Response{
					-1,
					fmt.Sprintf("occur err:%s", res.Error.Error()),
				})
				return
			}
			favourNum := video.FavouriteCount + 1
			models.DB.Table("video").Debug().Where("video_id=?", videoId).Update("favourite_count", favourNum)
			favour = models.Favourite{
				Identity: utils.GenerateUUID(),
				UserId:   userId,
				VideoId:  videoId,
			}
			models.DB.Table("favourite").Debug().Create(&favour)
			c.JSON(200, Response{
				0,
				"收藏成功",
			})
		} else if action == "2" {
			// 取消faveourite
			res := models.DB.Table("video").Debug().Where("video_id=?", videoId).Find(&video)
			if res.Error != nil {
				c.JSON(200, Response{
					-1,
					fmt.Sprintf("occur err:%s", res.Error.Error()),
				})
				return
			}
			num := video.FavouriteCount - 1
			models.DB.Table("video").Where("video_id=?", videoId).Update("favourite_count", num)
			models.DB.Table("favourite").Where("user_id=? AND video_id=?", userId, videoId).Delete(&favour)
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
	userIdentity := c.Query("identity")
	// todo: use token
	_ = c.Query("token")

	var videoIdList = make([]string, 10)
	models.DB.Select("video_id").Debug().Table("favourite").Where("user_id=?", userIdentity).Find(&videoIdList)
	videos := make([]models.Video, len(videoIdList))
	for i := 0; i < len(videoIdList); i++ {
		models.DB.Table("video").Debug().Select("video_id", "user_id", "title", "play_url", "cover_url", "favourite_count",
			"comment_count").Where("video_id=?", videoIdList[i]).Find(&videos[i])
	}
	var authorIdList = make([]string, 10)
	models.DB.Table("video").Debug().Select("user_id").Find(&authorIdList, videoIdList)
	var videoList = make([]Video, len(videos))
	for i := 0; i < len(videos); i++ {
		videoList[i].Id = videos[i].Id
		videoList[i].FavouriteCount = videos[i].FavouriteCount
		videoList[i].PlayUrl = videos[i].PlayUrl
		videoList[i].CoverUrl = videos[i].CoverUrl
		videoList[i].FavouriteCount = videos[i].FavouriteCount
		videoList[i].CommentCount = videos[i].CommentCount
		videoList[i].Author = videos[i].Author
	}
	c.JSON(200, VideoListResponse{
		Response: Response{
			StatusCode: 0,
			StatusMsg:  "get favourites success",
		},
		VideoList: videoList,
	})
}

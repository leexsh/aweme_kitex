package controller

import (
	"aweme_kitex/models"
	"aweme_kitex/utils"
	"context"
	"errors"
	"fmt"
	"net/http"
	"path/filepath"
	"sync"

	"github.com/gin-gonic/gin"
)

/*
发布作品
	1.视频上传本地./public
	2.视频上传COS
	3.视频信息写入mysql
*/

func Publish(c *gin.Context) {
	token := c.Query("token")
	user, err := CheckToken(token)

	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"message": "User doesn't exist",
		})
		return
	}

	data, err := c.FormFile("data")
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusMsg:  err.Error(),
			StatusCode: -1,
		})
		return
	}
	fileName := filepath.Base(data.Filename)
	finalName := fmt.Sprintf("%s_%s", user.Name, fileName)

	saveFile := filepath.Join("./public/", finalName)
	// 1.save local
	if err := c.SaveUploadedFile(data, saveFile); err != nil {
		c.JSON(200, Response{
			-1,
			err.Error(),
		})
		return
	}

	// 2.upload COS
	cosKey := fmt.Sprintf("%s/%s", user.Name, finalName)
	_, _, err = models.COSClient.Object.Upload(
		context.Background(), cosKey, saveFile, nil,
	)
	if err != nil {
		c.JSON(200, Response{
			-1,
			err.Error(),
		})
		return
	}
	// 3.wirte to mysql
	ourl := cos.Object.GetObjectURL(cosKey)

	title := c.Query("title")
	video := VideoRawData{
		VideoId: utils.GenerateUUID(),
		UserId:  user.Id,
		Title:   title,
		PlayUrl: ourl.String(),
	}
	if err := db.Table("video").Debug().Create(&video).Error; err != nil {
		c.JSON(200, Response{
			-1,
			err.Error(),
		})
		return
	}
	c.JSON(200, Response{
		0,
		fmt.Sprintf("%s uploaded successfully", title),
	})
}

// --------------------
type QueryUserVideoList struct {
	Token     string
	UserName  string
	VideoList []Video

	CurrentId    string
	VideoData    []*models.VideoRawData
	UserMap      map[string]*models.UserRawData
	FavouriteMap map[string]*models.FavouriteRaw
	RelationMap  map[string]*models.RelationRaw
}

func NewQueryUserVideoList(token string) *QueryUserVideoList {
	return &QueryUserVideoList{Token: token}
}

func (f *QueryUserVideoList) checkToken() error {
	user, err := CheckToken(f.Token)
	if err != nil {
		return err
	}
	f.CurrentId = user.Id
	f.UserName = user.Id
	return nil
}

func (f *QueryUserVideoList) prepareVideoInfo() error {
	videoData, err := models.NewVideoDaoInstance().QueryVideosByUserId(f.CurrentId)
	if err != nil {
		return err
	}
	f.VideoData = videoData

	videoIds := make([]string, 0)
	userIds := []string{f.CurrentId}
	for _, video := range f.VideoData {
		videoIds = append(videoIds, video.VideoId)
	}

	userMap, err := models.NewUserDaoInstance().QueryUserByIds(userIds)
	if err != nil {
		return err
	}
	f.UserMap = userMap

	var wg sync.WaitGroup
	wg.Add(2)
	var favoriteErr, relationErr error
	// 获取点赞信息
	go func() {
		defer wg.Done()
		favoriteMap, err := models.NewFavouriteDaoInstance().QueryFavoursByIds(f.CurrentId, videoIds)
		if err != nil {
			favoriteErr = err
			return
		}
		f.FavouriteMap = favoriteMap
	}()
	// 获取关注信息
	go func() {
		defer wg.Done()
		relationMap, err := models.NewRelationDaoInstance().QueryRelationByIds(f.CurrentId, userIds)
		if err != nil {
			relationErr = err
			return
		}
		f.RelationMap = relationMap
	}()
	wg.Wait()
	if favoriteErr != nil {
		return favoriteErr
	}
	if relationErr != nil {
		return relationErr
	}
	return nil
}

func (f *QueryUserVideoList) packVideoInfo() error {
	videoList := make([]Video, 0)
	for _, video := range f.VideoData {
		videoUser, ok := f.UserMap[video.UserId]
		if !ok {
			return errors.New("has no video user info for " + fmt.Sprint(video.UserId))
		}

		var isFavorite bool = false
		var isFollow bool = false

		_, ok = f.FavouriteMap[video.VideoId]
		if ok {
			isFavorite = true
		}
		_, ok = f.RelationMap[video.UserId]
		if ok {
			isFollow = true
		}
		videoList = append(videoList, Video{
			Id: video.VideoId,
			Author: User{
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

	f.VideoList = videoList
	return nil
}

func (f *QueryUserVideoList) do() ([]Video, error) {
	if err := f.checkToken(); err != nil {
		return nil, err
	}
	if err := f.prepareVideoInfo(); err != nil {
		return nil, err
	}
	if err := f.packVideoInfo(); err != nil {
		return nil, err
	}
	return f.VideoList, nil
}

func QueryUserVideos(token string) ([]Video, error) {
	return NewQueryUserVideoList(token).do()
}

func QueryVodeoList(token string) VideoListResponse {
	if token == "" {
	}
	videos, err := QueryUserVideos(token)
	if err != nil {
		return VideoListResponse{
			Response{
				-1,
				err.Error(),
			},
			nil,
		}
	}
	return VideoListResponse{
		Response{
			0, "success",
		},
		videos,
	}
}

func PublishList(c *gin.Context) {
	token := c.Query("token")
	videoRes := QueryVodeoList(token)
	c.JSON(200, videoRes)
}

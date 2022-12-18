package controller

import (
	"aweme_kitex/models"
	"aweme_kitex/utils"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
)

/*
direct get data from database
*/

// ------------service--------------------
// 该层负责鉴权  向repository获取视频数据和封装数据

func QueryVideoData(latestTime int64, token string) ([]Video, int64, error) {
	return NewQueryVideoDataFlow(latestTime, token).Do()
}

func NewQueryVideoDataFlow(latestTime int64, token string) *QueryVideoDataFlow {
	return &QueryVideoDataFlow{
		LatestTime: latestTime,
		Token:      token,
	}
}

// video data
type QueryVideoDataFlow struct {
	LatestTime int64
	Token      string
	VideoList  []Video
	NextTime   int64

	CurrentId   string
	VideoData   []*models.VideoRawData
	UserMap     map[string]*models.UserRawData
	FavoursMap  map[string]*models.FavouriteRaw
	RelationMap map[string]*models.RelationRaw
}

func (f *QueryVideoDataFlow) Do() ([]Video, int64, error) {
	user, err := CheckToken(f.Token)
	if err != nil {
		return nil, 0, err
	}
	f.CurrentId = user.Id

	if err := f.prepareVideoInfo(); err != nil {
		return nil, 0, err
	}
	if err := f.packVideoInfo(); err != nil {
		return nil, 0, err
	}
	return f.VideoList, f.NextTime, nil
}

// prepare video
func (f *QueryVideoDataFlow) prepareVideoInfo() error {
	// 1.get video
	videoData, err := models.NewVideoDaoInstance().QueryVideoByLatestTime(f.LatestTime)
	if err != nil {
		return err
	}
	f.VideoData = videoData

	// 2. get video_id and user_id
	videoIds := make([]string, 0)
	authorIds := make([]string, 0)
	for _, video := range f.VideoData {
		videoIds = append(videoIds, video.VideoId)
		authorIds = append(authorIds, video.UserId)
	}

	// 3. get user info
	usersMap, err := models.NewUserDaoInstance().QueryUserByIds(authorIds)
	if err != nil {
		return err
	}
	f.UserMap = usersMap

	// 4. should login
	if f.CurrentId == "" {
		return nil
	}

	var wg sync.WaitGroup
	wg.Add(2)
	var favourErr, relationErr error
	// 5.获取点赞信息
	go func() {
		defer wg.Done()
		favoursMap, err := models.NewFavouriteDaoInstance().QueryFavoursByIds(f.CurrentId, videoIds)
		if err != nil {
			favourErr = err
			return
		}
		f.FavoursMap = favoursMap
	}()
	// 6.获取关注信息
	go func() {
		defer wg.Done()
		relationMap, err := models.NewRelationDaoInstance().QueryRelationByIds(f.CurrentId, authorIds)
		if err != nil {
			relationErr = err
			return
		}
		f.RelationMap = relationMap

	}()
	wg.Wait()

	if favourErr != nil {
		return favourErr
	}
	if relationErr != nil {
		return relationErr
	}
	return nil
}

func (f *QueryVideoDataFlow) packVideoInfo() error {
	videoList := make([]Video, 0)
	for _, video := range f.VideoData {
		videoAuthor, ok := f.UserMap[video.UserId]
		if !ok {
			return errors.New("has no video user info for " + fmt.Sprint(video.UserId))
		}
		var isFavourite bool = false
		var isFollow bool = false
		if f.CurrentId != "" {
			if _, ok := f.FavoursMap[video.VideoId]; ok {
				isFavourite = true
			}
			if _, ok := f.RelationMap[video.UserId]; ok {
				isFollow = true
			}
		}
		videoList = append(videoList, Video{
			Id: video.VideoId,
			Author: User{
				UserId:        videoAuthor.UserId,
				Name:          videoAuthor.Name,
				FollowCount:   videoAuthor.FollowCount,
				FollowerCount: videoAuthor.FollowerCount,
				IsFollow:      isFollow,
			},
			PlayUrl:        video.PlayUrl,
			CoverUrl:       video.CoverUrl,
			Title:          video.Title,
			FavouriteCount: video.FavouriteCount,
			CommentCount:   video.CommentCount,
			IsFavourite:    isFavourite,
		})
	}
	f.VideoList = videoList
	if len(f.VideoList) == 0 {
		f.NextTime = 0
	} else {
		f.NextTime = f.VideoData[len(f.VideoData)-1].CreatedAt.Unix()
	}
	return nil
}

// --------- handler ---------------
// 该层功能包括处理传入参数，向service层获取视频信息，封装成响应信息
type FeedResponse struct {
	Response
	VideoList []Video `json:"video_list,omitempty"`
	NextTime  int64   `json:"next_time,omitempty"`
}

func QueryVideoFeed(token string, latestTimeStr string) FeedResponse {
	// 1.处理传人参数
	latestTime, err := strconv.ParseInt(latestTimeStr, 10, 64)
	if err != nil {
		return FeedResponse{
			Response: Response{
				StatusCode: -1,
				StatusMsg:  err.Error(),
			},
		}
	}

	// 2.ge video
	videoList, nextTime, err := QueryVideoData(latestTime, token)
	if err != nil {
		return FeedResponse{
			Response: Response{
				StatusCode: -1,
				StatusMsg:  err.Error(),
			},
		}
	}
	return FeedResponse{
		Response{
			StatusCode: 0,
			StatusMsg:  "获取video成功",
		},
		videoList,
		nextTime,
	}
}

// ----- controller------
// Feed same demo video list for every request
// 该层功能包括获取传入参数，向handler获取视频信息，返回响应信息
func Feed(c *gin.Context) {
	token := c.DefaultQuery("token", defaultToken)

	defaultTimeStr := strconv.Itoa(int(utils.GetUnix()))

	feedRes := QueryVideoFeed(token, defaultTimeStr)
	c.JSON(http.StatusOK, feedRes)
}

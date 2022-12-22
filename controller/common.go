package controller

import (
	"aweme_kitex/cfg"
	"aweme_kitex/model"
	"aweme_kitex/utils"
	"errors"

	"github.com/gin-gonic/gin"
)

type (
	// videoRawData model.VideoRawData
	// favouriteRaw model.FavouriteRaw
	// userRawData  model.UserRawData
	// relationRaw  model.RelationRaw
	// commetRaw    model.CommentRaw
	//

	video     model.Video
	favourite model.Favourite
	comment   model.Comment
	user      model.User

	// response             model.Response
	// userLoginResponse    model.UserLoginResponse
	// userRegisterResponse model.UserRegisterResponse
	// userResponse         model.UserResponse
	// userListResponse     model.UserListResponse
	// videoListResponse    model.VideoListResponse
	// commentListResponse  model.CommentListResponse
	// feedResponse         model.FeedResponse
)

var (
	db = cfg.DB

	defaultToken = "defaultToken"

	usersLoginInfo = map[string]model.User{
		"caiXuKun": {
			UserId:        "asdd",
			Name:          "caiXuKun",
			FollowerCount: 5,
			FollowCount:   20,
			IsFollow:      true,
		},
	}

	userIdSequeue = int64(1)

	u = model.UserRawData{}
)

// 鉴权
func CheckToken(token string) (*utils.UserClaim, error) {
	if token == defaultToken {
		return nil, errors.New("error: check token failed, please update Token")
	}
	uc, err := utils.AnalyzeToke(token)
	if err != nil {
		return nil, err
	}
	return uc, nil
}

func TokenErrorRes(c *gin.Context, err error) {
	c.JSON(200, model.Response{
		-1,
		err.Error(),
	})
	return
}

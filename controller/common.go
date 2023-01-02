package controller

import (
	"aweme_kitex/models"
	"errors"

	"github.com/gin-gonic/gin"
)

type (
	// videoRawData models.VideoRawData
	// favouriteRaw models.FavouriteRaw
	// userRawData  models.UserRawData
	// relationRaw  models.RelationRaw
	// commetRaw    models.CommentRaw
	//

	video     models.Video
	favourite models.Favourite
	comment   models.Comment
	user      models.User

	// response             models.Response
	// userLoginResponse    models.UserLoginResponse
	// userRegisterResponse models.UserRegisterResponse
	// userResponse         models.UserResponse
	// userListResponse     models.UserListResponse
	// videoListResponse    models.VideoListResponse
	// commentListResponse  models.CommentListResponse
	// feedResponse         models.FeedResponse
)

var (
	defaultToken = "defaultToken"
)

// 鉴权
func CheckToken(token string) (*models.UserClaim, error) {
	if token == defaultToken {
		return nil, errors.New("error: check token failed, please update Token")
	}
	uc, err := models.AnalyzeToken(token)
	if err != nil {
		return nil, err
	}
	return uc, nil
}

func TokenErrorRes(c *gin.Context, err error) {
	c.JSON(200, models.Response{
		-1,
		err.Error(),
	})
	return
}

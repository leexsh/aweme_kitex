package handlers

import (
	"aweme_kitex/cmd/api/rpc"
	"aweme_kitex/cmd/favourite/kitex_gen/favourite"
	"aweme_kitex/pkg/errno"
	"aweme_kitex/pkg/jwt"
	"context"

	"github.com/gin-gonic/gin"
)

// FavoriteAction implement like and unlike operations
func FavoriteAction(c *gin.Context) {
	token := c.Query("token")
	_, err := jwt.AnalyzeToken(token)
	if err != nil {
		SendResponse(c, errno.TokenInvalidErr)
		return
	}
	videoIdStr := c.Query("videoId")
	actionTypeStr := c.Query("actionType")

	favouriteReq := &favourite.FavouriteActionRequest{
		Token:      token,
		VideoId:    videoIdStr,
		ActionType: actionTypeStr,
	}
	err = rpc.FavoriteAction(context.Background(), favouriteReq)
	if err != nil {
		SendResponse(c, errno.ConvertErr(err))
		return
	}

	SendResponse(c, errno.Success)
}

// FavoriteList get favorite list info
func FavoriteList(c *gin.Context) {
	token := c.Query("token")
	_, err := jwt.AnalyzeToken(token)
	if err != nil {
		SendResponse(c, errno.TokenInvalidErr)
		return
	}
	favouriteListReq := &favourite.FavouriteListRequest{Token: token}
	videoList, err := rpc.FavoriteList(context.Background(), favouriteListReq)
	if err != nil {
		SendResponse(c, errno.ConvertErr(err))
		return
	}
	SendFavoriteListResponse(c, errno.Success, videoList)
}

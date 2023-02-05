package handlers

import (
	"aweme_kitex/cmd/api/rpc"
	"aweme_kitex/cmd/favourite/kitex_gen/favourite"
	constants "aweme_kitex/pkg/constant"
	"aweme_kitex/pkg/errno"
	"context"

	"github.com/gin-gonic/gin"
)

// FavoriteAction implement like and unlike operations
func FavoriteAction(c *gin.Context) {
	token := c.Query("token")
	videoIdStr := c.Query("video_id")
	actionTypeStr := c.Query("action_type")

	favouriteReq := &favourite.FavouriteActionRequest{
		Token:      token,
		VideoId:    videoIdStr,
		ActionType: actionTypeStr,
	}
	err := rpc.FavoriteAction(context.Background(), favouriteReq)
	if err != nil {
		SendResponse(c, errno.ConvertErr(err), nil)
		return
	}

	SendResponse(c, errno.Success, nil)
}

// FavoriteList get favorite list info
func FavoriteList(c *gin.Context) {
	token := c.Query("token")

	favouriteListReq := &favourite.FavouriteListRequest{Token: token}
	videoList, err := rpc.FavoriteList(context.Background(), favouriteListReq)
	if err != nil {
		SendResponse(c, errno.ConvertErr(err), nil)
		return
	}

	SendResponse(c, errno.Success, map[string]interface{}{constants.VideoList: videoList})
}

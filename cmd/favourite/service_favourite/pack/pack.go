package favPack

import (
	"aweme_kitex/cmd/favourite/kitex_gen/favourite"
	"aweme_kitex/cmd/favourite/kitex_gen/feed"
	"time"
)

func FavouriteActionResponse(code int64, msg string) (resp *favourite.FavouriteActionResponse) {
	resp = new(favourite.FavouriteActionResponse)
	resp.BaseResp.StatusCode = code
	resp.BaseResp.StatusMsg = msg
	resp.BaseResp.ServiceTime = time.Now().Unix()

	return
}

func FavouriteListResponse(code int64, msg string, videos []*feed.Video) (resp *favourite.FavouriteListResponse) {
	resp = new(favourite.FavouriteListResponse)
	resp.BaseResp.StatusCode = code
	resp.BaseResp.StatusMsg = msg
	resp.BaseResp.ServiceTime = time.Now().Unix()
	if videos != nil {
		resp.VideoList = videos
	}
	return
}

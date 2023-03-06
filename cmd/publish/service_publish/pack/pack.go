package publishPack

import (
	"aweme_kitex/cmd/publish/kitex_gen/feed"
	"aweme_kitex/cmd/publish/kitex_gen/publish"
	"time"
)

func PackPublishAction(code int64, msg string) (resp *publish.PublishActionResponse) {
	resp = new(publish.PublishActionResponse)
	resp.BaseResp.StatusCode = code
	resp.BaseResp.StatusMsg = msg
	resp.BaseResp.ServiceTime = time.Now().Unix()
	return
}

func PackPublishList(code int64, msg string, videos []*feed.Video) (resp *publish.PublishListResponse) {
	resp = new(publish.PublishListResponse)
	resp.BaseResp.StatusCode = code
	resp.BaseResp.StatusMsg = msg
	resp.BaseResp.ServiceTime = time.Now().Unix()
	if len(videos) > 0 {
		resp.VideoList = videos
	}
	return
}

package handlers

import (
	"aweme_kitex/cmd/api/rpc"
	"aweme_kitex/cmd/publish/kitex_gen/publish"
	"aweme_kitex/pkg/errno"
	"aweme_kitex/pkg/jwt"
	"bytes"
	"context"
	"io"

	"github.com/gin-gonic/gin"
)

func Publish(c *gin.Context) {
	token := c.Query("token")
	_, err := jwt.AnalyzeToken(token)
	if err != nil {
		SendResponse(c, errno.TokenInvalidErr)
		return
	}
	title := c.Query("title")
	if length := len(title); length <= 0 || length > 128 {
		SendResponse(c, errno.ParamErr)
		return
	}

	data, _, err := c.Request.FormFile("data")
	if err != nil {
		SendResponse(c, errno.ConvertErr(err))
		return
	}
	defer data.Close()
	// 处理视频数据
	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, data); err != nil {
		SendResponse(c, errno.VideoDataCopyErr)
		return
	}
	video := buf.Bytes()

	req := &publish.PublishActionRequest{
		Token: token,
		Data:  video,
		Title: title,
	}
	err = rpc.PublishVideoData(context.Background(), req)
	if err != nil {
		SendResponse(c, errno.ConvertErr(err))
		return
	}
	SendResponse(c, errno.Success)
}

func PublishList(c *gin.Context) {
	token := c.Query("token")
	_, err := jwt.AnalyzeToken(token)
	if err != nil {
		SendResponse(c, errno.TokenInvalidErr)
		return
	}
	videoList, err := rpc.PublishVideoList(context.Background(), &publish.PublishListRequest{Token: token})
	if err != nil {
		SendResponse(c, errno.ConvertErr(err))
		return
	}
	SendPublishListResponse(c, errno.Success, videoList)
}

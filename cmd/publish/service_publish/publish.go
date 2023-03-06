package service_publish

import (
	"aweme_kitex/cmd/publish/kitex_gen/publish"
	"aweme_kitex/cmd/publish/service_publish/db"
	"aweme_kitex/models"
	"aweme_kitex/pkg/jwt"
	"aweme_kitex/pkg/utils"
	"bytes"
	"context"
	"fmt"
	"image/jpeg"
	"os"

	"github.com/disintegration/imaging"
	ffmpeg_go "github.com/u2takey/ffmpeg-go"
)

type PublishService struct {
	ctx context.Context
}

// NewPublishService new PublishService
func NewPublishService(ctx context.Context) *PublishService {
	return &PublishService{ctx: ctx}
}

// Publish upload video info
func (s *PublishService) Publish(req *publish.PublishActionRequest) error {
	uc, _ := jwt.AnalyzeToken(req.Token)
	video := req.Data
	title := req.Title

	fileName := fmt.Sprintf(uc.Id + title)
	filePath := "/public/" + fileName
	// 1.将视频保存到本地文件夹
	err := db.NewCOSDaoInstance().PublishBinaryDataToPublic(s.ctx, video, filePath)
	if err != nil {
		return err
	}
	// 2.上传oss
	cosKey := fileName
	err = db.NewCOSDaoInstance().PublishVideoToCOS(context.Background(), cosKey, filePath)
	// 2.upload cos
	if err != nil {
		return err
	}

	// 3.获取封面
	// coverName := fileName + ".jpg"
	// coverData, err := s.getSnapshot(filePath)
	// if err != nil {
	// 	return err
	// }
	// coverKey := "cover/" + coverName
	ourl := db.NewCOSDaoInstance().GetCOSVideoURL(cosKey)
	// 3.获取播放链接
	video1 := &models.VideoRawData{
		VideoId: utils.GenerateUUID(),
		UserId:  uc.Id,
		Title:   title,
		PlayUrl: ourl.String(),
	}
	err = db.NewCOSDaoInstance().SaveVideoData(context.Background(), video1)
	if err != nil {
		return err
	}
	return nil
}

// 缩略图
func (s *PublishService) getSnapshot(videoUrl string) ([]byte, error) {
	buffer := bytes.NewBuffer(nil)
	err := ffmpeg_go.Input(videoUrl).Filter("select", ffmpeg_go.Args{fmt.Sprintf("gte(n, %d)", 1)}).
		Output("pipe:", ffmpeg_go.KwArgs{"vframes": 1, "format": "image2", "vcodec": "mjpeg"}).
		WithOutput(buffer, os.Stdout).
		Run()
	if err != nil {
		return nil, err
	}
	img, err := imaging.Decode(buffer)
	if err != nil {
		return nil, err
	}
	buf := new(bytes.Buffer)
	jpeg.Encode(buf, img, nil)
	return buf.Bytes(), nil
}

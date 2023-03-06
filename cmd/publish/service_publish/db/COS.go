package db

import (
	"aweme_kitex/cfg"
	"aweme_kitex/models"
	"aweme_kitex/pkg/logger"
	"context"
	"mime/multipart"
	"net/url"
	"os"
	"sync"

	"github.com/gin-gonic/gin"
)

type COSDao struct {
}

var (
	cosDao  *COSDao
	cosOnce sync.Once
)

func NewCOSDaoInstance() *COSDao {
	cosOnce.Do(
		func() {
			cosDao = &COSDao{}
		})
	return cosDao
}

// 将本地文件夹中的视频上传到COS
func (*COSDao) PublishVideoToCOS(ctx context.Context, cosKey string, saveFile string) error {
	_, _, err := cfg.COSClient.Object.Upload(
		ctx, cosKey, saveFile, nil,
	)
	if err != nil {
		return err
	}
	return nil
}

// 获取key的COS URL
func (*COSDao) GetCOSVideoURL(cosKey string) *url.URL {
	return cfg.COSClient.Object.GetObjectURL(cosKey)
}

// 上传封面到Oss

// 从oss上获取封面地址

// 将视频保存到本地文件夹中
func (*COSDao) PublishVideoToPublic(ctx context.Context, video *multipart.FileHeader, path string, c *gin.Context) error {
	if err := c.SaveUploadedFile(video, path); err != nil {
		logger.Error("save videos to local error : " + err.Error())
		return err
	}
	return nil
}

func (*COSDao) PublishBinaryDataToPublic(ctx context.Context, video []byte, filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		logger.Error("create %s faile, %v", filePath, err.Error())
		return err
	}
	defer file.Close()
	_, err = file.Write(video)
	if err != nil {
		logger.Error("write file err, %v", err.Error())
		return err
	}
	return nil
}

func (*COSDao) SaveVideoData(ctx context.Context, videoData *models.VideoRawData) error {
	err := cfg.DB.WithContext(ctx).Table("video").Debug().Create(videoData).Error
	if err != nil {
		logger.Error("create video error : " + err.Error())
		return err
	}
	return nil
}

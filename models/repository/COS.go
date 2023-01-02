package repository

import (
	"aweme_kitex/cfg"
	"context"
	"net/url"
	"sync"
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
func (*COSDao) PublishVideoToCOS(cosKey string, saveFile string) error {
	_, _, err := cfg.COSClient.Object.Upload(
		context.Background(), cosKey, saveFile, nil,
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

package main

import (
	"aweme_kitex/cfg"
	"aweme_kitex/pkg/logger"
	"aweme_kitex/routers"
)

func main() {
	if err := Init(); err != nil {
		panic(err)
	}
	defer logger.Sync() // 将缓存同步到文件中
	logger.Info("gin is running.....")
	routers.Include(routers.Video, routers.Follow, routers.User, routers.Comment, routers.Favourite, routers.Feed)
	r := routers.Init()
	if err := r.Run(); err != nil {
		logger.Error("startup service_user failed, err" + err.Error())
	}

	return
}

func Init() error {
	if err := logger.InitLogger(); err != nil {
		logger.Error("log init err: " + err.Error())
		return err
	}
	if err := cfg.Init(); err != nil {
		logger.Error("cfg init err: " + err.Error())
		return err
	}
	return nil
}

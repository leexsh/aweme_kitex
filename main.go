package main

import (
	"aweme_kitex/cfg"
	"aweme_kitex/routers"
	"aweme_kitex/utils"
)

func main() {
	if err := Init(); err != nil {
		panic(err)
	}
	defer utils.Sync() // 将缓存同步到文件中
	utils.Info("gin is running.....")
	routers.Include(routers.Video, routers.Follow, routers.User, routers.Comment, routers.Favourite, routers.Feed)
	r := routers.Init()
	if err := r.Run(); err != nil {
		utils.Error("startup service failed, err" + err.Error())
	}

	return
}

func Init() error {
	if err := utils.InitLogger(); err != nil {
		utils.Error("log init err: " + err.Error())
		return err
	}
	if err := cfg.Init(); err != nil {
		utils.Error("cfg init err: " + err.Error())
		return err
	}
	return nil
}

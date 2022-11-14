package main

import (
	"aweme_kitex/routers"
	"fmt"
)

func main() {
	routers.Include(routers.Video, routers.Follow, routers.User, routers.Comment, routers.Favourite)
	r := routers.Init()
	if err := r.Run(); err != nil {
		fmt.Printf("startup service failed, err is %#\n", err)
	}
	return
}

package main

import (
	feed "aweme_kitex/cmd/feed/kitex_gen/feed/feedservice"
	"log"
)

func main() {

	svr := feed.NewServer(new(FeedServiceImpl))

	err := svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}

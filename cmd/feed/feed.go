package main

import (
	"aweme_kitex/cmd/feed/kitex_gen/feed"
	"context"
)

type FeedService struct {
	ctx context.Context
}

func NewFeedService(ctx context.Context) *FeedService {
	return &FeedService{ctx: ctx}
}
func (f *FeedService) Feed(req *feed.FeedRequest) ([]*feed.Video, int64, error) {
	// todo
	// currentId, err := f.checkToken(req.Token)
	// if err != nil {
	// 	return nil, 0, err
	// }
	//
	// //get video info
	// videoData, err := db.QueryVideoByLatestTime(s.ctx, req.LatestTime)
	// if err != nil {
	// 	return  nil, 0, err
	// }
	//
	// //get video ids and user ids
	// videoIds, userIds := pack.Ids(videoData)
	//
	// //get user info
	// users, err :=
	return nil, 0, nil
}

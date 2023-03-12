package service_feed

import (
	"aweme_kitex/cmd/feed/kitex_gen/feed"
	"aweme_kitex/cmd/feed/kitex_gen/user"
	feedRPC "aweme_kitex/cmd/feed/rpc"
	videoDB "aweme_kitex/cmd/feed/service_feed/db"
	user2 "aweme_kitex/cmd/user/kitex_gen/user"
	"aweme_kitex/pkg/types"
	"context"
)

type CheckVideoService struct {
	ctx context.Context
}

// NewCheckVideoService new CheckVideoService
func NewCheckVideoService(ctx context.Context) *CheckVideoService {
	return &CheckVideoService{
		ctx: ctx,
	}
}

func (c *CheckVideoService) CheckVideoInvalid(vids []string) error {
	_, err := videoDB.NewVideoDaoInstance().QueryVideosByIs(c.ctx, vids)
	return err
}

func (c *CheckVideoService) GetVideos(vids []string) ([]*feed.Video, error) {
	vidos, err := videoDB.NewVideoDaoInstance().QueryVideosByIs(c.ctx, vids)
	if err != nil {
		return nil, err
	}
	uids := make([]string, len(vidos))
	for i := 0; i < len(vidos); i++ {
		uids = append(uids, vidos[i].UserId)
	}
	users, err := feedRPC.GetUserInfo(c.ctx, &user2.SingleUserInfoRequest{UserIds: uids})
	return c.packageinfo(vidos, users)
}

func (c *CheckVideoService) GetVideosByUserID(userId string) ([]*feed.Video, error) {
	vidos, err := videoDB.NewVideoDaoInstance().QueryVideosByUserId(c.ctx, userId)
	if err != nil {
		return nil, err
	}
	users, err := feedRPC.GetUserInfo(c.ctx, &user2.SingleUserInfoRequest{UserIds: []string{userId}})

	return c.packageinfo(vidos, users)
}

func (c *CheckVideoService) packageinfo(videos []*types.VideoRawData, users map[string]*user2.User) ([]*feed.Video, error) {
	videoList := make([]*feed.Video, len(videos))
	for i := 0; i < len(videos); i++ {
		v := videos[i]
		u := users[v.UserId]
		videoList = append(videoList, &feed.Video{
			VideoId: v.VideoId,
			Author: &user.User{
				UserId:        u.UserId,
				Name:          u.Name,
				FollowCount:   u.FollowCount,
				FollowerCount: u.FollowerCount,
			},
			PlayUrl:        v.PlayUrl,
			CoverUrl:       v.CoverUrl,
			FavouriteCount: v.FavouriteCount,
			CommentCount:   v.CommentCount,
			IsFavourite:    false,
			Title:          v.Title,
		})
	}
	return videoList, nil
}

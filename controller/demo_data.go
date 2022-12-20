package controller

import "aweme_kitex/model"

var DemoVideos = []model.Video{
	{
		Id:             "asd",
		Author:         DemoUser,
		PlayUrl:        "https://www.bilibili.com/video/BV1Ve4y147D2?t=4.7",
		CoverUrl:       "https://c-ssl.duitang.com/uploads/item/202006/13/20200613202923_flfxg.jpg",
		FavouriteCount: 2,
		CommentCount:   0,
		IsFavourite:    false,
	},
}

var DemoComments = []model.Comment{
	{
		Id:      "aaa",
		UserId:  "asd",
		VideoId: "qasd",
		Content: "test content",
	},
}

var DemoUser = model.User{
	UserId:        "XXX",
	Name:          "John",
	FollowCount:   20,
	FollowerCount: 12,
	IsFollow:      false,
}

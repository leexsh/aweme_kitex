package controller

import (
	"aweme_kitex/models"
)

var DemoVideos = []Video{
	{
		Id:             1,
		Author:         models.User(DemoUser),
		PlayUrl:        "https://www.bilibili.com/video/BV1Ve4y147D2?t=4.7",
		CoverUrl:       "https://c-ssl.duitang.com/uploads/item/202006/13/20200613202923_flfxg.jpg",
		FavouriteCount: 2,
		CommentCount:   0,
		IsFavourite:    false,
	},
}

var DemoComments = []Comment{
	{
		Id:         1,
		User:       models.User(DemoUser),
		Content:    "test content",
		CreateDate: "11-11",
	},
}

var DemoUser = User{
	Id:            1,
	Name:          "John",
	FollowCount:   20,
	FollowerCount: 12,
	IsFollow:      false,
}

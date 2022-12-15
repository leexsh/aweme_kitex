package controller

var DemoVideos = []Video{
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

var DemoComments = []Comment{
	{
		Id:      "aaa",
		UserId:  "asd",
		VideoId: "qasd",
		Content: "test content",
	},
}

var DemoUser = User{
	UserId:        "XXX",
	Name:          "John",
	FollowCount:   20,
	FollowerCount: 12,
	IsFollow:      false,
}

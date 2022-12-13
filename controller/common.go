package controller

import "aweme_kitex/models"

var (
	db = models.DB

	defaultToken = "defaultToken"

	address = "http://localhost:8080/aweme/"

	usersLoginInfo = map[string]User{
		"caiXuKun": {
			UserId:        "asdd",
			Name:          "caiXuKun",
			FollowerCount: 5,
			FollowCount:   20,
			IsFollow:      true,
		},
	}

	userIdSequeue = int64(1)

	u = UserRawData{}
)

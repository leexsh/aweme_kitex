package repository

import "aweme_kitex/models"

// 缓存的用户信息表，存储token到用户的映射
// 该缓存数据在服务重新启动自动清除
var UsersLoginInfo = map[string]*models.UserRawData{
	"JerryJerry123": {
		UserId:        "JerryJerry123",
		Name:          "Jerry",
		Password:      "Jerry123",
		FollowCount:   0,
		FollowerCount: 0,
	},
}

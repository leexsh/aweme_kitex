package controller

import (
	"aweme_kitex/models"
	"aweme_kitex/utils"
	"errors"
)

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

	u = models.UserRawData{}
)

// 鉴权
func CheckToken(token string) (*utils.UserClaim, error) {
	if token == defaultToken {
		return nil, errors.New("error: check token failed, please update Token")
	}
	uc, err := utils.AnalyzeToke(token)
	if err != nil {
		return nil, err
	}
	return uc, nil
}

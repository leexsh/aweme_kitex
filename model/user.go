package model

import (
	"aweme_kitex/cfg"
	"errors"
	"sync"
)

// ---------------------user----------------------------

type UserDao struct {
}

var (
	userDao  *UserDao
	userOnce sync.Once
)

func NewUserDaoInstance() *UserDao {
	userOnce.Do(
		func() {
			userDao = &UserDao{}
		})
	return userDao
}

func (u2 *UserDao) QueryUserByIds(uIds []string) (map[string]*UserRawData, error) {
	var users []*UserRawData
	err := cfg.DB.Debug().Where("user_id in (?)", uIds).Find(&users).Error
	if err != nil {
		return nil, errors.New("query users fail")
	}
	userMap := make(map[string]*UserRawData)
	for _, user := range users {
		userMap[user.UserId] = user
	}
	return userMap, nil
}

package model

import (
	"aweme_kitex/cfg"
	"errors"
	"sync"

	"gorm.io/gorm"
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

// 根据用户id获取用户信息
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

// 检查用户是否不存在
func (*UserDao) CheckUserNotExist(userId string) error {
	if _, exist := usersLoginInfo[userId]; exist {
		return errors.New("user already exists")
	}
	var user *UserRawData
	err := db.Table("user").Where("userId = ?", userId).First(&user).Error
	if err == nil {
		return errors.New("user already exists")
	}
	if err == gorm.ErrRecordNotFound {
		return nil
	}
	return err
}

// 上传用户信息到缓存的用户信息表和数据库
func (*UserDao) UploadUserData(user *UserRawData) error {
	usersLoginInfo[user.UserId] = *user
	err := db.Table("user").Create(&user).Error
	if err != nil {
		return err
	}
	return nil
}

// 通过token获取用户id和用户
func (*UserDao) QueryUserByUserId(userId string) (*UserRawData, error) {
	if userInfo, exist := usersLoginInfo[userId]; exist {
		return &userInfo, nil
	}

	var user *UserRawData
	err := db.Table("user").Where("user_id = ?", userId).First(&user).Error
	if err != nil {
		return nil, err
	}
	return user, nil

}

func (*UserDao) QueryUserByPassword(userName, password string) (*UserRawData, error) {
	if userInfo, exist := usersLoginInfo[userName]; exist {
		return &userInfo, nil
	}
	var usre *UserRawData
	err := db.Table("user").Where("name=? AND password=?", userName, password).First(&usre).Error
	if err != nil {
		return nil, err
	}
	return usre, nil

}

package repository

import (
	"aweme_kitex/models"
	"aweme_kitex/utils"
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
func (u2 *UserDao) QueryUserByIds(uIds []string) ([]*models.UserRawData, error) {
	var users []*models.UserRawData
	err := DB.Debug().Where("user_id in (?)", uIds).Find(&users).Error
	if err != nil {
		utils.Error("query user by Ids err: " + err.Error())
		return nil, errors.New("query users fail")
	}
	return users, nil
}

// 检查用户是否不存在
func (*UserDao) CheckUserNotExist(userId string) error {
	if _, exist := UsersLoginInfo[userId]; exist {
		return errors.New("user already exists")
	}
	var user *models.UserRawData
	err := DB.Table("user").Where("userId = ?", userId).First(&user).Error
	if err == gorm.ErrRecordNotFound {
		return nil
	}
	if err == nil {
		utils.Errorf("check user not exist fail, err:%s", err.Error())
		return errors.New("user already exists")
	}

	return err
}

// 上传用户信息到缓存的用户信息表和数据库
func (*UserDao) UploadUserData(user *models.UserRawData) error {
	UsersLoginInfo[user.UserId] = *user
	err := DB.Table("user").Create(&user).Error
	if err != nil {
		return err
	}
	return nil
}

// 通过token获取用户id和用户
func (*UserDao) QueryUserByUserId(userId string) (*models.UserRawData, error) {
	if userInfo, exist := UsersLoginInfo[userId]; exist {
		return &userInfo, nil
	}

	var user *models.UserRawData
	err := DB.Table("user").Where("user_id = ?", userId).First(&user).Error
	if err != nil {
		utils.Error("query user by Id err: " + err.Error())
		return nil, err
	}
	return user, nil

}

func (*UserDao) QueryUserByPassword(userName, password string) (*models.UserRawData, error) {
	if userInfo, exist := UsersLoginInfo[userName]; exist {
		return &userInfo, nil
	}
	var usre *models.UserRawData
	err := DB.Table("user").Where("name=? AND password=?", userName, password).First(&usre).Error
	if err != nil {
		utils.Error("query user by password err: " + err.Error())
		return nil, err
	}
	return usre, nil

}

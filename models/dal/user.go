package dal

import (
	"aweme_kitex/cfg"
	"aweme_kitex/models"
	"aweme_kitex/pkg/logger"
	"context"
	"encoding/json"
	"errors"
	"sync"
	"time"

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
func (u2 *UserDao) QueryUserByIds(ctx context.Context, uIds []string) ([]*models.UserRawData, error) {
	var users []*models.UserRawData
	err := cfg.DB.WithContext(ctx).Debug().Where("user_id in (?)", uIds).Find(&users).Error
	if err != nil {
		logger.Error("query user by Ids err: " + err.Error())
		return nil, errors.New("query users fail")
	}
	return users, nil
}

// 检查用户是否不存在
func (*UserDao) CheckUserNotExist(ctx context.Context, userId string) error {
	userRedis := &models.UserRawData{}
	val, err := redisGet(userId)
	_ = json.Unmarshal([]byte(val), userRedis)
	if userRedis.UserId == userId {
		return nil
	}
	var user *models.UserRawData
	err = cfg.DB.Table("user").WithContext(ctx).Where("userId = ?", userId).First(&user).Error
	if err == gorm.ErrRecordNotFound {
		return nil
	}
	if err == nil {
		logger.Errorf("check user not exist fail, err:%s", err.Error())
		return errors.New("user already exists")
	}

	return err
}

// 上传用户信息到缓存的用户信息表和数据库
func (*UserDao) UploadUserData(ctx context.Context, user *models.UserRawData) error {
	err := redisSet(user.UserId, marshal(user), time.Hour)
	err = cfg.DB.Table("user").WithContext(ctx).Create(&user).Error
	if err != nil {
		return err
	}
	return nil
}

// 通过token获取用户id和用户
func (*UserDao) QueryUserByUserId(ctx context.Context, userId string) (*models.UserRawData, error) {
	// 1. 有缓存，先从缓存中取出来
	userRedis := &models.UserRawData{}
	val, err := redisGet(userId)
	_ = unmarshal(val, userRedis)
	if userRedis.UserId == userId {
		return userRedis, nil
	}

	var user *models.UserRawData
	// 2. 没有缓存，先写数据库
	err = cfg.DB.Table("user").WithContext(ctx).Where("user_id = ?", userId).First(&user).Error
	if err != nil {
		logger.Error("query user by Id err: " + err.Error())
		return nil, err
	}
	{
		// 3.写缓存
		err = redisSet(userId, marshal(user), time.Hour)
		if err != nil {
			return nil, err
		}
	}
	return user, nil

}

func (*UserDao) QueryUserByPassword(ctx context.Context, userName, password string) (*models.UserRawData, error) {
	var usre *models.UserRawData
	err := cfg.DB.Table("user").WithContext(ctx).Where("name=? AND password=?", userName, password).First(&usre).Error
	if err != nil {
		logger.Error("query user by password err: " + err.Error())
		return nil, err
	}
	return usre, nil

}

// 检查用户是否存在
func (u *UserDao) CheckUserId(ctx context.Context, uids []string) ([]*models.UserRawData, error) {
	users, err := u.QueryUserByIds(ctx, uids)
	if err != nil {
		return nil, err
	}
	if len(users) == 0 {
		return nil, errors.New("userId not exist")
	}
	return users, nil
}

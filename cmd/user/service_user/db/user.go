package userDB

import (
	"aweme_kitex/cfg"
	"aweme_kitex/pkg/logger"
	"aweme_kitex/pkg/types"
	"context"
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
func (u2 *UserDao) QueryUserByIds(ctx context.Context, uIds []string) ([]*types.UserRawData, error) {
	var users []*types.UserRawData
	err := cfg.DB.WithContext(ctx).Debug().Where("user_id in (?)", uIds).Find(&users).Error
	if err != nil {
		logger.Error("query user by Ids err: " + err.Error())
		return nil, errors.New("query users fail")
	}
	return users, nil
}

// 检查用户是否不存在
func (*UserDao) CheckUserNotExist(ctx context.Context, userId string) error {
	userRedis := &types.UserRawData{}
	if userRedis.UserId == userId {
		return nil
	}
	var user *types.UserRawData
	err := cfg.DB.Table("user").WithContext(ctx).Where("userId = ?", userId).First(&user).Error
	if err == gorm.ErrRecordNotFound {
		return nil
	}
	if err == nil {
		logger.Error("check user not exist fail, err:%s", err.Error())
		return errors.New("user already exists")
	}

	return err
}

// 上传用户信息到缓存的用户信息表和数据库
func (*UserDao) UploadUserData(ctx context.Context, user *types.UserRawData) error {
	err := cfg.DB.Table("user").WithContext(ctx).Create(&user).Error
	if err != nil {
		return err
	}
	return nil
}

// 通过token获取用户id和用户
func (*UserDao) QueryUserByUserId(ctx context.Context, userId string) (*types.UserRawData, error) {

	var user *types.UserRawData
	err := cfg.DB.Table("user").WithContext(ctx).Where("user_id = ?", userId).First(&user).Error
	if err != nil {
		logger.Error("query user by Id err: " + err.Error())
		return nil, err
	}
	return user, nil

}

func (*UserDao) QueryUserByPassword(ctx context.Context, userName, password string) (*types.UserRawData, error) {
	var usre *types.UserRawData
	err := cfg.DB.Table("user").WithContext(ctx).Where("name=? AND password=?", userName, password).First(&usre).Error
	if err != nil {
		logger.Error("query user by password err: " + err.Error())
		return nil, err
	}
	return usre, nil

}

// 检查用户是否存在
func (u *UserDao) CheckUserId(ctx context.Context, uids []string) ([]*types.UserRawData, error) {
	users, err := u.QueryUserByIds(ctx, uids)
	if err != nil {
		return nil, err
	}
	if len(users) == 0 {
		return nil, errors.New("userId not exist")
	}
	return users, nil
}

// 关注操作--增加用户的关注数&增加被关注用户的粉丝数
func (u *UserDao) IncreaseFollowCount(ctx context.Context, userID, toUserID string) error {
	// 事务操作,保持连贯,一个完整的关注操作
	err := cfg.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 更新follow_count
		res := tx.Model(&types.UserRawData{}).Where("id = ?", userID).Update("follow_count", gorm.Expr("follow_count + ?", 1))
		if res.Error != nil {
			logger.Info("服务器增加follow_count失败")
			return res.Error
		}
		// 更新 follower_count 字段
		res = tx.Model(&types.UserRawData{}).Where("id = ?", toUserID).Update("follower_count", gorm.Expr("follower_count + ?", 1))
		if res.Error != nil {
			logger.Info("服务器增加follower_count失败")
			return res.Error
		}
		return nil
	})
	if err != nil {
		logger.Info("关注操作的事务出现问题！")
		return err
	}
	return nil
}

// 取关操作--减少用户的关注数&减少被关注用户的粉丝数
func (u *UserDao) DecreaseFollowCount(ctx context.Context, userID, toUserID string) error {
	// 事务操作,保持连贯,一个完整的关注操作
	err := cfg.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 更新follow_count
		res := tx.Model(&types.UserRawData{}).Where("id = ?", userID).Update("follow_count", gorm.Expr("follow_count - ?", 1))
		if res.Error != nil {
			logger.Info("服务器减少follow_count失败")
			return res.Error
		}
		// 更新 follower_count 字段
		res = tx.Model(&types.UserRawData{}).Where("id = ?", toUserID).Update("follower_count", gorm.Expr("follower_count - ?", 1))
		if res.Error != nil {
			logger.Info("服务器减少follower_count失败")
			return res.Error
		}
		return nil
	})
	if err != nil {
		logger.Info("取关操作的事务出现问题！")
		return err
	}
	return nil
}

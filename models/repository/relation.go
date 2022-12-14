package repository

import (
	"aweme_kitex/cfg"
	"aweme_kitex/models"
	"aweme_kitex/pkg/logger"
	"aweme_kitex/utils"
	"context"
	"errors"
	"sync"

	"gorm.io/gorm"
)

// -------------用户关系-------------------

type RelationDao struct {
}

var (
	relationDao  *RelationDao
	relationOnce sync.Once
)

func NewRelationDaoInstance() *RelationDao {
	relationOnce.Do(
		func() {
			relationDao = &RelationDao{}
		})
	return relationDao
}

// 1vN 根据当前用户Id和视频作者的id获取关注信息
func (r *RelationDao) QueryRelationByIds(ctx context.Context, currentUid string, userIds []string) (map[string]*models.RelationRaw, error) {
	var relations []*models.RelationRaw
	err := cfg.DB.WithContext(ctx).Table("relation").Where("user_id=? AND to_user_id IN ?", currentUid, userIds).Find(&relations).Error

	if err != nil {
		logger.Error("query relation by Id err: " + err.Error())
		return nil, errors.New("query relation record fail")
	}
	relationMap := make(map[string]*models.RelationRaw)
	for _, relation := range relations {
		relationMap[relation.ToUserId] = relation
	}
	return relationMap, nil
}

func (r *RelationDao) CreateRelation(ctx context.Context, userId, toUserId string) error {
	relation := &models.RelationRaw{
		Id:       utils.GenerateUUID(),
		UserId:   userId,
		ToUserId: toUserId,
		Status:   1,
	}
	cfg.DB.Transaction(func(tx *gorm.DB) error {
		err := tx.Debug().WithContext(ctx).Table("user").Where("user_id=?", userId).Update("follow_count", gorm.Expr("follow_count + ?", 1)).Error
		if err != nil {
			logger.Error("create relation err: " + err.Error())
			return err
		}
		err = tx.Debug().Table("user").Where("user_id=?", toUserId).Update("follower_count", gorm.Expr("follower_count + ?", 1)).Error
		if err != nil {
			return err
		}
		err = tx.Debug().Table("relation").Create(relation).Error
		if err != nil {
			return err
		}
		return nil
	})
	return nil
}

func (r *RelationDao) DeleteRelation(ctx context.Context, userId, toUserId string) error {
	var relationRaw *models.RelationRaw
	cfg.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err := tx.Table("user").Where("user_id = ?", userId).Update("follow_count", gorm.Expr("follow_count - ?", 1)).Error
		if err != nil {
			logger.Error("delete relation by Id err: " + err.Error())
			return err
		}

		err = tx.Table("user").Where("user_id = ?", toUserId).Update("follower_count", gorm.Expr("follower_count - ?", 1)).Error
		if err != nil {
			return err
		}

		err = tx.Table("relation").Where("user_id = ? AND to_user_id = ?", userId, toUserId).Delete(&relationRaw).Error
		if err != nil {
			return err
		}
		return nil
	})
	return nil
}

// get follow by uid
func (r *RelationDao) QueryFollowByUid(ctx context.Context, uid string) ([]*models.RelationRaw, error) {
	var relations []*models.RelationRaw
	err := cfg.DB.WithContext(ctx).Table("relation").Where("user_id = ?", uid).Find(&relations).Error
	if err != nil {
		return nil, err
	}
	return relations, nil
}

// 通过用户id，查询该用户的粉丝， 返回两者之间的关注记录
func (*RelationDao) QueryFollowerById(ctx context.Context, userId string) ([]*models.RelationRaw, error) {
	var relations []*models.RelationRaw
	err := cfg.DB.WithContext(ctx).Table("relation").Where("to_user_id = ?", userId).Find(&relations).Error
	if err != nil {
		return nil, err
	}
	return relations, nil
}

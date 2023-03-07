package db

import (
	"aweme_kitex/cfg"
	"aweme_kitex/pkg/logger"
	"aweme_kitex/pkg/types"
	"aweme_kitex/pkg/utils"
	"context"
	"errors"
	"sync"
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
func (r *RelationDao) QueryRelationByIds(ctx context.Context, currentUid string, userIds []string) (map[string]*types.RelationRaw, error) {
	var relations []*types.RelationRaw
	err := cfg.DB.WithContext(ctx).Table("relation").Where("user_id=? AND to_user_id IN ?", currentUid, userIds).Find(&relations).Error

	if err != nil {
		logger.Error("query relation by Id err: " + err.Error())
		return nil, errors.New("query relation record fail")
	}
	relationMap := make(map[string]*types.RelationRaw)
	for _, relation := range relations {
		relationMap[relation.ToUserId] = relation
	}
	return relationMap, nil
}

func (r *RelationDao) CreateRelation(ctx context.Context, userId, toUserId string) error {
	relation := &types.RelationRaw{
		Id:       utils.GenerateUUID(),
		UserId:   userId,
		ToUserId: toUserId,
		Status:   1,
	}
	return cfg.DB.WithContext(ctx).Debug().Table("relation").Create(relation).Error
}

func (r *RelationDao) DeleteRelation(ctx context.Context, userId, toUserId string) error {
	var relationRaw *types.RelationRaw
	return cfg.DB.WithContext(ctx).Debug().Table("relation").Where("user_id = ? AND to_user_id = ?", userId, toUserId).Delete(&relationRaw).Error
}

// get follow by uid
func (r *RelationDao) QueryFollowByUid(ctx context.Context, uid string) ([]*types.RelationRaw, error) {
	var relations []*types.RelationRaw
	err := cfg.DB.WithContext(ctx).Table("relation").Where("user_id = ?", uid).Find(&relations).Error
	if err != nil {
		return nil, err
	}
	return relations, nil
}

// 通过用户id，查询该用户的粉丝， 返回两者之间的关注记录
func (*RelationDao) QueryFollowerById(ctx context.Context, userId string) ([]*types.RelationRaw, error) {
	var relations []*types.RelationRaw
	err := cfg.DB.WithContext(ctx).Table("relation").Where("to_user_id = ?", userId).Find(&relations).Error
	if err != nil {
		return nil, err
	}
	return relations, nil
}

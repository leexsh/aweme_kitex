package model

import (
	"aweme_kitex/cfg"
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
func (r *RelationDao) QueryRelationByIds(currentUid string, userIds []string) (map[string]*RelationRaw, error) {
	var relations []*RelationRaw
	err := cfg.DB.Where("user_id=? AND to_user_id IN ? AND status IN ?", currentUid, userIds, []int64{0, -1}).
		Or("user_id IN ? AND to_user_id = ? AND status = ?", userIds, currentUid, 1).Find(&relations).Error

	if err == gorm.ErrRecordNotFound {
		return nil, errors.New("relation record not found")
	}
	if err != nil {
		return nil, errors.New("query relation record fail")
	}
	relationMap := make(map[string]*RelationRaw)
	for _, relation := range relations {
		if relation.Status == 1 {
			relationMap[relation.UserId] = relation
		} else {
			relationMap[relation.ToUserId] = relation
		}
	}
	return relationMap, nil
}

// 1v1 get relationShip
func (r *RelationDao) QueryRelationByUid(uid, toUid string) (*RelationRaw, error) {
	var relation *RelationRaw
	err := cfg.DB.Debug().Where("(user_id=? AND to_user_id=?) OR (user_id=? AND to_user_id=?)", uid, toUid, toUid, uid).Find(relation).Error
	return relation, err
}

func (r *RelationDao) InsertRaw(relation *RelationRaw) error {
	err := cfg.DB.Debug().Create(relation).Error
	return err
}

func (r *RelationDao) UpdateRaw(relation *RelationRaw, status int64) error {
	err := cfg.DB.Debug().Model(relation).Update("status", relation.Status).Error
	return err
}

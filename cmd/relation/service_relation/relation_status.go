package service_relation

import (
	"aweme_kitex/cmd/relation/service_relation/db"
	constants "aweme_kitex/pkg/constant"
	"aweme_kitex/pkg/logger"
	"context"
	"errors"
)

type RelationStatusService struct {
	ctx context.Context
}

// NewRelationActionService new RelationActionService
func NewRelationStatusService(ctx context.Context) *RelationStatusService {
	return &RelationStatusService{ctx: ctx}
}

func (r *RelationStatusService) QueryRelationStatus(uid, toUid string) (bool, error) {
	// 1.从redis获取关系
	flag, err := db.RelationClient.SIsMember(r.ctx, uid, toUid).Result()
	if flag {
		db.RelationClient.Expire(r.ctx, uid, constants.RedisExpireTime)
		return true, err
	}
	// 2.redis不存在 查数据库
	relationMap, err := db.NewRelationDaoInstance().QueryFollowerById(r.ctx, uid)
	if err != nil {
		logger.Error(err)
		return false, err
	}
	for _, relationRaw := range relationMap {
		if relationRaw.ToUserId == toUid {
			// 写redis
			db.AddRelation(r.ctx, uid, toUid)
			return true, nil
		}
	}
	return false, errors.New("NOT FOUND")
}

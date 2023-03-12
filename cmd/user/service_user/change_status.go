package service_user

import (
	"aweme_kitex/cmd/relation/kitex_gen/relation"
	userRPC "aweme_kitex/cmd/user/rpc"
	userDB "aweme_kitex/cmd/user/service_user/db"
	userKafka "aweme_kitex/cmd/user/service_user/kafka"
	constants "aweme_kitex/pkg/constant"
	"context"
	"errors"
	"time"
)

type ChangeFollowService struct {
	ctx context.Context
}

// NewRegisterUserService new RegisterUserService
func NewChangeFollowService(ctx context.Context) *ChangeFollowService {
	return &ChangeFollowService{
		ctx: ctx,
	}
}

func (u *ChangeFollowService) ChangeStatus(userId, toUserId string, isfollow bool) error {
	if userId == toUserId {
		return errors.New("userId is equal to toUserId")
	}
	// isFollow get
	msg := userId + "&" + toUserId
	flag, err := userRPC.QueryRelation(u.ctx, &relation.QueryRelationRequest{
		UserId:   userId,
		ToUserId: toUserId,
		IsFollow: false,
	})
	if err != nil {
		return err
	}
	// 关注操作
	if isfollow {
		// 已经关注了
		if flag {
			return nil
		}
		// redis删除
		userDB.DelCount(u.ctx, userId, toUserId)
		// 发送打消息队列
		err := userKafka.ProduceFollowMsg(constants.KafKaUserAddRelationTopic, msg)
		if err != nil {
			return err
		}
		time.Sleep(constants.SleepTime)
		userDB.DelCount(u.ctx, userId, toUserId)
	} else {
		// 取消关注操作
		if !flag {
			// 已经取消关注了
			return nil
		}
		// redis删除
		userDB.DelCount(u.ctx, userId, toUserId)
		// 发送打消息队列
		err := userKafka.ProduceFollowMsg(constants.KafKaUserDelRelationTopic, msg)
		if err != nil {
			return err
		}
		time.Sleep(constants.SleepTime)
		userDB.DelCount(u.ctx, userId, toUserId)
	}
	return nil
}

package service_relation

import (
	"aweme_kitex/cmd/relation/kitex_gen/relation"
	relationRPC "aweme_kitex/cmd/relation/rpc"
	relation_db "aweme_kitex/cmd/relation/service_relation/db"
	"aweme_kitex/cmd/relation/service_relation/kafka"
	user2 "aweme_kitex/cmd/user/kitex_gen/user"
	"aweme_kitex/cmd/user/service_user/db"
	constants "aweme_kitex/pkg/constant"
	"aweme_kitex/pkg/jwt"
	"context"
	"strings"
	"time"
)

type RelationActionService struct {
	ctx context.Context
}

// NewRelationActionService new RelationActionService
func NewRelationActionService(ctx context.Context) *RelationActionService {
	return &RelationActionService{ctx: ctx}
}

func (s *RelationActionService) RelationAction(req *relation.RelationActionRequest) error {
	uc, err := jwt.AnalyzeToken(req.Token)
	if err != nil {
		return err
	}
	err = relationAction(s.ctx, uc.Id, req.ToUserId, req.ActionType)
	return err
}

func relationAction(ctx context.Context, uid, toUid, action string) error {
	return newRelationActionDataFlow(ctx, uid, toUid, action).do()
}

func newRelationActionDataFlow(ctx context.Context, uid, toUid, action string) *relationActionDataFlow {
	return &relationActionDataFlow{
		ctx:        ctx,
		userId:     uid,
		toUserId:   toUid,
		actionType: action,
	}
}

type relationActionDataFlow struct {
	ctx        context.Context
	userId     string
	toUserId   string
	actionType string
}

func (r *relationActionDataFlow) do() error {
	if _, err := db.NewUserDaoInstance().CheckUserId(r.ctx, []string{r.toUserId}); err != nil {
		return err
	}
	if r.actionType == constants.Follow {
		err := r.createRelation()
		if err != nil {
			return err
		}
	} else if r.actionType == constants.UnFollow {
		if err := r.deleteRelation(); err != nil {
			return err
		}
	}
	return nil
}

func (r *relationActionDataFlow) createRelation() error {
	// 1.use user rpc  去修改关注数目
	err := relationRPC.ChangeFollowCount(r.ctx, &user2.ChangeFollowStatusRequest{
		UserId:   r.userId,
		ToUserId: r.toUserId,
		IsFollow: true,
	})
	if err != nil {
		return err
	}
	// 2. 修改本地数据关系
	// 2.1 修改redis
	relation_db.AddRelation(r.ctx, r.userId, r.toUserId)
	sb := strings.Builder{}
	sb.WriteString(r.userId)
	sb.WriteString("&")
	sb.WriteString(r.toUserId)
	// 2.2 写入kafka供给消费
	kafka.ProduceAddRelation(constants.KafKaRelationAddTopic, sb.String())
	// 2.3 sleep
	time.Sleep(constants.SleepTime)
	// 2.4 del redis
	relation_db.DelRelation(r.ctx, r.userId, r.toUserId)
	return nil
}

func (r *relationActionDataFlow) deleteRelation() error {
	// 1.use user rpc  修改关注数目
	err := relationRPC.ChangeFollowCount(r.ctx, &user2.ChangeFollowStatusRequest{
		UserId:   r.userId,
		ToUserId: r.toUserId,
		IsFollow: false,
	})
	if err != nil {
		return err
	}
	// 2. 修改本地数据关系
	// 2.1 修改redis
	sb := strings.Builder{}
	sb.WriteString(r.userId)
	sb.WriteString("&")
	sb.WriteString(r.toUserId)
	// 2.2 写入kafka供给消费
	kafka.ProduceAddRelation(constants.KafKaRelationDelTopic, sb.String())
	// 2.3 del redis
	relation_db.DelRelation(r.ctx, r.userId, r.toUserId)
	return nil
}

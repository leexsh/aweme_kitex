package service_relation

import (
	"aweme_kitex/cmd/relation/kitex_gen/relation"
	"aweme_kitex/models/dal"
	"aweme_kitex/pkg/jwt"
	"context"
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
	if _, err := dal.NewUserDaoInstance().CheckUserId(r.ctx, []string{r.toUserId}); err != nil {
		return err
	}
	if r.actionType == "1" {
		err := r.createRelation()
		if err != nil {
			return err
		}
	} else if r.actionType == "2" {
		if err := r.deleteRelation(); err != nil {
			return err
		}
	}
	return nil
}

func (r *relationActionDataFlow) createRelation() error {
	err := dal.NewRelationDaoInstance().CreateRelation(r.ctx, r.userId, r.toUserId)
	if err != nil {
		return err
	}
	return nil
}

func (r *relationActionDataFlow) deleteRelation() error {
	err := dal.NewRelationDaoInstance().DeleteRelation(r.ctx, r.userId, r.toUserId)
	if err != nil {
		return nil
	}
	return nil
}

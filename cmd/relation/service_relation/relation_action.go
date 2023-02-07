package service_relation

import (
	"aweme_kitex/cmd/relation/kitex_gen/relation"
	"aweme_kitex/pkg/jwt"
	"aweme_kitex/service"
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
	err = service.RelationAction(uc.Id, req.ToUserId, req.ActionType)
	return err
}

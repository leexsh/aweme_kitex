package main

import (
	relation "aweme_kitex/cmd/relation/kitex_gen/relation"
	"aweme_kitex/cmd/relation/service_relation"
	"aweme_kitex/cmd/relation/service_relation/pack"
	"context"
)

/*
关注操作信息流
	如果actionType等于1，表示当前用户关注其他用户，
		当前用户的关注总数增加，其他用户的粉丝总数增加，
		新建关注记录

	如果actionType等于2，表示当前用户取消关注其他用户
		当前用户的关注总数减少，其他用户的粉丝总数减少，
		删除关注记录
*/

// RelationServiceImpl implements the last service_user interface defined in the IDL.
type RelationServiceImpl struct{}

// RelationAction implements the RelationServiceImpl interface.
func (s *RelationServiceImpl) RelationAction(ctx context.Context, req *relation.RelationActionRequest) (resp *relation.RelationActionResponse, err error) {
	if len(req.Token) == 0 || len(req.ToUserId) == 0 || len(req.ActionType) == 0 {
		return pack.RelationActionResponse(-1, "relation action params error"), nil
	}
	err = service_relation.NewRelationActionService(ctx).RelationAction(req)
	if err != nil {
		return pack.RelationActionResponse(-1, err.Error()), err
	}
	return pack.RelationActionResponse(0, "Relation action success"), nil
}

// FollowList implements the RelationServiceImpl interface.
func (s *RelationServiceImpl) FollowList(ctx context.Context, req *relation.FollowListRequest) (resp *relation.FollowListResponse, err error) {
	userList, err := service_relation.NewFollowListService(ctx).FollowList(req)
	if err != nil {
		return pack.FollowListResponse(-1, err.Error(), nil), nil
	}
	return pack.FollowListResponse(0, "get follower list success", userList), nil
}

// FollowerList implements the RelationServiceImpl interface.
func (s *RelationServiceImpl) FollowerList(ctx context.Context, req *relation.FollowerListRequest) (resp *relation.FollowerListResponse, err error) {
	userList, err := service_relation.NewFollowerListService(ctx).FollowerList(req)
	if err != nil {
		return pack.FollowerListResponse(-1, err.Error(), nil), nil
	}
	return pack.FollowerListResponse(0, "get follower list success", userList), nil
}

// QueryRelation implements the RelationServiceImpl interface.
func (s *RelationServiceImpl) QueryRelation(ctx context.Context, req *relation.QueryRelationRequest) (resp *relation.QueryRelationResponse, err error) {
	isFollow, err := service_relation.NewRelationStatusService(ctx).QueryRelationStatus(req.UserId, req.ToUserId)
	if err != nil {
		return pack.RelationResponse(-1, err.Error(), false), err
	}
	resp.IsFollow = isFollow
	return pack.RelationResponse(0, "QueryInfo successed", isFollow), err
}

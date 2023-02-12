package main

import (
	relation "aweme_kitex/cmd/relation/kitex_gen/relation"
	"aweme_kitex/cmd/relation/service_relation"
	"context"
	"time"
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
	resp = new(relation.RelationActionResponse)
	if len(req.Token) == 0 || len(req.ToUserId) == 0 || len(req.ActionType) == 0 {
		resp.BaseResp.ServiceTime = time.Now().Unix()
		resp.BaseResp.StatusCode = -1
		resp.BaseResp.StatusMsg = "relation action params error"
		return resp, nil
	}
	err = service_relation.NewRelationActionService(ctx).RelationAction(req)
	if err != nil {
		resp.BaseResp.ServiceTime = time.Now().Unix()
		resp.BaseResp.StatusCode = -1
		resp.BaseResp.StatusMsg = err.Error()
		return resp, nil
	}
	return resp, nil
}

// FollowList implements the RelationServiceImpl interface.
func (s *RelationServiceImpl) FollowList(ctx context.Context, req *relation.FollowListRequest) (resp *relation.FollowListResponse, err error) {
	resp = new(relation.FollowListResponse)
	if len(req.Token) == 0 {
		resp.BaseResp.ServiceTime = time.Now().Unix()
		resp.BaseResp.StatusCode = -1
		resp.BaseResp.StatusMsg = "relation action params error"
		return resp, nil
	}
	userList, err := service_relation.NewFollowListService(ctx).FollowList(req)
	if err != nil {
		resp.BaseResp.ServiceTime = time.Now().Unix()
		resp.BaseResp.StatusCode = -1
		resp.BaseResp.StatusMsg = err.Error()
		return resp, nil
	}
	resp.BaseResp.StatusCode = 0
	resp.BaseResp.ServiceTime = time.Now().Unix()
	resp.BaseResp.StatusMsg = "get follow list success"
	resp.UserList = userList
	return resp, nil
}

// FollowerList implements the RelationServiceImpl interface.
func (s *RelationServiceImpl) FollowerList(ctx context.Context, req *relation.FollowerListRequest) (resp *relation.FollowerListResponse, err error) {
	resp = new(relation.FollowerListResponse)
	if len(req.Token) == 0 {
		resp.BaseResp.ServiceTime = time.Now().Unix()
		resp.BaseResp.StatusCode = -1
		resp.BaseResp.StatusMsg = "relation action params error"
		return resp, nil
	}
	userList, err := service_relation.NewFollowerListService(ctx).FollowerList(req)
	if err != nil {
		resp.BaseResp.ServiceTime = time.Now().Unix()
		resp.BaseResp.StatusCode = -1
		resp.BaseResp.StatusMsg = err.Error()
		return resp, nil
	}
	resp.BaseResp.StatusCode = 0
	resp.BaseResp.ServiceTime = time.Now().Unix()
	resp.BaseResp.StatusMsg = "get follow list success"
	resp.UserList = userList
	return resp, nil
	return
}

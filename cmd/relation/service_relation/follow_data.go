package service_relation

import (
	"aweme_kitex/cmd/relation/kitex_gen/user"
	relationRPC "aweme_kitex/cmd/relation/rpc"
	"aweme_kitex/cmd/relation/service_relation/db"
	user2 "aweme_kitex/cmd/user/kitex_gen/user"
	db2 "aweme_kitex/cmd/user/service_user/db"
	"aweme_kitex/pkg/types"
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

func getFollowList(ctx context.Context, userId string) ([]*user.User, error) {
	return newRelationListDataFlow(ctx, userId).getFollow()
}
func getFollowerList(ctx context.Context, userId string) ([]*user.User, error) {
	return newRelationListDataFlow(ctx, userId).getFollower()
}

type relationListDataFlow struct {
	ctx      context.Context
	UserId   string
	UserList []*user.User

	UserRaw     []*user2.User
	RelationMap map[string]*types.RelationRaw
}

func newRelationListDataFlow(ctx context.Context, userId string) *relationListDataFlow {
	return &relationListDataFlow{
		UserId: userId,
		ctx:    ctx,
	}
}

func (r *relationListDataFlow) getFollow() ([]*user.User, error) {
	if err := r.prepareFollowInfo(); err != nil {
		return nil, err
	}
	if err := r.packageFollowInfo(); err != nil {
		return nil, err
	}
	return r.UserList, nil
}

func (r *relationListDataFlow) prepareFollowInfo() error {
	relations, err := db.NewRelationDaoInstance().QueryFollowByUid(r.ctx, r.UserId)
	if err != nil {
		return err
	}
	toUids := make([]string, 0)
	for _, relation := range relations {
		toUids = append(toUids, relation.ToUserId)
	}
	us, err := relationRPC.GetUserInfo(r.ctx, &user2.SingleUserInfoRequest{
		UserIds: toUids,
	})
	// toUsers, err := db.NewUserDaoInstance().QueryUserByIds(r.ctx, toUids)
	if err != nil {
		return err
	}
	r.UserRaw = us
	relationMap, err := db.NewRelationDaoInstance().QueryRelationByIds(r.ctx, r.UserId, toUids)
	r.RelationMap = relationMap
	return nil
}
func (r *relationListDataFlow) packageFollowInfo() error {
	userList := make([]*user.User, 0)
	for _, u := range r.UserRaw {
		var isFollow bool = false

		_, ok := r.RelationMap[u.UserId]
		if ok {
			isFollow = true
		}

		curUer := &user.User{
			u.UserId,
			u.Name,
			u.FollowCount,
			u.FollowerCount,
			isFollow,
		}
		userList = append(userList, curUer)
	}

	r.UserList = userList
	return nil
}

func (r *relationListDataFlow) getFollower() ([]*user.User, error) {
	if _, err := db2.NewUserDaoInstance().CheckUserId(r.ctx, []string{r.UserId}); err != nil {
		return nil, err
	}
	if err := r.prepareFollowerInfo(); err != nil {
		return nil, err
	}
	if err := r.packageFollowerInfo(); err != nil {
		return nil, err
	}
	return r.UserList, nil
}

func (r *relationListDataFlow) prepareFollowerInfo() error {
	// 1.尝试使用redis获取
	// 查询目标用户的被关注记录
	relations, err := db.NewRelationDaoInstance().QueryFollowerById(r.ctx, r.UserId)
	if err != nil {
		return err
	}

	// 获取这些记录的关注方id
	userIds := make([]string, 0)
	for _, relation := range relations {
		userIds = append(userIds, relation.UserId)
	}

	// 获取关注方的信息
	us, err := relationRPC.GetUserInfo(r.ctx, &user2.SingleUserInfoRequest{
		UserIds: userIds,
	})
	r.UserRaw = us

	// 获取当前用户与关注方的关注记录
	relationMap, err := db.NewRelationDaoInstance().QueryRelationByIds(r.ctx, r.UserId, userIds)
	if err != nil {
		return err
	}
	r.RelationMap = relationMap

	return nil
}

func (r *relationListDataFlow) packageFollowerInfo() error {
	userList := make([]*user.User, 0)
	for _, u := range r.UserRaw {
		var isFollow bool = false

		_, ok := r.RelationMap[u.UserId]
		if ok {
			isFollow = true
		}
		curUer := &user.User{
			u.UserId,
			u.Name,
			u.FollowCount,
			u.FollowerCount,
			isFollow,
		}
		userList = append(userList, curUer)
	}

	r.UserList = userList
	return nil
}

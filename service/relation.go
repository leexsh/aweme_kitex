package service

import (
	"aweme_kitex/models"
	"aweme_kitex/models/repository"
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
func RelationAction(uid, toUid, action string) error {
	return newRelationActionDataFlow(uid, toUid, action).do()
}

func newRelationActionDataFlow(uid, toUid, action string) *relationActionDataFlow {
	return &relationActionDataFlow{
		userId:     uid,
		toUserId:   toUid,
		actionType: action,
	}
}

type relationActionDataFlow struct {
	userId     string
	toUserId   string
	actionType string
}

func (r *relationActionDataFlow) do() error {
	if _, err := checkUserId([]string{r.toUserId}); err != nil {
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
	err := repository.NewRelationDaoInstance().CreateRelation(context.Background(), r.userId, r.toUserId)
	if err != nil {
		return err
	}
	return nil
}

func (r *relationActionDataFlow) deleteRelation() error {
	err := repository.NewRelationDaoInstance().DeleteRelation(context.Background(), r.userId, r.toUserId)
	if err != nil {
		return nil
	}
	return nil
}

func GetFollowList(userId string) ([]*models.User, error) {
	return newRelationListDataFlow(userId).getFollow()
}
func GetFollowerList(userId string) ([]*models.User, error) {
	return newRelationListDataFlow(userId).getFollower()
}

type relationListDataFlow struct {
	UserId   string
	UserList []*models.User

	UserRaw     []*models.UserRawData
	RelationMap map[string]*models.RelationRaw
}

func newRelationListDataFlow(userId string) *relationListDataFlow {
	return &relationListDataFlow{
		UserId: userId,
	}
}

func (f *relationListDataFlow) getFollow() ([]*models.User, error) {
	if _, err := checkUserId([]string{f.UserId}); err != nil {
		return nil, err
	}
	if err := f.prepareFollowInfo(); err != nil {
		return nil, err
	}
	if err := f.packageFollowInfo(); err != nil {
		return nil, err
	}
	return f.UserList, nil
}

func (r *relationListDataFlow) prepareFollowInfo() error {
	relations, err := repository.NewRelationDaoInstance().QueryFollowByUid(context.Background(), r.UserId)
	if err != nil {
		return err
	}
	toUids := make([]string, 0)
	for _, relation := range relations {
		toUids = append(toUids, relation.ToUserId)
	}
	toUsers, err := repository.NewUserDaoInstance().QueryUserByIds(context.Background(), toUids)
	if err != nil {
		return err
	}
	r.UserRaw = toUsers
	relationMap, err := repository.NewRelationDaoInstance().QueryRelationByIds(context.Background(), r.UserId, toUids)
	r.RelationMap = relationMap
	return nil
}
func (r *relationListDataFlow) packageFollowInfo() error {
	userList := make([]*models.User, 0)
	for _, user := range r.UserRaw {
		var isFollow bool = false

		_, ok := r.RelationMap[user.UserId]
		if ok {
			isFollow = true
		}
		userList = append(userList, &models.User{
			UserId:        user.UserId,
			Name:          user.Name,
			FollowCount:   user.FollowCount,
			FollowerCount: user.FollowerCount,
			IsFollow:      isFollow,
		})
	}

	r.UserList = userList
	return nil
}

func (f *relationListDataFlow) getFollower() ([]*models.User, error) {
	if _, err := checkUserId([]string{f.UserId}); err != nil {
		return nil, err
	}
	if err := f.prepareFollowerInfo(); err != nil {
		return nil, err
	}
	if err := f.packageFollowerInfo(); err != nil {
		return nil, err
	}
	return f.UserList, nil
}

func (r *relationListDataFlow) prepareFollowerInfo() error {
	// 查询目标用户的被关注记录
	relations, err := repository.NewRelationDaoInstance().QueryFollowerById(context.Background(), r.UserId)
	if err != nil {
		return err
	}

	// 获取这些记录的关注方id
	userIds := make([]string, 0)
	for _, relation := range relations {
		userIds = append(userIds, relation.UserId)
	}

	// 获取关注方的信息
	users, err := repository.NewUserDaoInstance().QueryUserByIds(context.Background(), userIds)
	if err != nil {
		return err
	}
	r.UserRaw = users

	// 获取当前用户与关注方的关注记录
	relationMap, err := repository.NewRelationDaoInstance().QueryRelationByIds(context.Background(), r.UserId, userIds)
	if err != nil {
		return err
	}
	r.RelationMap = relationMap

	return nil
}

func (r *relationListDataFlow) packageFollowerInfo() error {
	userList := make([]*models.User, 0)
	for _, user := range r.UserRaw {
		var isFollow bool = false

		_, ok := r.RelationMap[user.UserId]
		if ok {
			isFollow = true
		}
		userList = append(userList, &models.User{
			UserId:        user.UserId,
			Name:          user.Name,
			FollowCount:   user.FollowCount,
			FollowerCount: user.FollowerCount,
			IsFollow:      isFollow,
		})
	}

	r.UserList = userList
	return nil
}

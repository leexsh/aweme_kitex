package service_user

import (
	"aweme_kitex/cmd/user/kitex_gen/user"
	"aweme_kitex/cmd/user/service_user/db"
	"aweme_kitex/pkg/types"
	"context"
	"errors"
	"strconv"
)

type UserInfoService struct {
	ctx context.Context
}

func NewUserInfoService(ctx context.Context) *UserInfoService {
	return &UserInfoService{
		ctx: ctx,
	}
}

func (s *UserInfoService) packUserInfo(u *types.UserRawData) (*user.User, error) {
	if u == nil {
		return nil, errors.New("user is nil")
	}
	return &user.User{
		UserId:        u.UserId,
		Name:          u.Name,
		FollowCount:   u.FollowCount,
		FollowerCount: u.FollowerCount,
	}, nil
}

func (s *UserInfoService) UserInfo(req *user.UserInfoRequest) (user *user.User, err error) {
	u, err := db.NewUserDaoInstance().QueryUserByUserId(s.ctx, req.UserId)
	if err != nil {
		return nil, err
	}
	return s.packUserInfo(u)
}

func (s *UserInfoService) SingleUserInfo(uid string) (user *user.User, err error) {

	followNum, _ := db.RedisClient.HGet(s.ctx, uid, db.FollowNum).Result()
	followerNum, _ := db.RedisClient.HGet(s.ctx, uid, db.FollowerNum).Result()
	name, _ := db.RedisClient.Get(s.ctx, uid).Result()
	// 1.Redis exist
	if len(name) > 0 && len(followNum) > 0 && len(followerNum) > 0 {
		fInt, _ := strconv.ParseInt(followNum, 10, 64)
		ferInt, _ := strconv.ParseInt(followerNum, 10, 64)

		user.UserId = uid
		user.Name = name
		user.FollowCount = fInt
		user.FollowerCount = ferInt
		return
	} else {
		// 2. redis不存在，查询sql，并写入redis
		u, err := db.NewUserDaoInstance().QueryUserByUserId(s.ctx, uid)
		if err != nil {
			return nil, err
		}
		db.AddName(s.ctx, u.UserId, u.Name)
		db.UpdateCount(s.ctx, u.UserId, u.FollowCount, u.FollowerCount)
		return s.packUserInfo(u)
	}
	return
}

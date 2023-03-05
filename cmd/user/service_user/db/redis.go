package db

import (
	constants "aweme_kitex/pkg/constant"
	"aweme_kitex/pkg/logger"
	"context"

	"github.com/go-redis/redis/v8"
)

const (
	FollowNum   = "FollowNum"
	FollowerNum = "FollowerNum"
)

var RedisClient *redis.Client

func InitRedis() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     "121.5.114.14:6379",
		Password: "123456",            // no password set
		DB:       constants.RedisUser, // use default DB
	})

	pong, _ := RedisClient.Ping(context.Background()).Result()
	logger.Info("pong is: " + pong)
}

func UpdateCount(ctx context.Context, userId string, followNum, followerNum int64) {
	RedisClient.HSet(ctx, userId, FollowNum, followNum)
	RedisClient.HSet(ctx, userId, FollowerNum, followerNum)
	RedisClient.Expire(ctx, userId, constants.RedisExpireTime)
}

func AddName(ctx context.Context, userId, userName string) {
	RedisClient.Set(ctx, userId, userName, constants.RedisExpireTime)
}

// 删除UserId的关注数和ToUserId的粉丝数
func DelCount(ctx context.Context, UserId string, ToUserId string) (bool, error) {
	// 删除粉丝缓存的两者关系
	//scard计算集合大小,SRem移除关系
	RedisClient.HDel(ctx, UserId, FollowNum)
	RedisClient.HDel(ctx, ToUserId, FollowerNum)
	RedisClient.Expire(ctx, UserId, constants.RedisExpireTime)
	RedisClient.Expire(ctx, ToUserId, constants.RedisExpireTime)
	return true, nil
}

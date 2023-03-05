package db

import (
	constants "aweme_kitex/pkg/constant"
	"context"

	"github.com/go-redis/redis/v8"
)

var RelationClient *redis.Client

func InitRedis() {
	RelationClient = redis.NewClient(&redis.Options{
		Addr:     "121.5.114.14:6379",
		Password: "123456",                 // no password set
		DB:       constants.RedisRelation1, // use default DB
	})

}

func AddRelation(ctx context.Context, userId, toUserId string) {
	RelationClient.SAdd(ctx, userId, toUserId)
	RelationClient.Expire(ctx, userId, constants.RedisExpireTime)
}

func DelRelation(ctx context.Context, userId, toUserId string) {
	RelationClient.SRem(ctx, userId, toUserId)
	RelationClient.Expire(ctx, userId, constants.RedisExpireTime)
}

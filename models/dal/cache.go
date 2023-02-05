package dal

import (
	"aweme_kitex/cfg"
	"time"
)

// redis
/*
	1.读
		- 先读缓存，存在则返回
		- 不存在，读取数据库
		- 读完数据库后写缓存
	2.写
		- 先写数据库，
		- 删除缓存
*/
func redisSet(key string, value interface{}, expiration time.Duration) error {
	return cfg.RedisClient.Set(key, value, expiration).Err()
}

func redisGet(key string) (string, error) {
	return cfg.RedisClient.Get(key).Result()
}

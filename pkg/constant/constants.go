package constants

import "time"

const (
	SecretKey   = "secret key"
	IdentiryKey = "id"

	VideoList = "video_list"
	NextTime  = "next_time"

	// 服务名
	ApiServiceName      = "api"
	FeedServiceName     = "feed"
	PublishServiceName  = "publish"
	UserServiceName     = "user"
	FavoriteServiceName = "favorite"
	CommentServiceName  = "comment"
	RelationServiceName = "relation"

	// MySQLDefaultDSN = "gorm:gorm@tcp(localhost:9910)/gorm?chazrset=utf8&parseTime=True&loc=local"
	EtcdAddress = "127.0.0.1:2379"
	ApiAddress  = "127.0.0.1:8080"

	CPURateLimit = 80.0
	DefaultLimit = 10

	// MySQL配置
	MySQLMaxIdleConns    = 10        // 空闲连接池中连接的最大数量
	MySQLMaxOpenConns    = 100       // 打开数据库连接的最大数量
	MySQLConnMaxLifetime = time.Hour // 连接可复用的最大时间
)

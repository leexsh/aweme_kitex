package constants

import "time"

const (
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

	// ETCD
	EtcdAddress     = "127.0.0.1:2379"
	ApiAddress      = "127.0.0.1:8080" // Api层 地址
	FeedAddress     = "127.0.0.1:8081" // Feed 服务地址
	PublishAddress  = "127.0.0.1:8082" // Publish 服务地址
	UserAddress     = "127.0.0.1:8083" // User服务地址
	FavoriteAddress = "127.0.0.1:8084" // Favorite服务地址
	CommentAddress  = "127.0.0.1:8085" // Comment服务地址
	RelationAddress = "127.0.0.1:8086" // Relation服务地址

	// KafKa
	KafkaAddress              = "127.0.0.1:9092" // kafka地址
	KafKaRelationAddTopic     = "relation_add"
	KafKaRelationDelTopic     = "relation_del"
	KafKaUserAddRelationTopic = "user_relation_add"
	KafKaUserDelRelationTopic = "user_relation_del"
	KafKaFavouriteAddTopic    = "favourite_add"
	KafKaFavouriteDelTopic    = "favourite_del"
	KafKaVideoCommentAddTopic = "video_comment_add"
	KafKaVideoCommentDelTopic = "video_comment_del"

	CPURateLimit = 80.0
	DefaultLimit = 10

	// MySQL配置
	MySQLMaxIdleConns    = 10        // 空闲连接池中连接的最大数量
	MySQLMaxOpenConns    = 100       // 打开数据库连接的最大数量
	MySQLConnMaxLifetime = time.Hour // 连接可复用的最大时间
	// kafka配置

	// logger
	LogFileName = "../../public/log/"

	// favorite actiontype,1是点赞，2是取消点赞
	Like   = "1"
	Unlike = "2"
	// comment actiontype,1是增加评论，2是删除评论
	AddComment = "1"
	DelComment = "2"
	// relation actiontypr,1是关注，2是取消关注
	Follow   = "1"
	UnFollow = "2"

	// Redis
	RedisExpireTime = time.Hour * 48
	SleepTime       = time.Millisecond * 500
	RedisFollower   = 0
	RedisFollow     = 1
	RedisRelation1  = 2
	RedisUser       = 3
	RedisCount2     = 4
)

package cfg

import (
	"aweme_kitex/models"
	constants "aweme_kitex/pkg/constant"
	"aweme_kitex/pkg/logger"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/go-redis/redis"
	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv"
	"github.com/tencentyun/cos-go-sdk-v5"
	"gopkg.in/ini.v1"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormopentracing "gorm.io/plugin/opentracing"
)

var (
	DB          *gorm.DB
	COSClient   *cos.Client
	RedisClient *redis.Client
	err         error
)

func Init() error {
	// 读取ini
	path, _ := os.Getwd()
	config, err := ini.Load(path + "/cfg/config.ini")
	if err != nil {
		return err
	}

	// -------------mysql----------------
	port := config.Section("mysql").Key("port").String()
	user := config.Section("mysql").Key("user").String()
	database := config.Section("mysql").Key("database").String()
	err = godotenv.Load("env")
	if err != nil {
		panic(err)
	}
	password := os.Getenv("MYSQL_PASSWD")
	ip := os.Getenv("MYSQL_IP")
	fmt.Println(port, user, password, database)
	dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8mb4&parseTime=True&loc=Local", "root", password, ip, port, database)
	DB, err = gorm.Open(
		mysql.Open(dsn),
		&gorm.Config{
			QueryFields:            true, // 打印sql
			SkipDefaultTransaction: true, // 禁用事务
		},
	)

	if err != nil {
		return err
	}
	if err = DB.Use(gormopentracing.New()); err != nil {
		panic(err)
	}

	err = DB.AutoMigrate(&models.UserRawData{}, &models.VideoRawData{},
		&models.FavouriteRaw{}, &models.CommentRaw{}, &models.RelationRaw{})
	if err != nil {
		panic(err)
	}

	sqlDB, err := DB.DB()
	if err != nil {
		panic(err)
	}

	if err := sqlDB.Ping(); err != nil {
		panic(err)
	}

	// SetMaxIdleConns 设置空闲连接池中连接的最大数量
	sqlDB.SetMaxIdleConns(constants.MySQLMaxIdleConns)

	// SetMaxOpenConns 设置打开数据库连接的最大数量
	sqlDB.SetMaxOpenConns(constants.MySQLMaxOpenConns)

	// SetConnMaxLifetime 设置了连接可复用的最大时间
	sqlDB.SetConnMaxLifetime(constants.MySQLConnMaxLifetime)

	// 	-----------COS--------
	cosAddr := config.Section("cos").Key("cosAddr").String()
	u, _ := url.Parse(cosAddr)
	b := &cos.BaseURL{BucketURL: u}
	COSClient = cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  os.Getenv("COS_ID"),
			SecretKey: os.Getenv("COS_KEY"),
		},
	})

	// -------------redis----------
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     "121.5.114.14:6379",
		Password: "123456", // no password set
		DB:       0,        // use default DB
	})

	pong, err := RedisClient.Ping().Result()
	logger.Info("pong is: " + pong)
	return nil
}

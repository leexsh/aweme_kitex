package cfg

import (
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv"
	"github.com/tencentyun/cos-go-sdk-v5"
	"gopkg.in/ini.v1"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	DB        *gorm.DB
	COSClient *cos.Client
	err       error
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
	return nil
}

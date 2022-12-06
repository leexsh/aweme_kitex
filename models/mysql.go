package models

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv"
	"gopkg.in/ini.v1"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB
var err error

func init() {
	// 读取ini
	path, _ := os.Getwd()
	// config, err := ini.Load(os.Getwd() + "/models/config.ini")
	config, err := ini.Load(path + "/models/config.ini")
	if err != nil {
		fmt.Println("Failed to read file:%v", err)
		os.Exit(-1)
	}
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
		panic(err)
	}
}

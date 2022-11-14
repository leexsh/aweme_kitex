package models

import (
	"fmt"
	"os"

	"gopkg.in/ini.v1"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB
var err error

func init() {
	// 读取ini
	config, err := ini.Load("./models/config.ini")
	if err != nil {
		fmt.Println("Failed to read file:%v", err)
		os.Exit(-1)
	}
	ip := config.Section("mysql").Key("ip").String()
	port := config.Section("mysql").Key("port").String()
	user := config.Section("mysql").Key("user").String()
	password := config.Section("mysql").Key("password").String()
	database := config.Section("mysql").Key("database").String()
	fmt.Println(ip, port, user, password, database)

	dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8mb4&parseTime=True&loc=Local", "root", os.Getenv("mysql_passwd"), "192.168.0.109", 3306, "aweme_community")
	DB, err = gorm.Open(
		mysql.Open(dsn),
		&gorm.Config{
			// 打印sql
			QueryFields: true,
		})
	DB.Debug()
	if err != nil {
		fmt.Println(err)
	}
}

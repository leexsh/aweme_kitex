package models

import (
	"fmt"
	"os"
	"testing"

	"gopkg.in/ini.v1"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestMysql(t *testing.T) {
	config, err := ini.Load("config.ini")
	if err != nil {
		fmt.Println("fail to read, ", err)
		t.Fatal(err)
	}
	ip := config.Section("mysql").Key("ip").String()
	port := config.Section("mysql").Key("port").String()
	user := config.Section("mysql").Key("user").String()
	password := config.Section("mysql").Key("password").String()
	database := config.Section("mysql").Key("database").String()
	fmt.Println(ip, port, user, password, database)

	dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8mb4&parseTime=True&loc=Local", "root", os.Getenv("mysql_passwd"), "121.5.114.14", 3306, "cloud_dist")
	DB, err := gorm.Open(
		mysql.Open(dsn),
		&gorm.Config{
			QueryFields:            true, // 打印sql
			SkipDefaultTransaction: true, // 禁用事务
		},
	)
	if err != nil {
		t.Fatal(err)
		panic(err)
	}
	fmt.Println(DB)
}

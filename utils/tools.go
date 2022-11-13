package utils

import (
	"crypto/md5"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	TEMPLATE = "2006-01-02 15:04:05"
)

// UnixToTime 时间戳->日期
func UnixToTime(timestamp int) string {
	t := time.Unix(int64(timestamp), 0)
	return t.Format(TEMPLATE)
}

// DateToUnix 日期->时间戳
func DateToUnix(str string) int64 {
	t, err := time.ParseInLocation(TEMPLATE, str, time.Local)
	if err != nil {
		return 0
	}
	return t.Unix()
}

// GetUnix 获取时间戳
func GetUnix() int64 {
	return time.Now().Unix()
}

// GetDate 获取当前时期
func GetDate() string {
	return time.Now().Format(TEMPLATE)
}

// 获取年月日
func GetDay() string {
	template := "20060102"
	return time.Now().Format(template)
}

// md5加密
func Md5(str string) string {
	h := md5.New()
	io.WriteString(h, str)
	return fmt.Sprintf("%x", h.Sum(nil))
}

func Int(str string) (int, error) {
	return strconv.Atoi(str)
}

func String(n int) string {
	return strconv.Itoa(n)
}

// 上传图片
func UploadImg(c *gin.Context, picName string) (string, error) {
	// 1.get file
	file, err := c.FormFile(picName)
	if err != nil {
		return "", err
	}

	// 2.check extName
	extName := path.Ext(file.Filename)
	allowExtMap := map[string]struct{}{
		".jgp":  {},
		".png":  {},
		".gif":  {},
		".jpeg": {},
	}
	if _, ok := allowExtMap[extName]; !ok {
		return "", errors.New("文件后缀不合法")
	}

	// 3.create file folder
	day := GetDay()
	dir := "./static/uploads/" + day
	err = os.Mkdir(dir, 0666)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	// 4.create fileName
	fileName := strconv.FormatInt(GetUnix(), 10) + extName

	// 5.upload
	dst := path.Join(dir, fileName)
	c.SaveUploadedFile(file, dst)
	return dst, nil
}

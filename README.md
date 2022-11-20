# Aweme-Kitex
## 1.文件说明

main.go                     加载路由，初始化路由，将路由器连接到 http.Server并开始侦听和服务 HTTP 请求

routers                     存放路由相关配置

controller                  存放具体实现代码

models                      存放数据库配置，共享结构和工具

public                      存放本地视频文件


model/example.sql           自动创建数据库

model/app.ini               通过修改相关信息，自动打开数据库


## 3. 踩坑记录
- os.getenv("MYSQL_PASSWD")读取失败
  - 原因：将MYSQL_PASSWD只配置到了当前用户的.bash中，goland运行的时候，读取不到，读到""
  - 解决方法：利用godotenv，将部分信息配置到.env中，参考https://www.cnblogs.com/zhangmingcheng/p/15802038.html
  
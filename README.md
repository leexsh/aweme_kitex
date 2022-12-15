# Aweme-Kitex
## 1.文件说明

main.go                     加载路由，初始化路由，将路由器连接到 http.Server并开始侦听和服务 HTTP 请求

routers                     存放路由相关配置

controller                  存放具体实现代码

models                      存放数据库配置，共享结构和工具

public                      存放本地视频文件


model/mysql_struct.sql      MYSQL DDL

model/app.ini               数据库配置


## 3. 踩坑记录
1. os.getenv("MYSQL_PASSWD")读取失败 
   - 原因：将MYSQL_PASSWD只配置到了当前用户的.bash中，goland运行的时候，读取不到，读到"";
   - 解决方法：利用godotenv，将部分信息配置到.env中; 
   - 参考: https://www.cnblogs.com/zhangmingcheng/p/15802038.html。


2. JSON omitempty的时候会自动忽略false等 
   - golang在处理json转换时，对于标签omitempty定义的field，如果给它赋得值恰好等于空值（比如：false、0、""、nil指针、nil接口、长度为0的数组、切片、映射），则在转为json之后不会输出这个field;
   - 所以如果属性用了omitempty，前端会不显示;
   - 参考：https://www.jianshu.com/p/ffff11015ccf。
## 4.todo
1. 整理sql信息到mysql_struct.sql


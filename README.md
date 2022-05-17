# douyin

1.文件说明
main.go                     加载路由，初始化路由，将路由器连接到 http.Server并开始侦听和服务 HTTP 请求
routers                     存放路由相关配置
controller                  存放具体实现代码
models                      存放数据库配置，共享结构和工具
public                      存放本地视频文件

model/example.sql           自动创建数据库
$ mysql –u用户名 –p密码 –D数据库<【sql脚本文件路径全名】

model/app.ini               通过修改相关信息，自动打开数据库
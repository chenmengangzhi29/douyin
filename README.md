# douyin
## 一.代码运行<br>
1.更改model/app.ini配置文件

2.初始化MySql数据库
```Shell
mysql -u mysql -p mysql -D douyin < ./model/example.sql
```

3.启动程序
```Shell
go run main.go
```

## 二.文件说明<br>

main.go 
加载路由，初始化路由，将路由器连接到 http.Server并开始侦听和服务 HTTP 请求<br>

routers
存放路由相关配置<br>

controller
存放具体实现代码<br>

models
存放数据库配置，共享结构和工具<br>

public
存放本地视频文件<br>

model/example.sql
自动创建数据库
mysql –u用户名 –p密码 –D数据库<【sql脚本文件路径全名】<br>

model/app.ini
通过修改相关信息，自动打开数据库<br>

## 三、接口文档<br>
https://www.apifox.cn/apidoc/shared-2a880467-5d93-4621-b152-a27bc722058c


# 极简版抖音服务端

## 一、介绍

- 1.基于RPC框架**Kitex**、HTTP框架**Gin**、ORM框架**GORM**的极简版抖音服务端项目
- 2.代码采用api层、service层、dal层三层结构
- 3.使用Kitex构建RPC微服务，Gin构建HTTP服务
- 4.GORM操作MySQL数据库，防止SQL注入，使用事务保证数据一致性，完整性
- 5.使用**ETCD**进行服务注册、服务发现，Jarger进行链路追踪
- 6.使用**MySQL**数据库进行数据存储，并建立索引
- 7.使用**OSS**进行视频对象存储，分片上传视频
- 8.使用**JWT**鉴权，**MD5**密码加密，**ffmpeg**获取视频第一帧当作视频封面
- 9.进行了**单元测试**，api**自动化测试**，[Apifox接口文档](https://www.apifox.cn/apidoc/shared-2a880467-5d93-4621-b152-a27bc722058c)
- 10.[演示视频及客户端apk](https://bytedancecampus1.feishu.cn/docx/doxcn8sH7FS6BWFPzG4GaTf187b)

| 服务名   | 用途                               | 框架                   | 协议       | 注册中心 | 链路追踪 | 数据存储      | 日志   |
| -------- | ---------------------------------- | ---------------------  | ---------- | -------- | -------- | ------------- | ------ |
| api      | http接口，通过RPC客户端调用RPC服务 | `gin` `kitex`          | `http`     | `etcd`   | `jaeger` |               | `klog` |
| feed     | 视频流RPC微服务                    | `kitex` `gorm`         | `protobuf` | `etcd`   | `jaeger` | `mysql`       | `klog` |
| publish  | 视频上传RPC微服务                  | `kitex` `gorm` `ffmpeg`| `protobuf` | `etcd`   | `jaeger` | `mysql` `oss` | `klog` |
| user     | 用户RPC微服务                      | `kitex` `gorm` `jwt`   | `protobuf` | `etcd`   | `jaeger` | `mysql`       | `klog` |
| favorite | 点赞RPC微服务                      | `kitex` `gorm`         | `protobuf` | `etcd`   | `jaeger` | `mysql`       | `klog` |
| comment  | 评论RPC微服务                      | `kitex` `gorm`         | `protobuf` | `etcd`   | `jaeger` | `mysql`       | `klog` |
| relation | 关注RPC微服务                      | `kitex` `gorm`         | `protobuf` | `etcd`   | `jaeger` | `mysql`       | `klog` |

## 二、调用关系图

<img src="https://github.com/chenmengangzhi29/douyin/blob/main/public/%E6%9C%8D%E5%8A%A1%E8%B0%83%E7%94%A8%E5%85%B3%E7%B3%BB%E6%88%AA%E5%9B%BE.jpg" width="2850"  alt=""/>

## 三、数据库ER图

<img src="https://github.com/chenmengangzhi29/douyin/blob/main/public/%E5%BE%AE%E6%9C%8D%E5%8A%A1%E6%95%B0%E6%8D%AE%E5%BA%93%E5%9B%BE.jpg" width="2850"  alt=""/>

## 四、文件目录结构

| 目录| 说明 |
| --- | ---- |
| [cmd/api](https://github.com/chenmengangzhi29/douyin/tree/main/cmd/api) | api层代码，处理外部http请求，通过rpc客户端发起rpc请求|
| [cmd/feed](https://github.com/chenmengangzhi29/douyin/tree/main/cmd/feed) | 视频流微服务，包含获取视频流功能，处理api层rpc请求，调用dal层处理数据 |
| [cmd/publish](https://github.com/chenmengangzhi29/douyin/tree/main/cmd/publish) | publish微服务，包含视频上传、视频列表等功能，处理api层rpc请求，调用dal层处理数据 |
| [cmd/user](https://github.com/chenmengangzhi29/douyin/tree/main/cmd/user) | user微服务，包含用户注册、用户登录、用户信息等功能，处理api层rpc请求，调用dal层处理数据 |
| [cmd/favorite](https://github.com/chenmengangzhi29/douyin/tree/main/cmd/favorite) | favorite微服务，包含点赞、取消点赞、点赞列表等功能，处理api层rpc请求，调用dal层处理数据 |
| [cmd/comment](https://github.com/chenmengangzhi29/douyin/tree/main/cmd/comment) | comment微服务，包含增加评论，删除评论，评论列表等功能，处理api层rpc请求，调用dal层处理数据 |
| [cmd/relation](https://github.com/chenmengangzhi29/douyin/tree/main/cmd/relation) | relation微服务，包含关注、取消关注、关注列表等功能，处理api层rpc请求，调用dal层处理数据 |
| [dal/db](https://github.com/chenmengangzhi29/douyin/tree/main/dal/db) | 使用gorm进行底层数据库操作 |
| [dal/pack](https://github.com/chenmengangzhi29/douyin/tree/main/dal/pack) | 封装gorm结构数据为protobuf结构数据 |
| [idl](https://github.com/chenmengangzhi29/douyin/tree/main/idl) | proto IDL文件 |
| [kitex_gen](https://github.com/chenmengangzhi29/douyin/tree/main/kitex_gen) | kitex框架生成的IDL内容相关代码 |
| [pkg/bound](https://github.com/chenmengangzhi29/douyin/tree/main/pkg/bound) | 限制CPU的相关代码 |
| [pkg/constants](https://github.com/chenmengangzhi29/douyin/tree/main/pkg/constants) | 项目中的配置及常量代码 |
| [pkg/errno](https://github.com/chenmengangzhi29/douyin/tree/main/pkg/errno) | 错误码的代码封装 |
| [pkg/jwt](https://github.com/chenmengangzhi29/douyin/tree/main/pkg/jwt) | jwt鉴权的代码封装 |
| [pkg/middleware](https://github.com/chenmengangzhi29/douyin/tree/main/pkg/middleware) | rpc的中间件 |
| [pkg/oss](https://github.com/chenmengangzhi29/douyin/tree/main/pkg/oss) | oss操作的相关代码 |
| [pkg/tracer](https://github.com/chenmengangzhi29/douyin/tree/main/pkg/tracer) | 初始化jaeger |
| [public](https://github.com/chenmengangzhi29/douyin/tree/main/public) | 存放本地视频 |

## 五、代码运行

### 1.更改配置
更改 constants/constant.go 中的地址配置

### 2.建立基础依赖
```shell
docker-compose up
```

### 3.运行feed微服务
```shell
cd cmd/feed
sh build.sh
sh output/bootstrap.sh
```

### 4.运行publish微服务
```shell
cd cmd/publish
sh build.sh
sh output/bootstrap.sh
```

### 5.运行user微服务
```shell
cd cmd/user
sh build.sh
sh output/bootstrap.sh
```

### 6.运行favorite微服务
```shell
cd cmd/favorite
sh build.sh
sh output/bootstrap.sh
```

### 7.运行comment微服务
```shell
cd cmd/comment
sh build.sh
sh output/bootstrap.sh
```

### 8.运行relation微服务
```shell
cd cmd/relation
sh build.sh
sh output/bootstrap.sh
```

### 9.运行api微服务
```shell
cd cmd/api
chmod +x run.sh
./run.sh
```

### 10.Jaeger链路追踪 

在浏览器上查看`http://127.0.0.1:16686/`

<img src="https://github.com/chenmengangzhi29/douyin/blob/main/public/jaeger%E5%B1%95%E7%A4%BA%E5%9B%BE.jpg" width="2850"  alt=""/>









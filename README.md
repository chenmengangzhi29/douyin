# douyin

## Quick Start

### 1.Change Configuration Information
Change constants/constant.go address config

### 2.Setup Basic Dependence
```shell
docker-compose up
```

### 3.Run Feed RPC Server
```shell
cd cmd/feed
sh build.sh
sh output/bootstrap.sh
```

### 4.Run Publish RPC Server
```shell
cd cmd/publish
sh build.sh
sh output/bootstrap.sh
```

### 5.Run User RPC Server
```shell
cd cmd/user
sh build.sh
sh output/bootstrap.sh
```

### 6.Run Favorite RPC Server
```shell
cd cmd/favorite
sh build.sh
sh output/bootstrap.sh
```

### 7.Run Comment RPC Server
```shell
cd cmd/comment
sh build.sh
sh output/bootstrap.sh
```

### 8.Run Relation RPC Server
```shell
cd cmd/relation
sh build.sh
sh output/bootstrap.sh
```

### 9.Run API Server
```shell
cd cmd/api
chmod +x run.sh
./run.sh
```

### 10.Jaeger 

visit `http://127.0.0.1:16686/` on  browser.


## Apifox Interface Document

https://www.apifox.cn/apidoc/shared-2a880467-5d93-4621-b152-a27bc722058c




package constants

import "time"

const (
	//数据库表名
	VideoTableName    = "video"
	UserTableName     = "user"
	FavoriteTableName = "favorite"
	CommentTableName  = "comment"
	ReltaionTableName = "relation"

	//jwt
	SecretKey   = "secret key"
	IdentiryKey = "id"

	//response字段名
	VideoList = "video_list"
	NextTime  = "next_time"
	UserId    = "user_id"
	Token     = "token"
	User      = "user"

	//rpc服务名
	ApiServiceName      = "api"
	FeedServiceName     = "feed"
	PublishServiceName  = "publish"
	UserServiceName     = "user"
	FavoriteServiceName = "favorite"
	CommentServiceName  = "comment"
	RelationServiceName = "relation"

	//地址
	MySQLDefaultDSN    = "gorm:gorm@tcp(localhost:9910)/gorm?charset=utf8&parseTime=True&loc=local" //MySQL DSN
	EtcdAddress        = "127.0.0.1:2379"                                                           //Etcd 地址
	ApiAddress         = "127.0.0.1:8080"                                                           //Api层 地址
	OssEndPoint        = "oss-cn-shenzhen.aliyuncs.com"                                             //Oss
	OssAccessKeyId     = "oss"
	OssAccessKeySecret = "oss"
	OssBucket          = "dousheng1"

	//Limit
	CPURateLimit = 80.0
	DefaultLimit = 10

	//MySQL配置
	MySQLMaxIdleConns    = 10        //空闲连接池中连接的最大数量
	MySQLMaxOpenConns    = 100       //打开数据库连接的最大数量
	MySQLConnMaxLifetime = time.Hour //连接可复用的最大时间

)

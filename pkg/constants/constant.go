package constants

import "time"

const (
	//地址配置
	MySQLDefaultDSN    = "gorm:gorm@tcp(localhost:9910)/gorm?charset=utf8&parseTime=True&loc=Local" //MySQL DSN
	EtcdAddress        = "127.0.0.1:2379"                                                           //Etcd 地址
	ApiAddress         = ":8080"                                                                    //Api层 地址
	FeedAddress        = "127.0.0.1:8081"                                                           //Feed 服务地址
	PublishAddress     = "127.0.0.1:8082"                                                           //Publish 服务地址
	UserAddress        = "127.0.0.1:8083"                                                           //User服务地址
	FavoriteAddress    = "127.0.0.1:8084"                                                           //Favorite服务地址
	CommentAddress     = "127.0.0.1:8085"                                                           //Comment服务地址
	RelationAddress    = "127.0.0.1:8086"                                                           //Relation服务地址
	OssEndPoint        = "oss-cn-shenzhen.aliyuncs.com"                                             //Oss
	OssAccessKeyId     = "oss"
	OssAccessKeySecret = "oss"
	OssBucket          = "dousheng1"

	//数据库表名
	VideoTableName    = "video"
	UserTableName     = "user"
	FavoriteTableName = "favorite"
	CommentTableName  = "comment"
	ReltaionTableName = "relation"

	//jwt
	SecretKey   = "secret key"
	IdentiryKey = "id"

	//时间字段格式
	TimeFormat = "2006-01-02 15:04:05"

	//favorite actiontype,1是点赞，2是取消点赞
	Like   = 1
	Unlike = 2
	//comment actiontype,1是增加评论，2是删除评论
	AddComment = 1
	DelComment = 2
	//relation actiontypr,1是关注，2是取消关注
	Follow   = 1
	UnFollow = 2

	//rpc服务名
	ApiServiceName      = "api"
	FeedServiceName     = "feed"
	PublishServiceName  = "publish"
	UserServiceName     = "user"
	FavoriteServiceName = "favorite"
	CommentServiceName  = "comment"
	RelationServiceName = "relation"

	//Limit
	CPURateLimit = 80.0
	DefaultLimit = 10

	//MySQL配置
	MySQLMaxIdleConns    = 10        //空闲连接池中连接的最大数量
	MySQLMaxOpenConns    = 100       //打开数据库连接的最大数量
	MySQLConnMaxLifetime = time.Hour //连接可复用的最大时间

)

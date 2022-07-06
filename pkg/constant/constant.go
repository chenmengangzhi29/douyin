package constant

const (
	//数据库表名
	VideoTableName    = "video"
	UserTableName     = "user"
	FavoriteTableName = "favorite"
	CommentTableName  = "comment"
	ReltaionTableName = "relation"

	SecretKey   = "secret key"
	IdentiryKey = "id"

	//服务名
	ApiServiceName      = "apiService"
	FeedServiceName     = "feedService"
	PublishServiceName  = "publishService"
	UserServiceName     = "userService"
	FavoriteServiceName = "favoriteService"
	CommentServiceName  = "commentService"
	RelationServiceName = "relationService"

	MySQLDefaultDSN = "gorm:gorm@tcp(localhost:9910)/gorm?charset=utf8&parseTime=True&loc=local"
	EtcdAddress     = "127.0.0.1:2379"
	ApiAddress      = "127.0.0.1:8080"
)

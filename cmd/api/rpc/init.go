package rpc

// InitRPC init rpc client
func InitRPC() {
	initFeedRpc()
	initPublishRpc()
	initUserRpc()
	initFavoriteRpc()
	initCommentRpc()
	initRelationRpc()
}

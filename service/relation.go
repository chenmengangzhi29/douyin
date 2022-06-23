package service

//关注操作信息流
//
//如果actionType等于1，表示当前用户关注其他用户，
//当前用户的关注总数增加，其他用户的粉丝总数增加，
//查询是否有其他用户关注当前用户的关注记录，有则在该记录上将状态位更新为1，没有则新建一条关注记录
//
//如果actionType等于2，表示当前用户取消关注其他用户
//当前用户的关注总数减少，其他用户的粉丝总数减少，
//找到当前用户关注其他用户的关注记录，若其他用户关注当前用户则改变状态位，否则删除该关注记录
func RelationActionData(token string, toUserId int64, actionType int64) error {
	return NewRelationActionDataFlow(token, toUserId, actionType).Do()
}

func NewRelationActionDataFlow(token string, toUserId int64, actionType int64) *RelationActionDataFlow {
	return &RelationActionDataFlow{
		Token:      token,
		ToUserId:   toUserId,
		ActionType: actionType,
	}
}

type RelationActionDataFlow struct {
	Token      string
	ToUserId   int64
	ActionType int64

	CurrentId int64
}

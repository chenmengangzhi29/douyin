package handler

import (
	"douyin/model"
	"strconv"
)

//处理关注操作传入参数
func MakeRelationAction(token string, toUserIdStr string, actionTypeStr string) model.Response {
	if len := len(token); len <= 0 || len > 64 {
		return model.Response{StatusCode: -1, StatusMsg: "token length out of range(0,64]"}
	}

	toUserId, err := strconv.ParseInt(toUserIdStr, 10, 64)
	if err != nil {
		return model.Response{StatusCode: -1, StatusMsg: err.Error()}
	}

	actionType, err := strconv.ParseInt(actionTypeStr, 10, 64)
	if err != nil {
		return model.Response{StatusCode: -1, StatusMsg: err.Error()}
	}

	err = service.RelationActionData(token, toUserId, actionType)
	if err != nil {
		return model.Response{StatusCode: -1, StatusMsg: err.Error()}
	}

	return model.Response{StatusCode: 0, StatusMsg: "relation action success"}
}

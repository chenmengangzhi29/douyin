package handler

import (
	"douyin/model"
	"douyin/service"
	"strconv"
)

//处理关注操作传入参数
func MakeRelationAction(token string, toUserIdStr string, actionTypeStr string) *model.Response {
	if len := len(token); len <= 0 || len > 64 {
		return &model.Response{StatusCode: -1, StatusMsg: "token length out of range(0,64]"}
	}

	toUserId, err := strconv.ParseInt(toUserIdStr, 10, 64)
	if err != nil {
		return &model.Response{StatusCode: -1, StatusMsg: err.Error()}
	}

	actionType, err := strconv.ParseInt(actionTypeStr, 10, 64)
	if err != nil {
		return &model.Response{StatusCode: -1, StatusMsg: err.Error()}
	}
	if actionType != 1 && actionType != 2 {
		return &model.Response{StatusCode: -1, StatusMsg: "action type error"}
	}

	err = service.RelationActionData(token, toUserId, actionType)
	if err != nil {
		return &model.Response{StatusCode: -1, StatusMsg: err.Error()}
	}

	return &model.Response{StatusCode: 0, StatusMsg: "relation action success"}
}

//处理关注列表传入参数
func ShowFollowList(token string, userIdStr string) *model.RelationListResponse {
	if len := len(token); len <= 0 || len > 64 {
		return &model.RelationListResponse{StatusCode: -1, StatusMsg: "token length out of range(0,64]"}
	}

	userId, err := strconv.ParseInt(userIdStr, 10, 64)
	if err != nil {
		return &model.RelationListResponse{StatusCode: -1, StatusMsg: err.Error()}
	}

	userList, err := service.FollowListData(token, userId)
	if err != nil {
		return &model.RelationListResponse{StatusCode: -1, StatusMsg: err.Error()}
	}

	return &model.RelationListResponse{
		StatusCode: 0,
		StatusMsg:  "follow list success",
		UserList:   userList,
	}
}

//处理粉丝列表传入参数
func ShowFollowerList(token string, userIdStr string) *model.RelationListResponse {
	if len := len(token); len <= 0 || len > 64 {
		return &model.RelationListResponse{StatusCode: -1, StatusMsg: "token length out of range(0,64]"}
	}

	userId, err := strconv.ParseInt(userIdStr, 10, 64)
	if err != nil {
		return &model.RelationListResponse{StatusCode: -1, StatusMsg: err.Error()}
	}

	userList, err := service.FollowerListData(token, userId)
	if err != nil {
		return &model.RelationListResponse{StatusCode: -1, StatusMsg: err.Error()}
	}

	return &model.RelationListResponse{
		StatusCode: 0,
		StatusMsg:  "follower list success",
		UserList:   userList,
	}
}

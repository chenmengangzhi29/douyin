package handlers

import (
	"context"
	"strconv"

	"github.com/chenmengangzhi29/douyin/cmd/api/rpc"
	"github.com/chenmengangzhi29/douyin/kitex_gen/relation"
	"github.com/chenmengangzhi29/douyin/pkg/constants"
	"github.com/chenmengangzhi29/douyin/pkg/errno"
	"github.com/gin-gonic/gin"
)

//RelationAction implement follow and unfollow actions
func RelationAction(c *gin.Context) {
	token := c.Query("token")
	toUserIdStr := c.Query("to_user_id")
	actionTypeStr := c.Query("action_type")

	if len(token) == 0 {
		SendResponse(c, errno.ParamErr)
		return
	}

	toUserId, err := strconv.ParseInt(toUserIdStr, 10, 64)
	if err != nil {
		SendResponse(c, errno.ParamParseErr)
		return
	}

	actionType, err := strconv.ParseInt(actionTypeStr, 10, 64)
	if err != nil {
		SendResponse(c, errno.ParamParseErr)
		return
	}
	if actionType != constants.Follow && actionType != constants.UnFollow {
		SendResponse(c, errno.ActionTypeErr)
		return
	}

	req := &relation.RelationActionRequest{Token: token, ToUserId: toUserId, ActionType: int32(actionType)}
	err = rpc.RelationAction(context.Background(), req)
	if err != nil {
		SendResponse(c, err)
		return
	}
	SendResponse(c, errno.Success)
}

// Followlist get user follow list info
func FollowList(c *gin.Context) {
	token := c.Query("token")
	userIdStr := c.Query("user_id")

	userId, err := strconv.ParseInt(userIdStr, 10, 64)
	if err != nil {
		SendResponse(c, errno.ParamParseErr)
		return
	}

	req := &relation.FollowListRequest{Token: token, UserId: userId}
	userList, err := rpc.FollowList(context.Background(), req)
	if err != nil {
		SendResponse(c, err)
		return
	}
	SendRelationListResponse(c, errno.Success, userList)
}

// FollowerList get user follower list info
func FollowerList(c *gin.Context) {
	token := c.Query("token")
	userIdStr := c.Query("user_id")

	userId, err := strconv.ParseInt(userIdStr, 10, 64)
	if err != nil {
		SendResponse(c, errno.ParamParseErr)
		return
	}

	req := &relation.FollowerListRequest{Token: token, UserId: userId}
	userList, err := rpc.FollowerList(context.Background(), req)
	if err != nil {
		SendResponse(c, err)
		return
	}
	SendRelationListResponse(c, errno.Success, userList)
}

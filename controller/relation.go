package controller

import (
	"douyin/handler"
	"net/http"

	"github.com/gin-gonic/gin"
)

//获取关注操作传入参数
func RelationAction(c *gin.Context) {
	token := c.Query("token")
	toUserIdStr := c.Query("to_user_id")
	actionTypeStr := c.Query("action_type")

	relationActionResponse := handler.MakeRelationAction(token, toUserIdStr, actionTypeStr)

	c.JSON(http.StatusOK, relationActionResponse)
}

//获取关注列表传入参数
func FollowList(c *gin.Context) {
	token := c.Query("token")
	userIdStr := c.Query("user_id")

	followListResponse := handler.ShowFollowList(token, userIdStr)

	c.JSON(http.StatusOK, followListResponse)

}

//获取粉丝列表传入参数
func FollowerList(c *gin.Context) {
	token := c.Query("token")
	userIdStr := c.Query("user_id")

	followerListResponse := handler.ShowFollowerList(token, userIdStr)

	c.JSON(http.StatusOK, followerListResponse)
}

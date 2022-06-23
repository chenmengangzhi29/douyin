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

func FollowList(c *gin.Context) {

}

func FollowerList(c *gin.Context) {

}

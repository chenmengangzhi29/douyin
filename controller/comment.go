package controller

import (
	"douyin/handler"
	"douyin/model"
	"net/http"
	"reflect"
	"strconv"

	"github.com/gin-gonic/gin"
)

//获取评论操作传入参数
func CommentAction(c *gin.Context) {
	token := c.Query("token")
	videoIdStr := c.Query("video_id")
	actionTypeStr := c.Query("action_type")

	actionType, err := strconv.ParseInt(actionTypeStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, model.Response{StatusCode: -1, StatusMsg: err.Error()})
	}

	if reflect.DeepEqual(actionType, 1) {
		commentText := c.Query("comment_text")
		commentActionResponse := handler.CreateComment(token, videoIdStr, commentText)
		c.JSON(http.StatusOK, commentActionResponse)
	} else if reflect.DeepEqual(actionType, 2) {
		commentIdStr := c.Query("comment_id")
		commentActionResponse := handler.DeleteComment(token, videoIdStr, commentIdStr)
		c.JSON(http.StatusOK, commentActionResponse)
	} else {
		c.JSON(http.StatusOK, model.Response{StatusCode: -1, StatusMsg: "action type error"})
	}
}

//获取评论列表的传入参数
func CommentList(c *gin.Context) {
	token := c.Query("token")
	videoIdStr := c.Query("video_id")

	commentListResponse := handler.ShowCommentList(token, videoIdStr)

	c.JSON(http.StatusOK, commentListResponse)
}

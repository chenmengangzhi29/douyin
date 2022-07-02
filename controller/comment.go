package controller

import (
	"douyin/handler"
	"douyin/model"
	"douyin/util/logger"
	"net/http"
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
		logger.Errorf("parse %v to int fail, %v", actionTypeStr, err.Error())
		c.JSON(http.StatusOK, &model.Response{StatusCode: -1, StatusMsg: err.Error()})
	}

	if actionType == 1 {
		commentText := c.Query("comment_text")
		commentActionResponse := handler.CreateComment(token, videoIdStr, commentText)
		c.JSON(http.StatusOK, commentActionResponse)
	} else if actionType == 2 {
		commentIdStr := c.Query("comment_id")
		commentActionResponse := handler.DeleteComment(token, videoIdStr, commentIdStr)
		c.JSON(http.StatusOK, commentActionResponse)
	} else {
		logger.Errorf("actionType = %v not equal 1 and 2", actionType)
		c.JSON(http.StatusOK, &model.Response{StatusCode: -1, StatusMsg: "action type error"})
	}
}

//获取评论列表的传入参数
func CommentList(c *gin.Context) {
	token := c.Query("token")
	videoIdStr := c.Query("video_id")

	commentListResponse := handler.ShowCommentList(token, videoIdStr)

	c.JSON(http.StatusOK, commentListResponse)
}

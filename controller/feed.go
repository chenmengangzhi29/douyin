package controller

import (
	"douyin/handler"
	"douyin/util/logger"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

//------------------controller--------------------------------------------
// Feed same demo video list for every request
// 该层功能包括获取传入参数，向handler获取视频信息，返回响应信息
func Feed(c *gin.Context) {
	logger.Info("/feed/")
	//获取传入参数
	token := c.DefaultQuery("token", "defaultToken")
	defaultTime := time.Now().Unix()
	defaultTimeStr := strconv.Itoa(int(defaultTime))
	latestTimeStr := c.DefaultQuery("latest_time", defaultTimeStr)

	//获取视频
	feedResponse := handler.QueryVideoFeed(token, latestTimeStr)

	logger.Info(&feedResponse)
	c.JSON(http.StatusOK, feedResponse)

}

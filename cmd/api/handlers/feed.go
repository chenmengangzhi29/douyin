package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/chenmengangzhi29/douyin/kitex_gen/feed"
	"github.com/gin-gonic/gin"
)

func Feed(c *gin.Context) {
	var FeedVar FeedRequest

	token := c.DefaultQuery("token", "")
	defaultTime := time.Now().Unix()
	defaultTimeStr := strconv.Itoa(int(defaultTime))
	latestTimeStr := c.DefaultQuery("latest_time", defaultTimeStr)

	//处理传入参数
	latestTime, err := strconv.ParseInt(latestTimeStr, 10, 64)
	if err != nil {
		SendResponse(c, err, nil)
		return
	}

	FeedVar.Token = token
	FeedVar.LatestTime = latestTime

	req := &feed.FeedRequest{LatestTime: FeedVar.LatestTime, Token: FeedVar.Token}

	c.JSON(http.StatusOK, feedResponse)
}

package handlers

import (
	"context"
	"strconv"
	"time"

	"github.com/chenmengangzhi29/douyin/cmd/api/rpc"
	"github.com/chenmengangzhi29/douyin/kitex_gen/feed"
	"github.com/chenmengangzhi29/douyin/pkg/errno"
	"github.com/gin-gonic/gin"
)

//Feed get video feed data
func Feed(c *gin.Context) {

	token := c.DefaultQuery("token", "")
	defaultTime := time.Now().UnixMilli()
	defaultTimeStr := strconv.Itoa(int(defaultTime))
	latestTimeStr := c.DefaultQuery("latest_time", defaultTimeStr)

	//处理传入参数
	latestTime, err := strconv.ParseInt(latestTimeStr, 10, 64)
	if err != nil {
		SendResponse(c, err)
		return
	}

	req := &feed.FeedRequest{LatestTime: latestTime, Token: token}
	video, nextTime, err := rpc.Feed(context.Background(), req)
	if err != nil {
		SendResponse(c, err)
		return
	}
	SendFeedResponse(c, errno.Success, video, nextTime)
}

package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func Feed(c *gin.Context) {
	var feedVar struct {
		Token      string `json:"token" form:"token"`
		LatestTime int64  `json:"latest_time" form:"latest_time"`
	}

	token := c.DefaultQuery("token", "")
	defaultTime := time.Now().Unix()
	defaultTimeStr := strconv.Itoa(int(defaultTime))
	latestTimeStr := c.DefaultQuery("latest_time", defaultTimeStr)

	c.JSON(http.StatusOK, feedResponse)
}

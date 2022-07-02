package controller

import (
	"douyin/handler"
	"douyin/model"
	"douyin/util/logger"
	"net/http"

	"github.com/gin-gonic/gin"
)

//获取传入参数，通过handler层上传视频数据
func Publish(c *gin.Context) {
	token := c.PostForm("token")
	title := c.PostForm("title")

	data, _, err := c.Request.FormFile("data")
	if err != nil {
		c.JSON(http.StatusOK, &model.Response{
			StatusCode: -1,
			StatusMsg:  err.Error(),
		})
		return
	}
	defer data.Close()

	logger.Info("publish video")
	publishVideoResponse := handler.PublishVideoData(token, data, title)

	logger.Info(&publishVideoResponse)
	c.JSON(http.StatusOK, publishVideoResponse)
}

//-------------------controller---------------------------------
//获取传入参数，调用handler层函数查询视频列表
func PublishList(c *gin.Context) {
	token := c.Query("token")
	userIdStr := c.Query("user_id")

	videoListResponse := handler.QueryVideoList(token, userIdStr)

	c.JSON(http.StatusOK, videoListResponse)

}

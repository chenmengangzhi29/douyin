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
	data, err := c.FormFile("data")
	if err != nil {
		c.JSON(http.StatusOK, &model.Response{
			StatusCode: -1,
			StatusMsg:  err.Error(),
		})
		return
	}
	title := c.PostForm("title")

	logger.Info("publish video")
	publishVideoResponse := handler.PublishVideoData(token, data, title, c)

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

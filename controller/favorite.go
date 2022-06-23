package controller

import (
	"douyin/handler"
	"net/http"

	"github.com/gin-gonic/gin"
)

//获取点赞操作的传入参数
func FavoriteAction(c *gin.Context) {
	token := c.Query("token")
	videoIdStr := c.Query("video_id")
	actionTypeStr := c.Query("action_type")

	response := handler.HandlerFavoriteAction(token, videoIdStr, actionTypeStr)

	c.JSON(http.StatusOK, response)

}

//获取点赞列表的传入参数
func FavoriteList(c *gin.Context) {
	userIdStr := c.Query("user_id")
	token := c.Query("token")

	favoriteListResponse := handler.HandlerFavoriteList(userIdStr, token)

	c.JSON(http.StatusOK, favoriteListResponse)
}

package routers

//路由定义

import (
	"douyin/controller"

	"github.com/gin-gonic/gin"
)

func Favorite(apiRouter *gin.RouterGroup) {
	apiRouter.POST("/favorite/action/", controller.FavoriteAction)
	apiRouter.GET("/favorite/list/", controller.FavoriteList)
}

package routers

//路由定义

import (
	"douyin/controller"

	"github.com/gin-gonic/gin"
)

func Comment(apiRouter *gin.RouterGroup) {
	apiRouter.POST("/comment/action/", controller.CommentAction)
	apiRouter.GET("/comment/list/", controller.CommentList)
}

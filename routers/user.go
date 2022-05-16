package routers

//路由定义

import (
	"douyin/controller"

	"github.com/gin-gonic/gin"
)

func User(apiRouter *gin.RouterGroup) {
	apiRouter.GET("/user/", controller.UserInfo)
	apiRouter.POST("/user/register/", controller.Register)
	apiRouter.POST("/user/login/", controller.Login)
}

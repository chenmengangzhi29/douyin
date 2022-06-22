package controller

import (
	"net/http"

	"douyin/handler"

	"github.com/gin-gonic/gin"
)

//获取用户注册传入参数
func Register(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	userRegisterResponse := handler.RegisterUser(username, password)

	c.JSON(http.StatusOK, userRegisterResponse)
}

//获取用户登陆传入参数
func Login(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	userLoginResponse := handler.LoginUser(username, password)

	c.JSON(http.StatusOK, userLoginResponse)
}

//获取用户信息传入参数
func UserInfo(c *gin.Context) {
	userIdStr := c.Query("user_id")
	token := c.DefaultQuery("token", "defaultToken")

	userInfoResponse := handler.UserInfoData(userIdStr, token)

	c.JSON(http.StatusOK, userInfoResponse)
}

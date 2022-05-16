package routers

import (
	"github.com/gin-gonic/gin"
)

type Option func(*gin.RouterGroup)

var options = []Option{}

// 注册app的路由配置
func Include(opts ...Option) {
	options = append(options, opts...)
}

// 初始化
func Init() *gin.Engine {
	r := gin.New()
	r.Static("/static", "./public")
	apiRouter := r.Group("/douyin")
	for _, opt := range options {
		opt(apiRouter)
	}
	return r
}

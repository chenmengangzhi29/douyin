package main

import (
	"context"
	"net/http"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/chenmengangzhi29/douyin/cmd/api/handlers"
	"github.com/chenmengangzhi29/douyin/cmd/api/rpc"
	"github.com/chenmengangzhi29/douyin/pkg/constant"
	"github.com/chenmengangzhi29/douyin/pkg/logger"
	"github.com/chenmengangzhi29/douyin/pkg/tracer"
	"github.com/gin-gonic/gin"
)

func Init() {
	tracer.InitJaeger(constant.ApiServiceName)
	rpc.InitRPC()
}

func main() {
	Init()
	r := gin.New()
	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Key:        []byte(constant.SecretKey),
		Timeout:    time.Hour,
		MaxRefresh: time.Hour,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(int64); ok {
				return jwt.MapClaims{
					constant.IdentiryKey: v,
				}
			}
			return jwt.MapClaims{}
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var loginVar handlers.UserLoginParam
			if err := c.ShouldBind(&loginVar); err != nil {
				return "", jwt.ErrMissingLoginValues
			}
			if len(loginVar.Username) == 0 || len(loginVar.Password) == 0 {
				return "", jwt.ErrMissingLoginValues
			}
			return rpc.CheckUser(context.Background(), &userService.CheckUserRequest{Username: loginVar.Username, Password: loginVar.Password})
		},
		TokenLookup:   "header: Authorization, query: token, cookie: jwt",
		TokenHeadName: "Bearer",
		TimeFunc:      time.Now,
	})

	if err != nil {
		logger.Fatal("JWT Error:" + err.Error())
	}

	r.NoRoute(authMiddleware.MiddlewareFunc(), func(c *gin.Context) {
		claims := jwt.ExtractClaims(c)
		logger.Infof("NoRoute claims: %#v\n", claims)
		c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
	})

	douyin := r.Group("/douyin")
	douyin.GET("/feed/", handlers.Feed)

	publish := douyin.Group("/publish")
	publish.Use(authMiddleware.MiddlewareFunc())
	publish.POST("/action/", handlers.Publish)
	publish.GET("/list/", handlers.PublishList)

	user := douyin.Group("/user")
	user.GET("/", handlers.UserInfo)
	user.POST("/register/", handlers.Register)
	user.POST("/login/", authMiddleware.LoginHandler)

	favorite := douyin.Group("/favorite")
	favorite.Use(authMiddleware.MiddlewareFunc())
	favorite.POST("/action/", handlers.FavoriteAction)
	favorite.GET("/list/", handlers.FavoriteList)

	comment := douyin.Group("/comment")
	comment.Use(authMiddleware.MiddlewareFunc())
	comment.POST("/action/", handlers.CommentAction)
	comment.GET("/list/", handlers.CommentList)

	relation := douyin.Group("/relation")
	relation.POST("/action/", handlers.RelationAction)
	relation.GET("/follow/list/", handlers.FollowList)
	relation.GET("/follower/list/", handlers.FollowerList)

	if err := http.ListenAndServe(constant.ApiAddress, r); err != nil {
		logger.Fatal(err)
	}
}

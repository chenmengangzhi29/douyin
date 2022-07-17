package main

import (
	"context"
	"net/http"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/chenmengangzhi29/douyin/cmd/api/handlers"
	"github.com/chenmengangzhi29/douyin/cmd/api/rpc"
	"github.com/chenmengangzhi29/douyin/kitex_gen/user"
	"github.com/chenmengangzhi29/douyin/pkg/constants"
	"github.com/chenmengangzhi29/douyin/pkg/errno"
	JWT "github.com/chenmengangzhi29/douyin/pkg/jwt"
	"github.com/chenmengangzhi29/douyin/pkg/tracer"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/gin-gonic/gin"
)

func Init() {
	tracer.InitJaeger(constants.ApiServiceName)
	rpc.InitRPC()
}

func main() {
	Init()
	r := gin.New()
	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Key:        []byte(constants.SecretKey),
		Timeout:    time.Hour,
		MaxRefresh: time.Hour,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(int64); ok {
				return jwt.MapClaims{
					constants.IdentiryKey: v,
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
			return rpc.CheckUser(context.Background(), &user.CheckUserRequest{Username: loginVar.Username, Password: loginVar.Password})
		},
		LoginResponse: func(c *gin.Context, code int, message string, time time.Time) {
			Jwt := JWT.NewJWT([]byte(constants.SecretKey))
			userId, err := Jwt.CheckToken(message)
			if err != nil {
				panic(err)
			}
			handlers.SendUserResponse(c, errno.Success, userId, message)
		},
		TokenLookup:   "header: Authorization, query: token, cookie: jwt",
		TokenHeadName: "Bearer",
		TimeFunc:      time.Now,
	})

	if err != nil {
		klog.Fatal("JWT Error:" + err.Error())
	}

	douyin := r.Group("/douyin")
	douyin.GET("/feed/", handlers.Feed)

	publish := douyin.Group("/publish")
	publish.POST("/action/", handlers.Publish)
	publish.GET("/list/", handlers.PublishList)

	user := douyin.Group("/user")
	user.GET("/", handlers.UserInfo)
	user.POST("/register/", handlers.Register)
	user.POST("/login/", authMiddleware.LoginHandler)

	favorite := douyin.Group("/favorite")
	favorite.POST("/action/", handlers.FavoriteAction)
	favorite.GET("/list/", handlers.FavoriteList)

	comment := douyin.Group("/comment")
	comment.POST("/action/", handlers.CommentAction)
	comment.GET("/list/", handlers.CommentList)

	relation := douyin.Group("/relation")
	relation.POST("/action/", handlers.RelationAction)
	relation.GET("/follow/list/", handlers.FollowList)
	relation.GET("/follower/list/", handlers.FollowerList)

	if err := http.ListenAndServe(constants.ApiAddress, r); err != nil {
		klog.Fatal(err)
	}
}

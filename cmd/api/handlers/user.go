package handlers

import (
	"context"
	"strconv"

	"github.com/chenmengangzhi29/douyin/cmd/api/rpc"
	"github.com/chenmengangzhi29/douyin/kitex_gen/user"
	"github.com/chenmengangzhi29/douyin/pkg/constants"
	"github.com/chenmengangzhi29/douyin/pkg/errno"
	"github.com/gin-gonic/gin"
)

// Register register uesr info
func Register(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	if len := len(username); len <= 0 || len > 32 {
		SendResponse(c, errno.UserNameValidationErr, nil)
		return
	}

	if len := len(password); len <= 0 || len > 32 {
		SendResponse(c, errno.PasswordValidationErr, nil)
		return
	}

	userId, token, err := rpc.RegisterUser(context.Background(), &user.RegisterUserRequest{
		Username: username,
		Password: password,
	})
	if err != nil {
		SendResponse(c, errno.ConvertErr(err), nil)
		return
	}
	SendResponse(c, errno.Success, map[string]interface{}{constants.Token: token, constants.UserId: userId})
}

// UserInfo get user info
func UserInfo(c *gin.Context) {
	userIdStr := c.Query("user_id")
	token := c.DefaultQuery("token", "")

	userId, err := strconv.ParseInt(userIdStr, 10, 64)
	if err != nil {
		SendResponse(c, errno.ParamParseErr, nil)
		return
	}

	user, err := rpc.UserInfo(context.Background(), &user.UserInfoRequest{
		UserId: userId,
		Token:  token,
	})
	if err != nil {
		SendResponse(c, errno.ConvertErr(err), nil)
	}
	SendResponse(c, errno.Success, map[string]interface{}{constants.User: user})
}

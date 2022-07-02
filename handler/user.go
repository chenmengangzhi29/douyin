package handler

import (
	"douyin/model"
	"douyin/service"
	"strconv"
)

//检查并处理用户注册传入参数
func RegisterUser(username string, password string) *model.UserResponse {
	if len := len(username); len <= 0 || len > 32 {
		return &model.UserResponse{
			StatusCode: -1,
			StatusMsg:  "username length out of range(0,32]",
		}
	}

	if len := len(password); len <= 0 || len > 32 {
		return &model.UserResponse{
			StatusCode: -1,
			StatusMsg:  "password length out of range(0,32]",
		}
	}

	userId, token, err := service.RegisterUserData(username, password)
	if err != nil {
		return &model.UserResponse{
			StatusCode: -1, StatusMsg: err.Error(),
		}
	}

	return &model.UserResponse{
		StatusCode: 0,
		StatusMsg:  "register user success",
		UserId:     userId,
		Token:      token,
	}
}

//检查并处理用户登陆传入参数
func LoginUser(username string, password string) *model.UserResponse {
	if len := len(username); len <= 0 || len > 32 {
		return &model.UserResponse{
			StatusCode: -1,
			StatusMsg:  "username length out of range(0,32]",
		}
	}

	if len := len(password); len <= 0 || len > 32 {
		return &model.UserResponse{
			StatusCode: -1,
			StatusMsg:  "password length out of range(0,32]",
		}
	}

	userId, token, err := service.LoginUserData(username, password)
	if err != nil {
		return &model.UserResponse{
			StatusCode: -1, StatusMsg: err.Error(),
		}
	}

	return &model.UserResponse{
		StatusCode: 0,
		StatusMsg:  "login user success",
		UserId:     userId,
		Token:      token,
	}

}

//检查并处理用户信息传入参数
func UserInfoData(userIdStr string, token string) *model.UserInfoResponse {
	userId, err := strconv.ParseInt(userIdStr, 10, 64)
	if err != nil {
		return &model.UserInfoResponse{
			StatusCode: -1, StatusMsg: err.Error(),
		}
	}

	if length := len(token); length <= 0 || length > 64 {
		return &model.UserInfoResponse{
			StatusCode: -1, StatusMsg: "token length out of range(0,64]",
		}
	}

	user, err := service.GetUserInfo(userId, token)
	if err != nil {
		return &model.UserInfoResponse{
			StatusCode: -1, StatusMsg: err.Error(),
		}
	}

	return &model.UserInfoResponse{
		StatusCode: 0,
		StatusMsg:  "get user info success",
		User:       user,
	}
}

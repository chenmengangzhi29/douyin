package handlers

import (
	"net/http"

	"github.com/chenmengangzhi29/douyin/pkg/errno"
	"github.com/gin-gonic/gin"
)

type UserLoginParam struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

type Response struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg"`
}

// SendResponse pack response
func SendResponse(c *gin.Context, err error) {
	Err := errno.ConvertErr(err)
	c.JSON(http.StatusOK, Response{
		StatusCode: Err.ErrCode,
		StatusMsg:  Err.ErrMsg,
	})
}

type FeedResponse struct {
	StatusCode int32       `json:"status_code"`
	StatusMsg  string      `json:"status_msg"`
	VideoList  interface{} `json:"video_list,omitempty"`
	NextTime   int64       `json:"next_time,omitempty"`
}

// SendFeedResponse pack feed response
func SendFeedResponse(c *gin.Context, err error, videoList interface{}, nextTime int64) {
	Err := errno.ConvertErr(err)
	c.JSON(http.StatusOK, FeedResponse{
		StatusCode: Err.ErrCode,
		StatusMsg:  Err.ErrMsg,
		VideoList:  videoList,
		NextTime:   nextTime,
	})
}

type VideoListResponse struct {
	StatusCode int32       `json:"status_code"`
	StatusMsg  string      `json:"status_msg"`
	VideoList  interface{} `json:"video_list,omitempty"`
}

// SendVideoListResponse pack video list response
func SendVideoListResponse(c *gin.Context, err error, videoList interface{}) {
	Err := errno.ConvertErr(err)
	c.JSON(http.StatusOK, VideoListResponse{
		StatusCode: Err.ErrCode,
		StatusMsg:  Err.ErrMsg,
		VideoList:  videoList,
	})
}

type UserResponse struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg"`
	UserId     int64  `json:"user_id,omitempty"`
	Token      string `json:"token,omitempty"`
}

// SendUserResponse pack user response
func SendUserResponse(c *gin.Context, err error, userId int64, token string) {
	Err := errno.ConvertErr(err)
	c.JSON(http.StatusOK, UserResponse{
		StatusCode: Err.ErrCode,
		StatusMsg:  Err.ErrMsg,
		UserId:     userId,
		Token:      token,
	})
}

type UserInfoResponse struct {
	StatusCode int32       `json:"status_code"`
	StatusMsg  string      `json:"status_msg"`
	User       interface{} `json:"user"`
}

// SendUserInfoResponse pack user info response
func SendUserInfoResponse(c *gin.Context, err error, user interface{}) {
	Err := errno.ConvertErr(err)
	c.JSON(http.StatusOK, UserInfoResponse{
		StatusCode: Err.ErrCode,
		StatusMsg:  Err.ErrMsg,
		User:       user,
	})
}

type FavoriteListResponse struct {
	StatusCode int32       `json:"status_code"`
	StatusMsg  string      `json:"status_msg"`
	VideoList  interface{} `json:"video_list,omitempty"`
}

// SendFavoriteListResponse pack favorite list response
func SendFavoriteListResponse(c *gin.Context, err error, videoList interface{}) {
	Err := errno.ConvertErr(err)
	c.JSON(http.StatusOK, FavoriteListResponse{
		StatusCode: Err.ErrCode,
		StatusMsg:  Err.ErrMsg,
		VideoList:  videoList,
	})
}

type CommentActionResponse struct {
	StatusCode int32       `json:"status_code"`
	StatusMsg  string      `json:"status_msg"`
	Comment    interface{} `json:"comment,omitempty"`
}

// SendCommentActionResponse pack comment action response
func SendCommentActionResponse(c *gin.Context, err error, comment interface{}) {
	Err := errno.ConvertErr(err)
	c.JSON(http.StatusOK, CommentActionResponse{
		StatusCode: Err.ErrCode,
		StatusMsg:  Err.ErrMsg,
		Comment:    comment,
	})
}

type CommentListResponse struct {
	StatusCode  int32       `json:"status_code"`
	StatusMsg   string      `json:"status_msg"`
	CommentList interface{} `json:"comment_list,omitempty"`
}

// SendCommentListResponse pack comment list response
func SendCommentListResponse(c *gin.Context, err error, commentList interface{}) {
	Err := errno.ConvertErr(err)
	c.JSON(http.StatusOK, CommentListResponse{
		StatusCode:  Err.ErrCode,
		StatusMsg:   Err.ErrMsg,
		CommentList: commentList,
	})
}

type RelationListResponse struct {
	StatusCode int32       `json:"status_code"`
	StatusMsg  string      `json:"status_msg,omitempty"`
	UserList   interface{} `json:"user_list,omitempty"`
}

// SendRelationListResponse pack relation list response
func SendRelationListResponse(c *gin.Context, err error, userList interface{}) {
	Err := errno.ConvertErr(err)
	c.JSON(http.StatusOK, RelationListResponse{
		StatusCode: Err.ErrCode,
		StatusMsg:  Err.ErrMsg,
		UserList:   userList,
	})
}

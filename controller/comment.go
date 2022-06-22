package controller

/*
because comment_id is the foreign key of each video，so we don't need to add video_id in comment table specially.
*/

import (
	"douyin/model"
	"net/http"
	"reflect"
	"strconv"
	"sync/atomic"

	"github.com/gin-gonic/gin"
)

type CommentListResponse struct {
	Response
	CommentList []Comment `json:"comment_list,omitempty"`
}

type CommentRaw struct {
	Id         int64  `gorm:"column:id"`
	UserId     int64  `gorm:"column:user_id"`
	VideoId    int64  `gorm:"column:video_id"`
	Contents   string `gorm:"column:contents"`
	CreateDate string `gorm:"column:create_date"`
}

// CommentAction no practical effect, just check if token is valid
func CommentAction(c *gin.Context) {
	//从客户端获取相应的数据
	userId, _ := strconv.ParseInt(c.Query("user_id"), 0, len(c.Query("user_id")))
	token := c.Query("token")
	actionType := c.Query("action_type")
	commentText := c.Query("content")
	videoId, _ := strconv.ParseInt(c.Query("video_id"), 0, len(c.Query("video_id")))
	commentId, _ := strconv.ParseInt(c.Query("comment_id"), 0, len(c.Query("comment_id")))
	commentCount, _ := strconv.ParseInt(c.Query("comment_count"), 0, len(c.Query("comment_count")))

	//用户权鉴
	if _, exist := usersLoginInfo[token]; exist {
		c.JSON(http.StatusOK, Response{StatusCode: 0})
	} else {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
		return
	}

	//如果操作数为1就将新增的这条评论返回到页面,并存入数据库
	if reflect.DeepEqual(actionType, 1) {

		//创建一条新的评论
		newComment := &CommentRaw{
			Id:         commentId,
			UserId:     userId,
			VideoId:    videoId,
			Contents:   commentText,
			CreateDate: model.GetDate(),
		}
		//将对应视频的评论总数加1
		atomic.AddInt64(&commentCount, 1)

		//打开数据库并传入数据
		db := model.DB
		db.Create(&newComment).Where("video_id = ?", videoId)

		//返回新增的评论到页面上
		c.JSON(http.StatusOK, newComment)

	}
	//如果操作数为2，就删除数据库中对应的评论
	if reflect.DeepEqual(actionType, 2) {

		db := model.DB
		//根据comment_id来删除数据库中对应的评论。
		db.Where("comment_id", commentId).Delete(&Comment{})
		//并将对应视频的评论总数减1
		atomic.AddInt64(&commentCount, -1)

	}
}

// CommentList all videos have same demo comment list
func CommentList(c *gin.Context) {
	token := c.Query("token")
	videoId := c.Query("video_id")
	//进行用户权鉴
	if _, exist := usersLoginInfo[token]; exist {
		c.JSON(http.StatusOK, Response{StatusCode: 0})
	} else {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
		return
	}

	//返回对应的数据库中的评论列表
	db := model.DB
	var CommentList = make([]Comment, 0)
	db.Where("video_id = ?", videoId).Find(&CommentList)
	c.JSON(http.StatusOK, CommentListResponse{
		Response:    Response{StatusCode: 0},
		CommentList: CommentList,
	})
}

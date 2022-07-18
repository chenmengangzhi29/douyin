package handlers

import (
	"bytes"
	"context"
	"io"
	"strconv"

	"github.com/chenmengangzhi29/douyin/cmd/api/rpc"
	"github.com/chenmengangzhi29/douyin/kitex_gen/publish"
	"github.com/chenmengangzhi29/douyin/pkg/errno"
	"github.com/gin-gonic/gin"
)

//Publish upload video datas
func Publish(c *gin.Context) {
	title := c.PostForm("title")
	token := c.PostForm("token")

	data, _, err := c.Request.FormFile("data")
	if err != nil {
		SendResponse(c, errno.ConvertErr(err))
		return
	}
	defer data.Close()

	if length := len(title); length <= 0 || length > 128 {
		SendResponse(c, errno.ParamErr)
		return
	}

	//处理视频数据
	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, data); err != nil {
		SendResponse(c, errno.VideoDataCopyErr)
		return
	}
	video := buf.Bytes()

	err = rpc.PublishVideoData(context.Background(), &publish.PublishActionRequest{
		Token: token,
		Title: title, Data: video,
	})
	if err != nil {
		SendResponse(c, errno.ConvertErr(err))
		return
	}

	SendResponse(c, errno.Success)
}

//PublishList get publish list
func PublishList(c *gin.Context) {
	token := c.Query("token")
	userIdStr := c.Query("user_id")

	userId, err := strconv.ParseInt(userIdStr, 10, 64)
	if err != nil {
		SendResponse(c, errno.ParamParseErr)
	}

	videoList, err := rpc.QueryVideoList(context.Background(), &publish.PublishListRequest{
		Token:  token,
		UserId: userId,
	})
	if err != nil {
		SendResponse(c, errno.ConvertErr(err))
		return
	}

	SendVideoListResponse(c, errno.Success, videoList)
}

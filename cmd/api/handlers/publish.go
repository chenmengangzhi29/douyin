package handlers

import (
	"bytes"
	"context"
	"io"
	"strconv"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/chenmengangzhi29/douyin/cmd/api/rpc"
	"github.com/chenmengangzhi29/douyin/kitex_gen/publish"
	"github.com/chenmengangzhi29/douyin/pkg/constants"
	"github.com/chenmengangzhi29/douyin/pkg/errno"
	"github.com/gin-gonic/gin"
)

//Publish upload video datas
func Publish(c *gin.Context) {
	title := c.PostForm("title")
	claims := jwt.ExtractClaims(c)
	userId := int64(claims[constants.IdentiryKey].(float64))

	data, _, err := c.Request.FormFile("data")
	if err != nil {
		SendResponse(c, errno.ConvertErr(err), nil)
		return
	}
	defer data.Close()

	if data == nil {
		SendResponse(c, errno.VideoDataGetErr, nil)
		return
	}

	if length := len(title); length <= 0 || length > 128 {
		SendResponse(c, errno.ParamErr, nil)
		return
	}

	//处理视频数据
	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, data); err != nil {
		SendResponse(c, errno.VideoDataCopyErr, nil)
		return
	}
	video := buf.Bytes()

	err = rpc.PublishVideoData(context.Background(), &publish.PublishActionRequest{
		UserId: userId,
		Title:  title, Data: video,
	})
	if err != nil {
		SendResponse(c, errno.ConvertErr(err), nil)
		return
	}

	SendResponse(c, errno.Success, nil)
}

//PublishList get publish list
func PublishList(c *gin.Context) {
	token := c.Query("token")
	userIdStr := c.Query("user_id")

	userId, err := strconv.ParseInt(userIdStr, 10, 64)
	if err != nil {
		SendResponse(c, errno.ParamParseErr, nil)
	}

	videoList, err := rpc.QueryVideoList(context.Background(), &publish.PublishListRequest{
		Token:  token,
		userId: userId,
	})
	if err != nil {
		SendResponse(c, errno.ConvertErr(err), nil)
		return
	}

	SendResponse(c, errno.Success, map[string]interface{}{constants.VideoList: videoList})
}

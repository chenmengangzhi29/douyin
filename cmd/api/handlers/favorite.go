package handlers

import (
	"context"
	"strconv"

	"github.com/chenmengangzhi29/douyin/cmd/api/rpc"
	"github.com/chenmengangzhi29/douyin/kitex_gen/favorite"
	"github.com/chenmengangzhi29/douyin/pkg/errno"
	"github.com/gin-gonic/gin"
)

// FavoriteAction implement like and unlike operations
func FavoriteAction(c *gin.Context) {
	token := c.Query("token")
	videoIdStr := c.Query("video_id")
	actionTypeStr := c.Query("action_type")

	videoId, err := strconv.ParseInt(videoIdStr, 10, 64)
	if err != nil {
		SendResponse(c, errno.ParamParseErr)
		return
	}

	actionType, err := strconv.ParseInt(actionTypeStr, 10, 64)
	if err != nil {
		SendResponse(c, errno.ParamParseErr)
		return
	}

	err = rpc.FavoriteAction(context.Background(), &favorite.FavoriteActionRequest{
		Token: token, VideoId: videoId, ActionType: int32(actionType),
	})
	if err != nil {
		SendResponse(c, errno.ConvertErr(err))
		return
	}

	SendResponse(c, errno.Success)
}

//FavoriteList get favorite list info
func FavoriteList(c *gin.Context) {
	userIdStr := c.Query("user_id")
	token := c.Query("token")

	userId, err := strconv.ParseInt(userIdStr, 10, 64)
	if err != nil {
		SendResponse(c, errno.ParamParseErr)
		return
	}

	videoList, err := rpc.FavoriteList(context.Background(), &favorite.FavoriteListRequest{
		Token: token, UserId: userId,
	})
	if err != nil {
		SendResponse(c, errno.ConvertErr(err))
		return
	}

	SendFavoriteListResponse(c, errno.Success, videoList)
}

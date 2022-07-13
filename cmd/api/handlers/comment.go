package handlers

import (
	"context"
	"strconv"
	"unicode/utf8"

	"github.com/chenmengangzhi29/douyin/cmd/api/rpc"
	"github.com/chenmengangzhi29/douyin/kitex_gen/comment"
	"github.com/chenmengangzhi29/douyin/pkg/constants"
	"github.com/chenmengangzhi29/douyin/pkg/errno"
	"github.com/gin-gonic/gin"
)

// CommentAction implement adding and deleting comments
func CommentAction(c *gin.Context) {
	token := c.Query("token")
	videoIdStr := c.Query("video_id")
	actionTypeStr := c.Query("action_type")

	videoId, err := strconv.ParseInt(videoIdStr, 10, 64)
	if err != nil {
		SendResponse(c, errno.ParamParseErr, nil)
		return
	}

	actionType, err := strconv.ParseInt(actionTypeStr, 10, 64)
	if err != nil {
		SendResponse(c, errno.ParamParseErr, nil)
		return
	}

	if actionType == 1 {
		commentText := c.Query("comment_text")

		if len := utf8.RuneCountInString(commentText); len > 20 {
			SendResponse(c, errno.CommentTextErr, nil)
			return
		}

		req := &comment.CreateCommentRequest{Token: token, VideoId: videoId, CommentText: commentText}
		comment, err := rpc.CreateComment(context.Background(), req)
		if err != nil {
			SendResponse(c, errno.ConvertErr(err), nil)
			return
		}
		SendResponse(c, errno.Success, map[string]interface{}{constants.Comment: comment})

	} else if actionType == 2 {
		commentIdStr := c.Query("comment_id")

		commentId, err := strconv.ParseInt(commentIdStr, 10, 64)
		if err != nil {
			SendResponse(c, errno.ParamParseErr, nil)
		}

		req := &comment.DeleteCommentRequest{Token: token, VideoId: videoId, CommentId: commentId}
		comment, err := rpc.DeleteComment(context.Background(), req)
		if err != nil {
			SendResponse(c, errno.ConvertErr(err), nil)
			return
		}
		SendResponse(c, errno.Success, map[string]interface{}{constants.Comment: comment})

	} else {
		SendResponse(c, errno.ParamErr, nil)
	}
}

//CommentList get comment list info
func CommentList(c *gin.Context) {
	token := c.Query("token")
	videoIdStr := c.Query("video_id")

	videoId, err := strconv.ParseInt(videoIdStr, 10, 64)
	if err != nil {
		SendResponse(c, errno.ParamParseErr, nil)
	}

	req := &comment.CommentListRequest{Token: token, VideoId: videoId}
	commentList, err := rpc.CommentList(context.Background(), req)
	if err != nil {
		SendResponse(c, errno.ConvertErr(err), nil)
		return
	}

	SendResponse(c, errno.Success, map[string]interface{}{constants.CommentList: commentList})
}
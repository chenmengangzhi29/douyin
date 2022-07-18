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
		SendResponse(c, errno.ParamParseErr)
		return
	}

	actionType, err := strconv.ParseInt(actionTypeStr, 10, 64)
	if err != nil {
		SendResponse(c, errno.ParamParseErr)
		return
	}

	if actionType == constants.AddComment {
		commentText := c.Query("comment_text")

		if len := utf8.RuneCountInString(commentText); len > 20 {
			SendResponse(c, errno.CommentTextErr)
			return
		}

		req := &comment.CreateCommentRequest{Token: token, VideoId: videoId, CommentText: commentText}
		comment, err := rpc.CreateComment(context.Background(), req)
		if err != nil {
			SendResponse(c, errno.ConvertErr(err))
			return
		}
		SendCommentActionResponse(c, errno.Success, comment)

	} else if actionType == constants.DelComment {
		commentIdStr := c.Query("comment_id")

		commentId, err := strconv.ParseInt(commentIdStr, 10, 64)
		if err != nil {
			SendResponse(c, errno.ParamParseErr)
		}

		req := &comment.DeleteCommentRequest{Token: token, VideoId: videoId, CommentId: commentId}
		comment, err := rpc.DeleteComment(context.Background(), req)
		if err != nil {
			SendResponse(c, errno.ConvertErr(err))
			return
		}
		SendCommentActionResponse(c, errno.Success, comment)

	} else {
		SendResponse(c, errno.ParamErr)
	}
}

//CommentList get comment list info
func CommentList(c *gin.Context) {
	token := c.Query("token")
	videoIdStr := c.Query("video_id")

	videoId, err := strconv.ParseInt(videoIdStr, 10, 64)
	if err != nil {
		SendResponse(c, errno.ParamParseErr)
	}

	req := &comment.CommentListRequest{Token: token, VideoId: videoId}
	commentList, err := rpc.CommentList(context.Background(), req)
	if err != nil {
		SendResponse(c, errno.ConvertErr(err))
		return
	}

	SendCommentListResponse(c, errno.Success, commentList)
}

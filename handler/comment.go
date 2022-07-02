package handler

import (
	"douyin/model"
	"douyin/service"
	"strconv"
	"unicode/utf8"
)

//处理创建评论的传入参数
func CreateComment(token string, videoIdStr string, commentText string) *model.CommentActionResponse {
	if len := len(token); len <= 0 || len > 64 {
		return &model.CommentActionResponse{
			StatusCode: -1, StatusMsg: "token length out of range(0,64]",
		}
	}

	videoId, err := strconv.ParseInt(videoIdStr, 10, 64)
	if err != nil {
		return &model.CommentActionResponse{
			StatusCode: -1, StatusMsg: err.Error(),
		}
	}

	if len := utf8.RuneCountInString(commentText); len > 20 {
		return &model.CommentActionResponse{
			StatusCode: -1, StatusMsg: "comment text too long (must be <= 20)",
		}
	}

	comment, err := service.CreateCommentData(token, videoId, commentText)
	if err != nil {
		return &model.CommentActionResponse{
			StatusCode: -1, StatusMsg: err.Error(),
		}
	}

	return &model.CommentActionResponse{
		StatusCode: 0,
		StatusMsg:  "create comment success",
		Comment:    comment,
	}
}

//处理删除评论的传入参数
func DeleteComment(token string, videoIdStr string, commentIdStr string) *model.CommentActionResponse {
	if len := len(token); len <= 0 || len > 64 {
		return &model.CommentActionResponse{
			StatusCode: -1, StatusMsg: "token length out of range(0,64]",
		}
	}

	videoId, err := strconv.ParseInt(videoIdStr, 10, 64)
	if err != nil {
		return &model.CommentActionResponse{
			StatusCode: -1, StatusMsg: err.Error(),
		}
	}

	commentId, err := strconv.ParseInt(commentIdStr, 10, 64)
	if err != nil {
		return &model.CommentActionResponse{
			StatusCode: -1, StatusMsg: err.Error(),
		}
	}

	comment, err := service.DeleteCommentData(token, videoId, commentId)
	if err != nil {
		return &model.CommentActionResponse{
			StatusCode: -1, StatusMsg: err.Error(),
		}
	}

	return &model.CommentActionResponse{
		StatusCode: 0,
		StatusMsg:  "delete comment success",
		Comment:    comment,
	}

}

//处理评论列表的传入参数

func ShowCommentList(token string, videoIdStr string) *model.CommentListResponse {
	if len := len(token); len <= 0 || len > 64 {
		return &model.CommentListResponse{
			StatusCode: -1, StatusMsg: "token length out of range(0,64]",
		}
	}

	videoId, err := strconv.ParseInt(videoIdStr, 10, 64)
	if err != nil {
		return &model.CommentListResponse{
			StatusCode: -1, StatusMsg: err.Error(),
		}
	}

	commentList, err := service.CommentListData(token, videoId)
	if err != nil {
		return &model.CommentListResponse{
			StatusCode: -1, StatusMsg: err.Error(),
		}
	}

	return &model.CommentListResponse{
		StatusCode:  0,
		StatusMsg:   "comment list success",
		CommentList: commentList,
	}
}

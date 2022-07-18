package main

import (
	"context"
	"unicode/utf8"

	"github.com/chenmengangzhi29/douyin/cmd/comment/service"
	"github.com/chenmengangzhi29/douyin/dal/pack"
	"github.com/chenmengangzhi29/douyin/kitex_gen/comment"
	"github.com/chenmengangzhi29/douyin/pkg/errno"
)

// CommentServiceImpl implements the last service interface defined in the IDL.
type CommentServiceImpl struct{}

// CreateComment implements the CommentServiceImpl interface.
func (s *CommentServiceImpl) CreateComment(ctx context.Context, req *comment.CreateCommentRequest) (resp *comment.CreateCommentResponse, err error) {
	resp = new(comment.CreateCommentResponse)

	if len(req.Token) == 0 || req.VideoId == 0 || utf8.RuneCountInString(req.CommentText) > 20 {
		resp.BaseResp = pack.BuilCommentBaseResp(errno.ParamErr)
		return resp, nil
	}

	comment, err := service.NewCreateCommentService(ctx).CreateComment(req)
	if err != nil {
		resp.BaseResp = pack.BuilCommentBaseResp(err)
		return resp, nil
	}

	resp.BaseResp = pack.BuilCommentBaseResp(errno.Success)
	resp.Comment = comment
	return resp, nil
}

// DeleteComment implements the CommentServiceImpl interface.
func (s *CommentServiceImpl) DeleteComment(ctx context.Context, req *comment.DeleteCommentRequest) (resp *comment.DeleteCommentResponse, err error) {
	resp = new(comment.DeleteCommentResponse)

	if len(req.Token) == 0 || req.VideoId == 0 || req.CommentId == 0 {
		resp.BaseResp = pack.BuilCommentBaseResp(errno.ParamErr)
		return resp, nil
	}

	comment, err := service.NewDeleteCommentService(ctx).DeleteComment(req)
	if err != nil {
		resp.BaseResp = pack.BuilCommentBaseResp(err)
		return resp, nil
	}

	resp.BaseResp = pack.BuilCommentBaseResp(errno.Success)
	resp.Comment = comment
	return resp, nil
}

// CommentList implements the CommentServiceImpl interface.
func (s *CommentServiceImpl) CommentList(ctx context.Context, req *comment.CommentListRequest) (resp *comment.CommentListResponse, err error) {
	resp = new(comment.CommentListResponse)

	if req.VideoId == 0 {
		resp.BaseResp = pack.BuilCommentBaseResp(errno.ParamErr)
		return resp, nil
	}

	commentList, err := service.NewCommentListService(ctx).CommentList(req)
	if err != nil {
		resp.BaseResp = pack.BuilCommentBaseResp(err)
		return resp, nil
	}

	resp.BaseResp = pack.BuilCommentBaseResp(errno.Success)
	resp.CommentList = commentList
	return resp, nil
}

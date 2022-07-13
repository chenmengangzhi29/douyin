package service

import (
	"context"
	"errors"
	"sync"

	"github.com/chenmengangzhi29/douyin/dal/db"
	"github.com/chenmengangzhi29/douyin/dal/pack"
	"github.com/chenmengangzhi29/douyin/kitex_gen/comment"
	"github.com/chenmengangzhi29/douyin/pkg/constants"
	"github.com/chenmengangzhi29/douyin/pkg/jwt"
)

type DeleteCommentService struct {
	ctx context.Context
}

// NewDeleteCommentService new DeleteCommentService
func NewDeleteCommentService(ctx context.Context) *DeleteCommentService {
	return &DeleteCommentService{ctx: ctx}
}

// DeleteComment delete comment
func (s *DeleteCommentService) DeleteComment(req *comment.DeleteCommentRequest) (*comment.Comment, error) {
	Jwt := jwt.NewJWT([]byte(constants.SecretKey))
	currentId, err := Jwt.CheckToken(req.Token)
	if err != nil {
		return nil, err
	}

	videos, err := db.QueryVideoByVideoIds(s.ctx, []int64{req.VideoId})
	if err != nil {
		return nil, err
	}
	if len(videos) == 0 {
		return nil, errors.New("videoId not exist")
	}
	comments, err := db.QueryCommentByCommentIds(s.ctx, []int64{req.CommentId})
	if err != nil {
		return nil, err
	}
	if len(comments) == 0 {
		return nil, errors.New("commentId not exist")
	}

	var wg sync.WaitGroup
	wg.Add(2)
	var commentRaw *db.CommentRaw
	var userRaw *db.UserRaw
	var commentErr, userErr error
	//删除评论记录并减少视频评论数
	go func() {
		defer wg.Done()
		commentRaw, err = db.DeleteComment(s.ctx, req.CommentId)
		if err != nil {
			commentErr = err
			return
		}
	}()
	//获取用户信息
	go func() {
		defer wg.Done()
		users, err := db.QueryUserByIds(s.ctx, []int64{currentId})
		if err != nil {
			userErr = err
			return
		}
		userRaw = users[0]
	}()
	wg.Wait()
	if commentErr != nil {
		return nil, commentErr
	}
	if userErr != nil {
		return nil, userErr
	}

	comment := pack.CommentInfo(commentRaw, userRaw)
	return comment, nil
}

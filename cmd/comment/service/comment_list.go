package service

import (
	"context"
	"errors"

	"github.com/chenmengangzhi29/douyin/dal/db"
	"github.com/chenmengangzhi29/douyin/dal/pack"
	"github.com/chenmengangzhi29/douyin/kitex_gen/comment"
	"github.com/chenmengangzhi29/douyin/pkg/constants"
	"github.com/chenmengangzhi29/douyin/pkg/jwt"
)

type CommentListService struct {
	ctx context.Context
}

// NewCommentListService new CommentListService
func NewCommentListService(ctx context.Context) *CommentListService {
	return &CommentListService{ctx: ctx}
}

// CommentList get comment list info
func (s *CommentListService) CommentList(req *comment.CommentListRequest) ([]*comment.Comment, error) {
	//获取当前用户id号
	Jwt := jwt.NewJWT([]byte(constants.SecretKey))
	currentId, _ := Jwt.CheckToken(req.Token)

	//验证视频Id是否存在
	videos, err := db.QueryVideoByVideoIds(s.ctx, []int64{req.VideoId})
	if err != nil {
		return nil, err
	}
	if len(videos) == 0 {
		return nil, errors.New("video not exist")
	}

	//获取一系列评论信息
	comments, err := db.QueryCommentByVideoId(s.ctx, req.VideoId)
	if err != nil {
		return nil, err
	}

	//获取评论信息的用户id
	userIds := make([]int64, 0)
	for _, comment := range comments {
		userIds = append(userIds, comment.UserId)
	}

	//获取一系列用户信息
	users, err := db.QueryUserByIds(s.ctx, userIds)
	if err != nil {
		return nil, err
	}
	userMap := make(map[int64]*db.UserRaw)
	for _, user := range users {
		userMap[int64(user.ID)] = user
	}

	var relationMap map[int64]*db.RelationRaw
	if currentId == -1 {
		relationMap = nil
	} else {
		//获取一系列关注信息
		relationMap, err = db.QueryRelationByIds(s.ctx, currentId, userIds)
		if err != nil {
			return nil, err
		}
	}

	commentList := pack.CommentList(currentId, comments, userMap, relationMap)
	return commentList, nil
}

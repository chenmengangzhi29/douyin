package service

import (
	"context"
	"errors"

	"github.com/chenmengangzhi29/douyin/dal/db"
	"github.com/chenmengangzhi29/douyin/kitex_gen/relation"
	"github.com/chenmengangzhi29/douyin/pkg/constants"
	"github.com/chenmengangzhi29/douyin/pkg/jwt"
)

type RelationActionService struct {
	ctx context.Context
}

// NewRelationActionService new RelationActionService
func NewRelationActionService(ctx context.Context) *RelationActionService {
	return &RelationActionService{ctx: ctx}
}

// RelationAction implement follow and unfollow action
//如果actionType等于1，表示当前用户关注其他用户，
//当前用户的关注总数增加，其他用户的粉丝总数增加，
//新建一条关注记录
//
//如果actionType等于2，表示当前用户取消关注其他用户
//当前用户的关注总数减少，其他用户的粉丝总数减少，
//删除该关注记录
func (s *RelationActionService) RelationAction(req *relation.RelationActionRequest) error {
	Jwt := jwt.NewJWT([]byte(constants.SecretKey))
	currentId, err := Jwt.CheckToken(req.Token)
	if err != nil {
		return err
	}

	users, err := db.QueryUserByIds(s.ctx, []int64{req.ToUserId})
	if err != nil {
		return err
	}
	if len(users) == 0 {
		return errors.New("toUserId not exist")
	}

	if req.ActionType == constants.Follow {
		err := db.Create(s.ctx, currentId, req.ToUserId)
		if err != nil {
			return err
		}
		return nil
	}
	if req.ActionType == constants.UnFollow {
		err := db.Delete(s.ctx, currentId, req.ToUserId)
		if err != nil {
			return err
		}
		return nil
	}
	return errors.New("ActionType Err")
}

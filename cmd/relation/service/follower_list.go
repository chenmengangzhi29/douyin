package service

import (
	"context"
	"errors"

	"github.com/chenmengangzhi29/douyin/dal/db"
	"github.com/chenmengangzhi29/douyin/dal/pack"
	"github.com/chenmengangzhi29/douyin/kitex_gen/relation"
	"github.com/chenmengangzhi29/douyin/pkg/constants"
	"github.com/chenmengangzhi29/douyin/pkg/jwt"
)

type FollowerListService struct {
	ctx context.Context
}

// NewFollowerListService new FollowerListService
func NewFollowerListService(ctx context.Context) *FollowerListService {
	return &FollowerListService{ctx: ctx}
}

// FollowerList get user follower list info
func (s *FollowerListService) FollowerList(req *relation.FollowerListRequest) ([]*relation.User, error) {
	Jwt := jwt.NewJWT([]byte(constants.SecretKey))
	currentId, _ := Jwt.CheckToken(req.Token)

	user, err := db.QueryUserByIds(s.ctx, []int64{req.UserId})
	if err != nil {
		return nil, err
	}
	if len(user) == 0 {
		return nil, errors.New("userId not exist")
	}

	//查询目标用户的被关注记录
	relations, err := db.QueryFollowerById(s.ctx, req.UserId)
	if err != nil {
		return nil, err
	}

	//获取这些记录的关注方id
	userIds := make([]int64, 0)
	for _, relation := range relations {
		userIds = append(userIds, relation.UserId)
	}

	//获取关注方的信息
	users, err := db.QueryUserByIds(s.ctx, userIds)
	if err != nil {
		return nil, err
	}

	var relationMap map[int64]*db.RelationRaw
	if currentId == -1 {
		relationMap = nil
	} else {
		//获取当前用户与关注方的关注记录
		relationMap, err = db.QueryRelationByIds(s.ctx, currentId, userIds)
		if err != nil {
			return nil, err
		}
	}

	userList := pack.UserList(currentId, users, relationMap)
	return userList, nil
}

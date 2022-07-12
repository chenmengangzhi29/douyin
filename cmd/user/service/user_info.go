package service

import (
	"context"
	"errors"

	"github.com/chenmengangzhi29/douyin/dal/db"
	"github.com/chenmengangzhi29/douyin/dal/pack"
	"github.com/chenmengangzhi29/douyin/kitex_gen/user"
	"github.com/chenmengangzhi29/douyin/pkg/jwt"
)

type UserInfoService struct {
	ctx context.Context
}

// NewUserInfoService new UserInfoService
func NewUserInfoService(ctx context.Context) *UserInfoService {
	return &UserInfoService{
		ctx: ctx,
	}
}

func (s *UserInfoService) UserInfo(req *user.UserInfoRequest) (*user.User, error) {
	currentId, err := s.checkToken(req.Token)
	if err != nil {
		return nil, err
	}

	userIds := []int64{req.UserId}
	users, err := db.QueryUserByIds(s.ctx, userIds)
	if err != nil {
		return nil, err
	}
	if len(users) == 0 {
		return nil, errors.New("user not exist")
	}
	user := users[0]

	relationMap, err := db.QueryRelationByIds(s.ctx, currentId, userIds)
	var isFollow bool
	_, ok := relationMap[req.UserId]
	if ok {
		isFollow = true
	} else {
		isFollow = false
	}

	userInfo := pack.UserInfo(user, isFollow)
	return userInfo, nil
}

//checkToken get userId by token
func (s *UserInfoService) checkToken(token string) (int64, error) {
	if token == "" {
		return -1, nil
	}
	var Jwt *jwt.JWT
	claim, err := Jwt.ParseToken(token)
	if err != nil {
		return 0, jwt.ErrTokenInvalid
	}
	return claim.Id, nil
}

package main

import (
	"context"

	"github.com/chenmengangzhi29/douyin/cmd/user/service"
	"github.com/chenmengangzhi29/douyin/dal/pack"
	"github.com/chenmengangzhi29/douyin/kitex_gen/user"
	"github.com/chenmengangzhi29/douyin/pkg/errno"
	"github.com/chenmengangzhi29/douyin/pkg/jwt"
)

// UserServiceImpl implements the last service interface defined in the IDL.
type UserServiceImpl struct{}

// CheckUser implements the UserServiceImpl interface.
func (s *UserServiceImpl) CheckUser(ctx context.Context, req *user.CheckUserRequest) (resp *user.CheckUserResponse, err error) {
	resp = new(user.CheckUserResponse)

	if len(req.Username) == 0 || len(req.Password) == 0 {
		resp.BaseResp = pack.BuildUserBaseResp(errno.ParamErr)
		return resp, nil
	}

	uid, err := service.NewCheckUserService(ctx).CheckUser(req)
	if err != nil {
		resp.BaseResp = pack.BuildUserBaseResp(err)
		return resp, nil
	}
	resp.UserId = uid
	resp.BaseResp = pack.BuildUserBaseResp(errno.Success)
	return resp, nil
}

// RegisterUser implements the UserServiceImpl interface.
func (s *UserServiceImpl) RegisterUser(ctx context.Context, req *user.RegisterUserRequest) (resp *user.RegisterUserResponse, err error) {
	resp = new(user.RegisterUserResponse)

	if len(req.Username) == 0 || len(req.Password) == 0 {
		resp.BaseResp = pack.BuildUserBaseResp(errno.ParamErr)
	}

	userId, err := service.NewRegisterUserService(ctx).RegisterUser(req)
	if err != nil {
		resp.BaseResp = pack.BuildUserBaseResp(err)
		return
	}

	token, err := Jwt.CreateToken(jwt.CustomClaims{
		Id: int64(userId),
	})
	if err != nil {
		resp.BaseResp = pack.BuildUserBaseResp(err)
		return
	}

	resp.BaseResp = pack.BuildUserBaseResp(errno.Success)
	resp.UserId = userId
	resp.Token = token
	return resp, nil
}

// UserInfo implements the UserServiceImpl interface.
func (s *UserServiceImpl) UserInfo(ctx context.Context, req *user.UserInfoRequest) (resp *user.UserInfoResponse, err error) {
	resp = new(user.UserInfoResponse)

	if req.UserId == 0 {
		resp.BaseResp = pack.BuildUserBaseResp(errno.ParamErr)
	}

	user, err := service.NewUserInfoService(ctx).UserInfo(req)
	if err != nil {
		resp.BaseResp = pack.BuildUserBaseResp(err)
		return
	}

	resp.BaseResp = pack.BuildUserBaseResp(errno.Success)
	resp.User = user
	return resp, nil
}

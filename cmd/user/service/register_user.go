package service

import (
	"context"
	"crypto/md5"
	"fmt"
	"io"

	"github.com/chenmengangzhi29/douyin/dal/db"
	"github.com/chenmengangzhi29/douyin/kitex_gen/user"
	"github.com/chenmengangzhi29/douyin/pkg/errno"
)

type RegisterUserService struct {
	ctx context.Context
}

// NewRegisterUserService new RegisterUserService
func NewRegisterUserService(ctx context.Context) *RegisterUserService {
	return &RegisterUserService{
		ctx: ctx,
	}
}

// RegisterUser register user info
func (s *RegisterUserService) RegisterUser(req *user.RegisterUserRequest) (int64, error) {
	users, err := db.QueryUserByName(s.ctx, req.Username)
	if err != nil {
		return 0, err
	}
	if len(users) != 0 {
		return 0, errno.UserAlreadyExistErr
	}

	h := md5.New()
	if _, err = io.WriteString(h, req.Password); err != nil {
		return 0, err
	}
	password := fmt.Sprintf("%x", h.Sum(nil))

	userId, err := db.UploadUserData(s.ctx, req.Username, password)
	if err != nil {
		return 0, err
	}

	return userId, nil

}

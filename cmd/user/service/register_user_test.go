package service

import (
	"bytes"
	"context"
	"io"
	"os"
	"testing"

	"github.com/chenmengangzhi29/douyin/dal/db"
	"github.com/chenmengangzhi29/douyin/kitex_gen/user"
	"github.com/chenmengangzhi29/douyin/pkg/constants"
	"github.com/chenmengangzhi29/douyin/pkg/jwt"
	"github.com/chenmengangzhi29/douyin/pkg/logger"
	"github.com/chenmengangzhi29/douyin/pkg/oss"
)

var File []byte
var Token string

func TestMain(m *testing.M) {
	err := logger.Init()
	if err != nil {
		panic(err)
	}

	Jwt := jwt.NewJWT([]byte(constants.SecretKey))
	Token, err = Jwt.CreateToken(jwt.CustomClaims{
		Id: int64(1),
	})
	if err != nil {
		logger.Errorf("create token fail, %v", err.Error())
		panic(err)
	}

	db.Init()
	oss.Init()

	path := oss.Path + "/public/girl.mp4"
	file, err := os.Open(path)
	if err != nil {
		logger.Errorf("open local file %v fail", path)
		panic(err)
	}
	defer file.Close()
	//处理视频数据
	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, file); err != nil {
		panic(err)
	}
	File = buf.Bytes()

	m.Run()
}

//测试用户注册
func TestRegisterUser(t *testing.T) {
	type args struct {
		username string
		password string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "测试注册不存在的用户",
			args: args{
				username: "xiangjiang",
				password: "xiaojiang",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userId, err := NewRegisterUserService(context.Background()).RegisterUser(&user.RegisterUserRequest{Username: tt.args.username, Password: tt.args.password})
			if (err != nil) != tt.wantErr {
				t.Errorf("RegisterUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			logger.Info(tt.name + " success")
			logger.Info(userId)
		})
	}
}

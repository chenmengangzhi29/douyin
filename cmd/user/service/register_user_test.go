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
	"github.com/chenmengangzhi29/douyin/pkg/oss"
	"github.com/cloudwego/kitex/pkg/klog"
)

var File []byte
var Token string

func TestMain(m *testing.M) {

	Jwt := jwt.NewJWT([]byte(constants.SecretKey))
	token, err := Jwt.CreateToken(jwt.CustomClaims{
		Id: int64(1),
	})
	if err != nil {
		klog.Errorf("create token fail, %v", err.Error())
		panic(err)
	}
	Token = token

	db.Init()
	oss.Init()

	path := oss.Path + "/public/girl.mp4"
	file, err := os.Open(path)
	if err != nil {
		klog.Errorf("open local file %v fail", path)
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
				username: "hhh",
				password: "hhh123",
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
			klog.Info(tt.name + " success")
			klog.Info(userId)
		})
	}
}

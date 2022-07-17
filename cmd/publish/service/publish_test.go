package service

import (
	"bytes"
	"context"
	"io"
	"os"
	"testing"

	"github.com/chenmengangzhi29/douyin/dal/db"
	"github.com/chenmengangzhi29/douyin/kitex_gen/publish"
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

func TestPublish(t *testing.T) {
	type args struct {
		token string
		data  []byte
		title string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "测试上传视频的正常操作",
			args: args{
				token: Token,
				data:  File,
				title: "girl",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NewPublishService(context.Background()).Publish(&publish.PublishActionRequest{Token: tt.args.token, Title: tt.args.title, Data: tt.args.data})
			if (err != nil) != tt.wantErr {
				t.Errorf("%v fail, %v, wantErr %v", tt.name, err, tt.wantErr)
				return
			}
			klog.Info(tt.name + " success")
		})
	}
}

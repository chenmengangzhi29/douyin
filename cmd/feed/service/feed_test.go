package service

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/chenmengangzhi29/douyin/dal/db"
	"github.com/chenmengangzhi29/douyin/kitex_gen/feed"
	"github.com/chenmengangzhi29/douyin/pkg/constants"
	"github.com/chenmengangzhi29/douyin/pkg/jwt"
	"github.com/chenmengangzhi29/douyin/pkg/oss"
	"github.com/cloudwego/kitex/pkg/klog"
)

var File *os.File
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
	File = file

	m.Run()
}

func TestFeed(t *testing.T) {
	type args struct {
		latestTime int64
		token      string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "测试视频流的默认token",
			args: args{
				latestTime: time.Now().Unix(),
				token:      "",
			},
			wantErr: false,
		},
		{
			name: "测试视频流的正常token",
			args: args{
				latestTime: time.Now().Unix(),
				token:      Token,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := NewFeedService(context.Background()).Feed(&feed.FeedRequest{LatestTime: tt.args.latestTime, Token: tt.args.token})
			if (err != nil) != tt.wantErr {
				t.Errorf("Feed() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			klog.Info(tt.name + " success")
		})
	}
}

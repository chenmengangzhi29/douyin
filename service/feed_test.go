package service

import (
	"douyin/model"
	"douyin/util"
	"fmt"
	"os"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	if err := model.Init(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if err := util.InitLogger(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	m.Run()
}

func TestQueryVideoData(t *testing.T) {
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
			name: "测试视频流的不存在token",
			args: args{
				latestTime: time.Now().Unix(),
				token:      "feedTest",
			},
			wantErr: true,
		},
		{
			name: "测试视频流的默认token",
			args: args{
				latestTime: time.Now().Unix(),
				token:      "defaultToken",
			},
			wantErr: false,
		},
		{
			name: "测试视频流的正常token",
			args: args{
				latestTime: time.Now().Unix(),
				token:      "JerryJerry123",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := QueryVideoData(tt.args.latestTime, tt.args.token)
			if (err != nil) != tt.wantErr {
				t.Errorf("QueryVideoData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

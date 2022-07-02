package service

import (
	"douyin/util/logger"
	"testing"
)

//该测试程序不包含像TestMain初始化数据库的程序，所以要与feed_test.go一样运行，因为feed_test.go包含TestMain

// func TestPublishUserVideoData(t *testing.T) {
// 	type args struct {
// 		token string
// 		data  multipart.File
// 		title string
// 	}
// 	tests := []struct {
// 		name    string
// 		args    args
// 		wantErr bool
// 	}{
// 		{
// 			name: "测试上传视频的正常操作",
// 			args: args{
// 				token: "JerryJerry123",
// 				data:  File,
// 				title: "girl",
// 			},
// 			wantErr: false,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			err := PublishUserVideoData(tt.args.token, tt.args.data, tt.args.title)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("%v fail, %v, wantErr %v", tt.name, err, tt.wantErr)
// 				return
// 			}
// 			logger.Info(tt.name + " success")
// 		})
// 	}
// }

func TestQueryUserVideoList(t *testing.T) {
	type args struct {
		token  string
		userId int64
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "测试视频列表的不正确token",
			args: args{
				token:  "unToken",
				userId: 1,
			},
			wantErr: true,
		},
		{
			name: "测试视频列表的默认token",
			args: args{
				token:  "defaultToken",
				userId: 1,
			},
			wantErr: false,
		},
		{
			name: "测试视频列表的正常token",
			args: args{
				token:  "JerryJerry123",
				userId: 2,
			},
			wantErr: false,
		},
		{
			name: "测试视频列表的不正确用户id",
			args: args{
				token:  "JerryJerry123",
				userId: 99999999,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			videoList, err := QueryUserVideoList(tt.args.token, tt.args.userId)
			if (err != nil) != tt.wantErr {
				t.Errorf("QueryUserVideoList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			logger.Info(tt.name + " success")
			logger.Info(videoList)
		})
	}
}

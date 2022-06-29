package service

import (
	"douyin/util/logger"
	"testing"
)

//该测试程序不包含像TestMain初始化数据库的程序，所以要与feed_test.go一样运行，因为feed_test.go包含TestMain

func TestFavoriteActionData(t *testing.T) {
	type args struct {
		token      string
		videoId    int64
		actionType int64
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "测试点赞操作的actionType=1",
			args: args{
				token:      "JerryJerry123",
				videoId:    2,
				actionType: 1,
			},
			wantErr: false,
		},
		{
			name: "测试点赞操作的actionType=2",
			args: args{
				token:      "JerryJerry123",
				videoId:    2,
				actionType: 2,
			},
			wantErr: false,
		},
		{
			name: "测试点赞操作的actionType=3",
			args: args{
				token:      "JerryJerry123",
				videoId:    2,
				actionType: 3,
			},
			wantErr: true,
		},
		{
			name: "测试点赞列表的不存在视频id",
			args: args{
				token:      "JerryJerry123",
				videoId:    99999999,
				actionType: 1,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := FavoriteActionData(tt.args.token, tt.args.videoId, tt.args.actionType)
			if (err != nil) != tt.wantErr {
				t.Errorf("FavoriteActionData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			logger.Info(tt.name + " success")
		})
	}
}

func TestFavoriteListData(t *testing.T) {
	type args struct {
		userId int64
		token  string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "测试点赞列表的存在用户id",
			args: args{
				token:  "JerryJerry123",
				userId: 1,
			},
			wantErr: false,
		},
		{
			name: "测试点赞列表的不存在用户id",
			args: args{
				token:  "JerryJerry123",
				userId: 99999999,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			videoList, err := FavoriteListData(tt.args.userId, tt.args.token)
			if (err != nil) != tt.wantErr {
				t.Errorf("FavoriteListData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			logger.Info(tt.name + " success")
			logger.Info(videoList)
		})
	}
}

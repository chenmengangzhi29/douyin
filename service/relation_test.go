package service

import (
	"douyin/util/logger"
	"testing"
)

//该测试程序不包含像TestMain初始化数据库的程序，所以要与feed_test.go一样运行，因为feed_test.go包含TestMain

func TestRelationActionData(t *testing.T) {
	type args struct {
		token      string
		toUserId   int64
		actionType int64
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "测试关注操作的正常操作",
			args: args{
				token:      "JerryJerry123",
				toUserId:   2,
				actionType: 1,
			},
			wantErr: false,
		},
		{
			name: "测试取消关注的正常操作",
			args: args{
				token:      "JerryJerry123",
				toUserId:   2,
				actionType: 2,
			},
			wantErr: false,
		},
		{
			name: "测试关注操作的错误用户id",
			args: args{
				token:      "JerryJerry123",
				toUserId:   99999999,
				actionType: 1,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := RelationActionData(tt.args.token, tt.args.toUserId, tt.args.actionType)
			if (err != nil) != tt.wantErr {
				t.Errorf("RelationActionData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			logger.Info(tt.name + " success")
		})
	}
}

func TestFollowListData(t *testing.T) {
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
			name: "测试关注列表的正常操作",
			args: args{
				token:  "JerryJerry123",
				userId: 2,
			},
			wantErr: false,
		},
		{
			name: "测试关注列表的错误用户id",
			args: args{
				token:  "JerryJerry123",
				userId: 99999999,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := FollowListData(tt.args.token, tt.args.userId)
			if (err != nil) != tt.wantErr {
				t.Errorf("FollowListData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			logger.Info(tt.name + " success")
		})
	}
}

func TestFollowerListData(t *testing.T) {
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
			name: "测试粉丝列表的正常操作",
			args: args{
				token:  "JerryJerry123",
				userId: 3,
			},
			wantErr: false,
		},
		{
			name: "测试粉丝列表的没有粉丝的用户",
			args: args{
				token:  "JerryJerry123",
				userId: 1,
			},
			wantErr: false,
		},
		{
			name: "测试粉丝列表的错误用户id",
			args: args{
				token:  "JerryJerry123",
				userId: 99999999,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := FollowerListData(tt.args.token, tt.args.userId)
			if (err != nil) != tt.wantErr {
				t.Errorf("FollowerListData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			logger.Info(tt.name + " success")
		})
	}
}

package service

import (
	"context"
	"testing"

	"github.com/chenmengangzhi29/douyin/kitex_gen/publish"
	"github.com/cloudwego/kitex/pkg/klog"
)

func TestPublishList(t *testing.T) {
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
			name: "测试视频列表的默认token",
			args: args{
				token:  "",
				userId: 1,
			},
			wantErr: false,
		},
		{
			name: "测试视频列表的正常token",
			args: args{
				token:  Token,
				userId: 1,
			},
			wantErr: false,
		},
		{
			name: "测试视频列表的不正确用户id",
			args: args{
				token:  Token,
				userId: 99999999,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			videoList, err := NewPublishListService(context.Background()).PublishList(&publish.PublishListRequest{UserId: tt.args.userId, Token: tt.args.token})
			if (err != nil) != tt.wantErr {
				t.Errorf("PublishList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			klog.Info(tt.name + " success")
			klog.Info(videoList)
		})
	}
}

package service

import (
	"context"
	"testing"

	"github.com/chenmengangzhi29/douyin/kitex_gen/relation"
	"github.com/cloudwego/kitex/pkg/klog"
)

func TestFollowerList(t *testing.T) {
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
				token:  Token,
				userId: 2,
			},
			wantErr: false,
		},
		{
			name: "测试粉丝列表的没有粉丝的用户",
			args: args{
				token:  Token,
				userId: 1,
			},
			wantErr: false,
		},
		{
			name: "测试粉丝列表的错误用户id",
			args: args{
				token:  Token,
				userId: 99999999,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewFollowerListService(context.Background()).FollowerList(&relation.FollowerListRequest{Token: tt.args.token, UserId: tt.args.userId})
			if (err != nil) != tt.wantErr {
				t.Errorf("FollowerList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			klog.Info(tt.name + " success")
		})
	}
}

package service

import (
	"context"
	"testing"

	"github.com/chenmengangzhi29/douyin/kitex_gen/favorite"
	"github.com/cloudwego/kitex/pkg/klog"
)

func TestFavoriteList(t *testing.T) {
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
				token:  Token,
				userId: 1,
			},
			wantErr: false,
		},
		{
			name: "测试点赞列表的不存在用户id",
			args: args{
				token:  Token,
				userId: 99999999,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			videoList, err := NewFavoriteListService(context.Background()).FavoriteList(&favorite.FavoriteListRequest{Token: tt.args.token, UserId: tt.args.userId})
			if (err != nil) != tt.wantErr {
				t.Errorf("FavoriteList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			klog.Info(tt.name + " success")
			klog.Info(videoList)
		})
	}
}

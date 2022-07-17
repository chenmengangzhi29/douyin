package service

import (
	"context"
	"testing"

	"github.com/chenmengangzhi29/douyin/kitex_gen/user"
	"github.com/cloudwego/kitex/pkg/klog"
)

//测试用户信息
func TestUserInfo(t *testing.T) {
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
			name: "测试获取存在用户的信息",
			args: args{
				userId: 1,
				token:  Token,
			},
			wantErr: false,
		},
		{
			name: "测试获取不存在用户的信息",
			args: args{
				userId: 99999999,
				token:  "JerryJerry123",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := NewUserInfoService(context.Background()).UserInfo(&user.UserInfoRequest{UserId: tt.args.userId, Token: tt.args.token})
			if (err != nil) != tt.wantErr {
				t.Errorf("UserInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			klog.Info(user)
			klog.Info(tt.name + " success")
		})
	}
}

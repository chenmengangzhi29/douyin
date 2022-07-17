package service

import (
	"context"
	"testing"

	"github.com/chenmengangzhi29/douyin/kitex_gen/user"
	"github.com/cloudwego/kitex/pkg/klog"
)

func TestCheckUser(t *testing.T) {
	type args struct {
		username string
		password string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "测试登陆存在的用户",
			args: args{
				username: "xiaohuang",
				password: "xiaohuang",
			},
			wantErr: false,
		},
		{
			name: "测试登陆不存在的用户",
			args: args{
				username: "UnExist",
				password: "UnExist",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewCheckUserService(context.Background()).CheckUser(&user.CheckUserRequest{Username: tt.args.username, Password: tt.args.password})
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			klog.Info(tt.name + " success")
		})
	}
}

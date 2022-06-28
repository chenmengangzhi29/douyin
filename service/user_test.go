package service

import (
	"testing"
)

//该测试程序不包含像TestMain初始化数据库的程序，所以要与feed_test.go一样运行，因为feed_test.go包含TestMain

//测试用户注册
func TestRegisterUserData(t *testing.T) {
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
			name: "测试注册存在的用户",
			args: args{
				username: "Jerry",
				password: "Jerry123",
			},
			wantErr: true,
		},
		{
			name: "测试注册不存在的用户",
			args: args{
				username: "xiaohuang",
				password: "xiaohuang",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := RegisterUserData(tt.args.username, tt.args.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("RegisterUserData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

//测试用户登陆
func TestLoginUserData(t *testing.T) {
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
				username: "Jerry",
				password: "Jerry123",
			},
			wantErr: false,
		},
		{
			name: "测试登陆不存在的用户",
			args: args{
				username: "unExist",
				password: "unExist",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := LoginUserData(tt.args.username, tt.args.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoginUserData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

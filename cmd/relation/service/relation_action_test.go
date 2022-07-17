package service

import (
	"bytes"
	"context"
	"io"
	"os"
	"testing"

	"github.com/chenmengangzhi29/douyin/dal/db"
	"github.com/chenmengangzhi29/douyin/kitex_gen/relation"
	"github.com/chenmengangzhi29/douyin/pkg/constants"
	"github.com/chenmengangzhi29/douyin/pkg/jwt"
	"github.com/chenmengangzhi29/douyin/pkg/oss"
	"github.com/cloudwego/kitex/pkg/klog"
)

var File []byte
var Token string

func TestMain(m *testing.M) {
	Jwt := jwt.NewJWT([]byte(constants.SecretKey))
	token, err := Jwt.CreateToken(jwt.CustomClaims{
		Id: int64(1),
	})
	if err != nil {
		klog.Errorf("create token fail, %v", err.Error())
		panic(err)
	}
	Token = token

	db.Init()
	oss.Init()

	path := oss.Path + "/public/girl.mp4"
	file, err := os.Open(path)
	if err != nil {
		klog.Errorf("open local file %v fail", path)
		panic(err)
	}
	defer file.Close()
	//处理视频数据
	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, file); err != nil {
		panic(err)
	}
	File = buf.Bytes()

	m.Run()
}

func TestRelationAction(t *testing.T) {
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
				token:      Token,
				toUserId:   2,
				actionType: 1,
			},
			wantErr: false,
		},
		{
			name: "测试取消关注的正常操作",
			args: args{
				token:      Token,
				toUserId:   2,
				actionType: 2,
			},
			wantErr: false,
		},
		{
			name: "测试关注操作的错误用户id",
			args: args{
				token:      Token,
				toUserId:   99999999,
				actionType: 1,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NewRelationActionService(context.Background()).RelationAction(&relation.RelationActionRequest{Token: tt.args.token, ToUserId: tt.args.toUserId, ActionType: int32(tt.args.actionType)})
			if (err != nil) != tt.wantErr {
				t.Errorf("RelationAction() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			klog.Info(tt.name + " success")
		})
	}
}

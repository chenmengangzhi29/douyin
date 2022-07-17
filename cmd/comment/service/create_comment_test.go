package service

import (
	"bytes"
	"context"
	"io"
	"os"
	"testing"

	"github.com/chenmengangzhi29/douyin/dal/db"
	"github.com/chenmengangzhi29/douyin/kitex_gen/comment"
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

func TestCreateComment(t *testing.T) {
	type args struct {
		token       string
		videoId     int64
		commentText string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "测试创建评论的正常操作",
			args: args{
				token:       Token,
				videoId:     1,
				commentText: "hello world",
			},
			wantErr: false,
		},
		{
			name: "测试创建评论的不正确视频id",
			args: args{
				token:       Token,
				videoId:     99999999,
				commentText: "hello world",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			comment, err := NewCreateCommentService(context.Background()).CreateComment(&comment.CreateCommentRequest{Token: tt.args.token, VideoId: tt.args.videoId, CommentText: tt.args.commentText})
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateComment() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			klog.Info(comment)
			klog.Info(tt.name + " success")
		})
	}
}

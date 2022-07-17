package service

import (
	"context"
	"testing"

	"github.com/chenmengangzhi29/douyin/kitex_gen/comment"
	"github.com/cloudwego/kitex/pkg/klog"
)

func TestCommentList(t *testing.T) {
	type args struct {
		token   string
		videoId int64
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "测试评论列表的正常操作",
			args: args{
				token:   Token,
				videoId: 1,
			},
			wantErr: false,
		},
		{
			name: "测试评论列表的不正确视频id",
			args: args{
				token:   Token,
				videoId: 99999999,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			commentList, err := NewCommentListService(context.Background()).CommentList(&comment.CommentListRequest{Token: tt.args.token, VideoId: tt.args.videoId})
			if (err != nil) != tt.wantErr {
				t.Errorf("CommentList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			klog.Info(commentList)
			klog.Info(tt.name + " success")
		})
	}
}

package service

import (
	"context"
	"testing"

	"github.com/chenmengangzhi29/douyin/kitex_gen/comment"
	"github.com/cloudwego/kitex/pkg/klog"
)

func TestDeleteComment(t *testing.T) {
	type args struct {
		token     string
		videoId   int64
		commentId int64
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "测试删除评论的正常操作",
			args: args{
				token:     Token,
				videoId:   1,
				commentId: 5,
			},
			wantErr: false,
		},
		{
			name: "测试删除评论的不正确视频id",
			args: args{
				token:     Token,
				videoId:   99999999,
				commentId: 1,
			},
			wantErr: true,
		},
		{
			name: "测试删除评论的不正确评论id",
			args: args{
				token:     Token,
				videoId:   1,
				commentId: 99999999,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewDeleteCommentService(context.Background()).DeleteComment(&comment.DeleteCommentRequest{Token: tt.args.token, VideoId: tt.args.videoId, CommentId: tt.args.commentId})
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteComment() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			klog.Info(tt.name + " success")
		})
	}
}

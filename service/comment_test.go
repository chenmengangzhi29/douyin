package service

import (
	"douyin/util/logger"
	"testing"
)

//该测试程序不包含像TestMain初始化数据库的程序，所以要与feed_test.go一样运行，因为feed_test.go包含TestMain

func TestCreateCommentData(t *testing.T) {
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
				token:       "JerryJerry123",
				videoId:     1,
				commentText: "hello world",
			},
			wantErr: false,
		},
		{
			name: "测试创建评论的不正确视频id",
			args: args{
				token:       "JerryJerry123",
				videoId:     99999999,
				commentText: "hello world",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := CreateCommentData(tt.args.token, tt.args.videoId, tt.args.commentText)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateCommentData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			logger.Info(tt.name + " success")
		})
	}
}

func TestDeleteCommentData(t *testing.T) {
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
				token:     "JerryJerry123",
				videoId:   1,
				commentId: 1,
			},
			wantErr: false,
		},
		{
			name: "测试删除评论的不正确视频id",
			args: args{
				token:     "JerryJerry123",
				videoId:   99999999,
				commentId: 1,
			},
			wantErr: true,
		},
		{
			name: "测试删除评论的不正确评论id",
			args: args{
				token:     "JerryJerry123",
				videoId:   1,
				commentId: 99999999,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := DeleteCommentData(tt.args.token, tt.args.videoId, tt.args.commentId)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteCommentData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			logger.Info(tt.name + " success")
		})
	}
}

func TestCommentListData(t *testing.T) {
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
				token:   "JerryJerry123",
				videoId: 1,
			},
			wantErr: false,
		},
		{
			name: "测试评论列表的不正确视频id",
			args: args{
				token:   "JerryJerry123",
				videoId: 99999999,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := CommentListData(tt.args.token, tt.args.videoId)
			if (err != nil) != tt.wantErr {
				t.Errorf("CommentListData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			logger.Info(tt.name + " success")
		})
	}
}

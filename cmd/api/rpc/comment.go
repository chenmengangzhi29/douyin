package rpc

import (
	"context"
	"time"

	"github.com/chenmengangzhi29/douyin/kitex_gen/comment"
	"github.com/chenmengangzhi29/douyin/kitex_gen/comment/commentservice"
	"github.com/chenmengangzhi29/douyin/pkg/constants"
	"github.com/chenmengangzhi29/douyin/pkg/errno"
	"github.com/chenmengangzhi29/douyin/pkg/middleware"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/retry"
	etcd "github.com/kitex-contrib/registry-etcd"
	trace "github.com/kitex-contrib/tracer-opentracing"
)

var commentClient commentservice.Client

func initCommentRpc() {
	r, err := etcd.NewEtcdResolver([]string{constants.EtcdAddress})
	if err != nil {
		panic(err)
	}

	c, err := commentservice.NewClient(
		constants.CommentServiceName,
		client.WithMiddleware(middleware.CommonMiddleware),
		client.WithInstanceMW(middleware.ClientMiddleware),
		client.WithMuxConnection(1),
		client.WithRPCTimeout(3*time.Second),
		client.WithConnectTimeout(50*time.Millisecond),
		client.WithFailureRetry(retry.NewFailurePolicy()),
		client.WithSuite(trace.NewDefaultClientSuite()),
		client.WithResolver(r),
	)
	if err != nil {
		panic(err)
	}
	commentClient = c
}

// CreateComment add comment
func CreateComment(ctx context.Context, req *comment.CreateCommentRequest) (*comment.Comment, error) {
	resp, err := commentClient.CreateComment(ctx, req)
	if err != nil {
		return nil, err
	}
	if resp.BaseResp.StatusCode != 0 {
		return nil, errno.NewErrNo(resp.BaseResp.StatusCode, resp.BaseResp.StatusMessage)
	}
	return resp.Comment, nil
}

// DeleteComment delete comment
func DeleteComment(ctx context.Context, req *comment.DeleteCommentRequest) (*comment.Comment, error) {
	resp, err := commentClient.DeleteComment(ctx, req)
	if err != nil {
		return nil, err
	}
	if resp.BaseResp.StatusCode != 0 {
		return nil, errno.NewErrNo(resp.BaseResp.StatusCode, resp.BaseResp.StatusMessage)
	}
	return resp.Comment, nil
}

// CommentList get comment list info
func CommentList(ctx context.Context, req *comment.CommentListRequest) ([]*comment.Comment, error) {
	resp, err := commentClient.CommentList(ctx, req)
	if err != nil {
		return nil, err
	}
	if resp.BaseResp.StatusCode != 0 {
		return nil, errno.NewErrNo(resp.BaseResp.StatusCode, resp.BaseResp.StatusMessage)
	}
	return resp.CommentList, nil
}

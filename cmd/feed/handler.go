package main

import (
	"context"

	"github.com/chenmengangzhi29/douyin/cmd/feed/service"
	"github.com/chenmengangzhi29/douyin/dal/pack"
	"github.com/chenmengangzhi29/douyin/kitex_gen/feed"
	"github.com/chenmengangzhi29/douyin/pkg/errno"
)

// FeedServiceImpl implements the last service interface defined in the IDL.
type FeedServiceImpl struct{}

// Feed implements the FeedServiceImpl interface.
func (s *FeedServiceImpl) Feed(ctx context.Context, req *feed.FeedRequest) (resp *feed.FeedResponse, err error) {
	resp = new(feed.FeedResponse)

	if req.LatestTime <= 0 {
		resp.BaseResp = pack.BuildFeedBaseResp(errno.ParamErr)
		return resp, nil
	}

	videos, nextTime, err := service.NewFeedService(ctx).Feed(req)
	if err != nil {
		resp.BaseResp = pack.BuildFeedBaseResp(err)
		return resp, nil
	}
	resp.BaseResp = pack.BuildFeedBaseResp(errno.Success)
	resp.VideoList = videos
	resp.NextTime = nextTime
	return resp, nil
}

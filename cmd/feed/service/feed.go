package service

import (
	"context"

	"github.com/chenmengangzhi29/douyin/kitex_gen/feed"
	"github.com/chenmengangzhi29/douyin/pkg/jwt"
)

type FeedService struct {
	ctx context.Context
}

// NewFeedService new FeedService
func NewFeedService(ctx context.Context) *FeedService {
	return &FeedService{ctx: ctx}
}

// Feed multiple get list of video info
func (s *FeedService) Feed(req *feed.FeedRequest) ([]*feed.Video, int64, error) {
	currentId, err := s.checkToken(req.Token)
	if err != nil {
		return nil, 0, err
	}

	//get video info
	videoData, err := db.QueryVideoByLatestTime(s.ctx, req.LatestTime)
	if err != nil {
		return  nil, 0, err
	}

	//get video ids and user ids
	videoIds, userIds := pack.Ids(videoData)

	//get user info
	users, err := 




}

//checkToken get userId by token
func (s *FeedService) checkToken(token string) (int64, error) {
	if token == "" {
		return -1, nil
	}
	var Jwt *jwt.JWT
	claim, err := Jwt.ParseToken(token)
	if err != nil {
		return 0, jwt.ErrTokenInvalid
	}
	return claim.Id, nil
}

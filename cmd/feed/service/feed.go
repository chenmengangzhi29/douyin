package service

import (
	"context"
	"sync"

	"github.com/chenmengangzhi29/douyin/dal/db"
	"github.com/chenmengangzhi29/douyin/dal/pack"
	"github.com/chenmengangzhi29/douyin/kitex_gen/feed"
	"github.com/chenmengangzhi29/douyin/pkg/constants"
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
	Jwt := jwt.NewJWT([]byte(constants.SecretKey))
	currentId, _ := Jwt.CheckToken(req.Token)

	//get video info
	videoData, err := db.QueryVideoByLatestTime(s.ctx, req.LatestTime)
	if err != nil {
		return nil, 0, err
	}

	//get video ids and user ids
	videoIds := make([]int64, 0)
	userIds := make([]int64, 0)
	for _, video := range videoData {
		videoIds = append(videoIds, int64(video.ID))
		userIds = append(userIds, video.UserId)
	}

	//get user info
	users, err := db.QueryUserByIds(s.ctx, userIds)
	if err != nil {
		return nil, 0, err
	}
	userMap := make(map[int64]*db.UserRaw)
	for _, user := range users {
		userMap[int64(user.ID)] = user
	}

	var favoriteMap map[int64]*db.FavoriteRaw
	var relationMap map[int64]*db.RelationRaw
	//if user not logged in
	if currentId == -1 {
		favoriteMap = nil
		relationMap = nil
	} else {
		var wg sync.WaitGroup
		wg.Add(2)
		var favoriteErr, relationErr error
		//获取点赞信息
		go func() {
			defer wg.Done()
			favoriteMap, err = db.QueryFavoriteByIds(s.ctx, currentId, videoIds)
			if err != nil {
				favoriteErr = err
				return
			}
		}()
		//获取关注信息
		go func() {
			defer wg.Done()
			relationMap, err = db.QueryRelationByIds(s.ctx, currentId, userIds)
			if err != nil {
				relationErr = err
				return
			}
		}()
		wg.Wait()
		if favoriteErr != nil {
			return nil, 0, favoriteErr
		}
		if relationErr != nil {
			return nil, 0, relationErr
		}

	}

	videos, nextTime := pack.VideoInfo(currentId, videoData, userMap, favoriteMap, relationMap)
	return videos, nextTime, nil
}

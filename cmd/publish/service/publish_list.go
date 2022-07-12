package service

import (
	"context"
	"errors"
	"sync"

	"github.com/chenmengangzhi29/douyin/dal/db"
	"github.com/chenmengangzhi29/douyin/dal/pack"
	"github.com/chenmengangzhi29/douyin/kitex_gen/publish"
	"github.com/chenmengangzhi29/douyin/pkg/jwt"
)

type PublishListService struct {
	ctx context.Context
}

// NewPublishService new PublishService
func NewPublishListService(ctx context.Context) *PublishListService {
	return &PublishListService{ctx: ctx}
}

// PublishList get publish list by userid
func (s *PublishListService) PublishList(req *publish.PublishListRequest) ([]*publish.Video, error) {
	currentId, err := s.checkToken(req.Token)
	if err != nil {
		return nil, err
	}

	videoData, err := db.QueryVideoByUserId(s.ctx, req.UserId)
	if err != nil {
		return nil, err
	}

	videoIds := make([]int64, 0)
	userIds := []int64{req.UserId}
	for _, video := range videoData {
		videoIds = append(videoIds, int64(video.ID))
	}

	users, err := db.QueryUserByIds(s.ctx, userIds)
	if err != nil {
		return nil, err
	}
	if len(users) == 0 {
		return nil, errors.New("user not exist")
	}
	userMap := make(map[int64]*db.UserRaw)
	for _, user := range users {
		userMap[int64(user.ID)] = user
	}

	var favoriteMap map[int64]*db.FavoriteRaw
	var relationMap map[int64]*db.RelationRaw
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
			return nil, favoriteErr
		}
		if relationErr != nil {
			return nil, relationErr
		}
	}

	videoList := pack.PublishInfo(currentId, videoData, userMap, favoriteMap, relationMap)
	return videoList, nil
}

//checkToken get userId by token
func (s *PublishListService) checkToken(token string) (int64, error) {
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

package service

import (
	"context"
	"errors"
	"sync"

	"github.com/chenmengangzhi29/douyin/dal/db"
	"github.com/chenmengangzhi29/douyin/dal/pack"
	"github.com/chenmengangzhi29/douyin/kitex_gen/favorite"
	"github.com/chenmengangzhi29/douyin/pkg/constants"
	"github.com/chenmengangzhi29/douyin/pkg/jwt"
)

type FavoriteListService struct {
	ctx context.Context
}

// NewFavoriteListService new FavoriteListService
func NewFavoriteListService(ctx context.Context) *FavoriteListService {
	return &FavoriteListService{ctx: ctx}
}

// FavoriteList get video information that users like
func (s *FavoriteListService) FavoriteList(req *favorite.FavoriteListRequest) ([]*favorite.Video, error) {
	//获取用户id
	Jwt := jwt.NewJWT([]byte(constants.SecretKey))
	currentId, _ := Jwt.CheckToken(req.Token)

	//检查用户是否存在
	user, err := db.QueryUserByIds(s.ctx, []int64{req.UserId})
	if err != nil {
		return nil, err
	}
	if len(user) == 0 {
		return nil, errors.New("user not exist")
	}

	//获取目标用户的点赞视频id号
	videoIds, err := db.QueryFavoriteById(s.ctx, req.UserId)
	if err != nil {
		return nil, err
	}

	//获取点赞视频的信息
	videoData, err := db.QueryVideoByVideoIds(s.ctx, videoIds)
	if err != nil {
		return nil, err
	}

	//获取点赞视频的用户id号
	userIds := make([]int64, 0)
	for _, video := range videoData {
		userIds = append(userIds, video.UserId)
	}

	//获取点赞视频的用户信息
	users, err := db.QueryUserByIds(s.ctx, userIds)
	if err != nil {
		return nil, err
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
			return nil, favoriteErr
		}
		if relationErr != nil {
			return nil, relationErr
		}

	}

	videoList := pack.VideoList(currentId, videoData, userMap, favoriteMap, relationMap)
	return videoList, nil

}

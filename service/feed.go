package service

import (
	"douyin/model"
	"douyin/repository"
	"errors"
	"fmt"
	"sync"
)

//--------------------------------service---------------------------------------
//该层功能包括鉴权，向repository层获取视频数据，封装视频数据

func QueryVideoData(latestTime int64, token string) ([]*model.Video, int64, error) {
	return NewQueryVideoDataFlow(latestTime, token).Do()
}

func NewQueryVideoDataFlow(latestTime int64, token string) *QueryVideoDataFlow {
	return &QueryVideoDataFlow{
		LatestTime: latestTime,
		Token:      token,
	}
}

type QueryVideoDataFlow struct {
	LatestTime int64
	Token      string
	VideoList  []*model.Video
	NextTime   int64

	CurrentId   int64
	VideoData   []*model.VideoRaw
	UserMap     map[int64]*model.UserRaw
	FavoriteMap map[int64]*model.FavoriteRaw
	RelationMap map[int64]*model.RelationRaw
}

func (f *QueryVideoDataFlow) Do() ([]*model.Video, int64, error) {
	if err := f.checkToken(); err != nil {
		return nil, 0, err
	}
	if err := f.prepareVideoInfo(); err != nil {
		return nil, 0, err
	}
	if err := f.packVideoInfo(); err != nil {
		return nil, 0, err
	}
	return f.VideoList, f.NextTime, nil
}

//鉴权
func (f *QueryVideoDataFlow) checkToken() error {
	if f.Token == "defaultToken" {
		f.CurrentId = -1
		return nil
	}
	currentId, err := repository.NewTokenDaoInstance().QueryUserIdByToken(f.Token)
	if err != nil {
		return err
	}
	f.CurrentId = currentId
	return nil
}

func (f *QueryVideoDataFlow) prepareVideoInfo() error {
	//获取视频信息
	videoData, err := repository.NewVideoDaoInstance().QueryVideoByLatestTime(f.LatestTime)
	if err != nil {
		return err
	}
	f.VideoData = videoData

	//获取视频Id和用户Id
	videoIds := make([]int64, 0)
	userIds := make([]int64, 0)
	for _, video := range f.VideoData {
		videoIds = append(videoIds, video.Id)
		userIds = append(userIds, video.UserId)
	}

	//获取用户信息
	users, err := repository.NewUserDaoInstance().QueryUserByIds(userIds)
	if err != nil {
		return err
	}
	userMap := make(map[int64]*model.UserRaw)
	for _, user := range users {
		userMap[user.Id] = user
	}
	f.UserMap = userMap

	//如果用户未登陆，直接返回
	if f.CurrentId == -1 {
		return nil
	}

	var wg sync.WaitGroup
	wg.Add(2)
	var favoriteErr, relationErr error
	//获取点赞信息
	go func() {
		defer wg.Done()
		favoriteMap, err := repository.NewFavoriteDaoInstance().QueryFavoriteByIds(f.CurrentId, videoIds)
		if err != nil {
			favoriteErr = err
			return
		}
		f.FavoriteMap = favoriteMap
	}()
	//获取关注信息
	go func() {
		defer wg.Done()
		relationMap, err := repository.NewRelationDaoInstance().QueryRelationByIds(f.CurrentId, userIds)
		if err != nil {
			relationErr = err
			return
		}
		f.RelationMap = relationMap
	}()
	wg.Wait()
	if favoriteErr != nil {
		return favoriteErr
	}
	if relationErr != nil {
		return relationErr
	}

	return nil
}

func (f *QueryVideoDataFlow) packVideoInfo() error {
	videoList := make([]*model.Video, 0)
	for _, video := range f.VideoData {
		videoUser, ok := f.UserMap[video.UserId]
		if !ok {
			return errors.New("has no video user info for " + fmt.Sprint(video.UserId))
		}

		var isFavorite bool = false
		var isFollow bool = false

		if f.CurrentId != -1 {
			_, ok := f.FavoriteMap[video.Id]
			if ok {
				isFavorite = true
			}
			_, ok = f.RelationMap[video.UserId]
			if ok {
				isFollow = true
			}
		}
		videoList = append(videoList, &model.Video{
			Id: video.Id,
			Author: &model.User{
				Id:            videoUser.Id,
				Name:          videoUser.Name,
				FollowCount:   videoUser.FollowCount,
				FollowerCount: videoUser.FollowerCount,
				IsFollow:      isFollow,
			},
			PlayUrl:       video.PlayUrl,
			CoverUrl:      video.CoverUrl,
			FavoriteCount: video.FavoriteCount,
			CommentCount:  video.CommentCount,
			IsFavorite:    isFavorite,
			Title:         video.Title,
		})
	}

	f.VideoList = videoList
	if len(f.VideoData) == 0 {
		f.NextTime = 0
	} else {
		f.NextTime = f.VideoData[len(f.VideoData)-1].CreateTime
	}

	return nil
}

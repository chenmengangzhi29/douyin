package service

import (
	"douyin/model"
	"douyin/repository"
	"errors"
	"fmt"
	"sync"
	"time"
)

//点赞操作流，包括鉴权，调用repository层执行点赞操作和取消点赞操作
func FavoriteActionData(token string, videoId int64, actionType int64) error {
	return NewFavoriteActionDataFlow(token, videoId, actionType).Do()
}

func NewFavoriteActionDataFlow(token string, videoId int64, actionType int64) *FavoriteActionDataFlow {
	return &FavoriteActionDataFlow{
		Token:      token,
		VideoId:    videoId,
		ActionType: actionType,
	}
}

type FavoriteActionDataFlow struct {
	Token      string
	VideoId    int64
	ActionType int64

	CurrentId int64
	Favorite  *model.FavoriteRaw
}

func (f *FavoriteActionDataFlow) Do() error {
	if err := f.checkToken(); err != nil {
		return err
	}
	if err := f.checkVideoId(); err != nil {
		return err
	}
	if err := f.prepareFavoriteInfo(); err != nil {
		return err
	}
	return nil
}

//鉴权，登陆的用户才能执行点赞操作
func (f *FavoriteActionDataFlow) checkToken() error {
	currentId, err := repository.NewTokenDaoInstance().QueryUserIdByToken(f.Token)
	if err != nil {
		return err
	}
	f.CurrentId = currentId
	return nil
}

func (f *FavoriteActionDataFlow) checkVideoId() error {
	videos, err := repository.NewVideoDaoInstance().QueryVideoByVideoIds([]int64{f.VideoId})
	if err != nil {
		return err
	}
	if len(videos) == 0 {
		return errors.New("video not exist")
	}
	return nil
}

//若ActionType（操作类型）等于1，则向favorite表创建一条记录，同时向video表的目标video增加点赞数
//若ActionType等于2，则向favorite表删除一条记录，同时向video表的目标video减少点赞数
//若ActionType不等于1和2，则返回错误
func (f *FavoriteActionDataFlow) prepareFavoriteInfo() error {
	if f.ActionType == 1 {
		favorite := &model.FavoriteRaw{
			Id:      time.Now().Unix(),
			UserId:  f.CurrentId,
			VideoId: f.VideoId,
		}
		f.Favorite = favorite

		err := repository.NewFavoriteDaoInstance().CreateFavorite(f.Favorite, f.VideoId)
		if err != nil {
			return err
		}
	}
	if f.ActionType == 2 {
		err := repository.NewFavoriteDaoInstance().DeleteFavorite(f.CurrentId, f.VideoId)
		if err != nil {
			return err
		}

	}
	if f.ActionType != 1 && f.ActionType != 2 {
		return errors.New("action type no equal 1 and 2")
	}
	return nil
}

//点赞列表流，包括鉴权，
//通过repository层准备需要的数据，包括视频数据、用户数据、点赞数据和关注数据
//打包所有数据到videoList
func FavoriteListData(userId int64, token string) ([]model.Video, error) {
	return NewFavoriteListDataFlow(userId, token).Do()
}

func NewFavoriteListDataFlow(userId int64, token string) *FavoriteListDataFlow {
	return &FavoriteListDataFlow{
		UserId: userId,
		Token:  token,
	}
}

type FavoriteListDataFlow struct {
	UserId    int64
	Token     string
	VideoList []model.Video

	CurrentId   int64
	VideoData   []*model.VideoRaw
	UserMap     map[int64]*model.UserRaw
	FavoriteMap map[int64]*model.FavoriteRaw
	RelationMap map[int64]model.RelationRaw
}

func (f *FavoriteListDataFlow) Do() ([]model.Video, error) {
	if err := f.checkToken(); err != nil {
		return nil, err
	}
	if err := f.checkUserId(); err != nil {
		return nil, err
	}
	if err := f.prepareVideoInfo(); err != nil {
		return nil, err
	}
	if err := f.packVideoInfo(); err != nil {
		return nil, err
	}
	return f.VideoList, nil
}

//鉴权
func (f *FavoriteListDataFlow) checkToken() error {
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

func (f *FavoriteListDataFlow) checkUserId() error {
	user, err := repository.NewUserDaoInstance().QueryUserByIds([]int64{f.UserId})
	if err != nil {
		return err
	}
	if len(user) == 0 {
		return errors.New("user not exist")
	}
	return nil
}

//通过repository层准备需要的数据，包括视频数据、用户数据、点赞数据和关注数据
func (f *FavoriteListDataFlow) prepareVideoInfo() error {
	//获取目标用户的点赞视频id号
	videoIds, err := repository.NewFavoriteDaoInstance().QueryFavoriteById(f.UserId)
	if err != nil {
		return err
	}

	//获取点赞视频的信息
	videoData, err := repository.NewVideoDaoInstance().QueryVideoByVideoIds(videoIds)
	if err != nil {
		return err
	}
	f.VideoData = videoData

	//获取点赞视频的用户id号
	userIds := make([]int64, 0)
	for _, video := range f.VideoData {
		userIds = append(userIds, video.UserId)
	}

	//获取点赞视频的用户信息
	users, err := repository.NewUserDaoInstance().QueryUserByIds(userIds)
	if err != nil {
		return err
	}
	userMap := make(map[int64]*model.UserRaw)
	for _, user := range users {
		userMap[user.Id] = &user
	}
	f.UserMap = userMap

	//如果用户未登陆，直接返回
	if f.CurrentId == -1 {
		return nil
	}

	var wg sync.WaitGroup
	wg.Add(2)
	var favoriteErr, relationErr error
	//获取当前登陆用户对目标用户点赞视频的点赞信息
	go func() {
		defer wg.Done()
		favoriteMap, err := repository.NewFavoriteDaoInstance().QueryFavoriteByIds(f.CurrentId, videoIds)
		if err != nil {
			favoriteErr = err
			return
		}
		f.FavoriteMap = favoriteMap
	}()
	//获取当前登陆用户对点赞视频的用户的关注信息
	go func() {
		defer wg.Done()
		relatinMap, err := repository.NewRelationDaoInstance().QueryRelationByIds(f.CurrentId, userIds)
		if err != nil {
			relationErr = err
			return
		}
		f.RelationMap = relatinMap
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

//打包所有数据到videoList
func (f *FavoriteListDataFlow) packVideoInfo() error {
	videoList := make([]model.Video, 0)
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

		videoList = append(videoList, model.Video{
			Id: video.Id,
			Author: model.User{
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
	return nil
}

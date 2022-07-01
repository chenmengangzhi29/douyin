package service

import (
	"douyin/model"
	"douyin/repository"
	"errors"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

//--------------------service-----------------------------------
//上传视频数据流，包括获取当前用户id，处理视频名，确定视频保存的文件夹路径，处理Oss上的object路径，处理好要上传到mysql的视频结构
func PublishUserVideoData(token string, data *multipart.FileHeader, title string, c *gin.Context) error {
	return NewPublishUserVideoDataFlow(token, data, title, c).Do()
}

func NewPublishUserVideoDataFlow(token string, data *multipart.FileHeader, title string, c *gin.Context) *PublishUserVideoDataFlow {
	return &PublishUserVideoDataFlow{
		Token: token,
		Video: data,
		Title: title,
		Gin:   c,
	}
}

type PublishUserVideoDataFlow struct {
	Token string
	Video *multipart.FileHeader
	Title string
	Gin   *gin.Context

	CurrentId int64
	VideoData model.VideoRaw
}

func (f *PublishUserVideoDataFlow) Do() error {
	if err := f.checkToken(); err != nil {
		return err
	}
	if err := f.publishVideo(); err != nil {
		return err
	}
	return nil
}

func (f *PublishUserVideoDataFlow) checkToken() error {
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

func (f *PublishUserVideoDataFlow) publishVideo() error {
	//处理视频名
	filename := filepath.Base(f.Video.Filename)
	finalName := fmt.Sprintf("%d_%s", f.CurrentId, filename)

	//将视频保存到本地文件夹
	saveFile := filepath.Join("./public/", finalName)
	err := repository.NewVideoDaoInstance().PublishVideoToPublic(f.Video, saveFile, f.Gin)
	if err != nil {
		return err
	}

	//将本地视频上传到oss同时将视频信息上传到mysql
	object := "video/" + finalName

	video := model.VideoRaw{
		Id:         time.Now().Unix(),
		UserId:     f.CurrentId,
		Title:      f.Title,
		PlayUrl:    "https://dousheng1.oss-cn-shenzhen.aliyuncs.com/" + object,
		CreateTime: time.Now().Unix(),
	}
	f.VideoData = video

	var wg sync.WaitGroup
	wg.Add(2)
	var OssErr, MysqlErr error

	go func() {
		defer wg.Done()
		err = repository.NewVideoDaoInstance().PublishVideoToOss(object, saveFile)
		if err != nil {
			OssErr = err
			return
		}
	}()
	go func() {
		defer wg.Done()
		err = repository.NewVideoDaoInstance().PublishVideoData(f.VideoData)
		if err != nil {
			MysqlErr = err
			return
		}
	}()

	wg.Wait()
	if OssErr != nil {
		return OssErr
	}
	if MysqlErr != nil {
		return MysqlErr
	}

	return nil
}

//查询用户视频列表流，包括鉴权操作，通过repository层获取需要的各种信息（包括视频，用户，点赞，关注信息），将所有信息转换成视频列表。
func QueryUserVideoList(token string, userId int64) ([]model.Video, error) {
	return NewQueryUserVideoListFlow(token, userId).Do()
}

func NewQueryUserVideoListFlow(token string, userId int64) *QueryUserVideoListFlow {
	return &QueryUserVideoListFlow{
		Token:  token,
		UserId: userId,
	}
}

type QueryUserVideoListFlow struct {
	Token     string
	UserId    int64
	VideoList []model.Video

	CurrentId   int64
	VideoData   []*model.VideoRaw
	UserMap     map[int64]*model.UserRaw
	FavoriteMap map[int64]*model.FavoriteRaw
	RelationMap map[int64]model.RelationRaw
}

func (f *QueryUserVideoListFlow) Do() ([]model.Video, error) {
	if err := f.checkToken(); err != nil {
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

func (f *QueryUserVideoListFlow) checkToken() error {
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

func (f *QueryUserVideoListFlow) prepareVideoInfo() error {
	videoData, err := repository.NewVideoDaoInstance().QueryVideoByUserId(f.UserId)
	if err != nil {
		return err
	}
	f.VideoData = videoData

	videoIds := make([]int64, 0)
	userIds := []int64{f.UserId}
	for _, video := range f.VideoData {
		videoIds = append(videoIds, video.Id)
	}

	users, err := repository.NewUserDaoInstance().QueryUserByIds(userIds)
	if err != nil {
		return err
	}
	if len(users) == 0 {
		return errors.New("user not exist")
	}
	userMap := make(map[int64]*model.UserRaw)
	for _, user := range users {
		userMap[user.Id] = user
	}
	f.UserMap = userMap

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

func (f *QueryUserVideoListFlow) packVideoInfo() error {
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

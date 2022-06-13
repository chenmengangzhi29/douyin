package controller

import (
	"douyin/models"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

//--------------------------------repository------------------------------------
//该层功能是直接向数据库获取信息

//------------视频--------------
type VideoRaw struct {
	Id            int64  `gorm:"column:id"`
	UserId        int64  `gorm:"column:user_id"`
	Title         string `gorm:"column:title"`
	PlayUrl       string `gorm:"column:play_url"`
	CoverUrl      string `gorm:"column:cover_url"`
	FavoriteCount int64  `gorm:"column:favorite_count"`
	CommentCount  int64  `gorm:"column:comment_count"`
	CreateTime    int64  `gorm:"column:create_time"`
}

func (VideoRaw) TableName() string {
	return "video"
}

type VideoDao struct {
}

var videoDao *VideoDao
var videoOnce sync.Once

func NewVideoDaoInstance() *VideoDao {
	videoOnce.Do(
		func() {
			videoDao = &VideoDao{}
		})
	return videoDao
}

//根据最新时间戳获取视频信息
func (*VideoDao) QueryVideoByLatestTime(latestTime int64) ([]*VideoRaw, error) {
	var videos []*VideoRaw
	err := models.DB.Limit(30).Order("create_time desc").Where("create_time < ?", latestTime).Find(&videos).Error
	if err == gorm.ErrRecordNotFound {
		return nil, errors.New("do not find video")
	}
	if err != nil {
		return nil, errors.New("find video error")
	}
	return videos, nil
}

//-------------用户--------------
type UserRaw struct {
	Id            int64  `gorm:"column:id"`
	Name          string `gorm:"column:name"`
	Password      string `gorm:"column:password"`
	FollowCount   int64  `gorm:"column:follow_count"`
	FollowerCount int64  `gorm:"column:follower_count"`
	Token         string `gorm:"column:token"`
}

func (UserRaw) TableName() string {
	return "user"
}

type TokenDao struct {
}
type UserDao struct {
}

var tokenDao *TokenDao
var tokenOnce sync.Once

var userDao *UserDao
var userOnce sync.Once

func NewTokenDaoInstance() *TokenDao {
	tokenOnce.Do(
		func() {
			tokenDao = &TokenDao{}
		})
	return tokenDao
}

func NewUserDaoInstance() *UserDao {
	userOnce.Do(
		func() {
			userDao = &UserDao{}
		})
	return userDao
}

//根据token获取用户id
func (*TokenDao) QueryUserIdByToken(token string) (int64, error) {
	var user UserRaw
	err := models.DB.Where("token = ?", token).Find(&user).Error
	if err != nil {
		return -1, errors.New("check token fail")
	}
	return user.Id, nil
}

//根据用户id获取用户信息
func (*UserDao) QueryUserByIds(userIds []int64) (map[int64]*UserRaw, error) {
	var users []*UserRaw
	err := models.DB.Where("id in (?)", userIds).Find(&users).Error
	if err != nil {
		return nil, errors.New("query user fail")
	}
	userMap := make(map[int64]*UserRaw)
	for _, user := range users {
		userMap[user.Id] = user
	}
	return userMap, nil
}

//-------------点赞--------------
type FavoriteRaw struct {
	Id      int64 `gorm:"column:id"`
	UserId  int64 `gorm:"column:user_id"`
	VideoId int64 `gorm:"column:video_id"`
}

func (FavoriteRaw) TableName() string {
	return "favorite"
}

type FavoriteDao struct {
}

var favoriteDao *FavoriteDao
var favoriteOnce sync.Once

func NewFavoriteDaoInstance() *FavoriteDao {
	favoriteOnce.Do(
		func() {
			favoriteDao = &FavoriteDao{}
		})
	return favoriteDao
}

//根据当前用户id和视频id获取点赞信息
func (*FavoriteDao) QueryFavoriteByIds(currentId int64, videoIds []int64) (map[int64]*FavoriteRaw, error) {
	var favorites []*FavoriteRaw
	err := models.DB.Where("user_id = ? AND video_id IN ?", currentId, videoIds).Find(&favorites).Error
	if err == gorm.ErrRecordNotFound {
		return nil, errors.New("favorite record not found")
	}
	if err != nil {
		return nil, errors.New("query favorite record fail")
	}
	favoriteMap := make(map[int64]*FavoriteRaw)
	for _, favorite := range favorites {
		favoriteMap[favorite.VideoId] = favorite
	}
	return favoriteMap, nil
}

//-------------关注--------------
type RelationRaw struct {
	Id       int64 `gorm:"column:id"`
	UserId   int64 `gorm:"column:user_id"`
	ToUserId int64 `gorm:"column:to_user_id"`
	Status   int64 `gorm:"column:status"`
}

func (RelationRaw) TableName() string {
	return "relation"
}

type RelationDao struct {
}

var relationDao *RelationDao
var relationOnce sync.Once

func NewRelationDaoInstance() *RelationDao {
	relationOnce.Do(
		func() {
			relationDao = &RelationDao{}
		})
	return relationDao
}

//根据当前用户id和视频拥有者id获取关注信息
func (*RelationDao) QueryRelationByIds(currentId int64, userIds []int64) (map[int64]*RelationRaw, error) {
	var relations []*RelationRaw
	err := models.DB.Where("user_id = ? AND to_user_id IN ? AND status IN ?", currentId, userIds, []int64{0, -1}).Or("user_id IN ? AND to_user_id = ? AND status = ?", userIds, currentId, 1).Find(&relations).Error
	if err == gorm.ErrRecordNotFound {
		return nil, errors.New("relation record not found")
	}
	if err != nil {
		return nil, errors.New("query relation record fail")
	}
	relationMap := make(map[int64]*RelationRaw)
	for _, relation := range relations {
		if relation.Status == 1 {
			relationMap[relation.UserId] = relation
		} else {
			relationMap[relation.ToUserId] = relation
		}
	}
	return relationMap, nil
}

//--------------------------------service---------------------------------------
//该层功能包括鉴权，向repository层获取视频数据，封装视频数据

func QueryVideoData(latestTime int64, token string) ([]Video, int64, error) {
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
	VideoList  []Video
	NextTime   int64

	CurrentId   int64
	VideoData   []*VideoRaw
	UserMap     map[int64]*UserRaw
	FavoriteMap map[int64]*FavoriteRaw
	RelationMap map[int64]*RelationRaw
}

func (f *QueryVideoDataFlow) Do() ([]Video, int64, error) {
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
	currentId, err := NewTokenDaoInstance().QueryUserIdByToken(f.Token)
	if err != nil {
		return err
	}
	f.CurrentId = currentId
	return nil
}

func (f *QueryVideoDataFlow) prepareVideoInfo() error {
	//获取视频信息
	videoData, err := NewVideoDaoInstance().QueryVideoByLatestTime(f.LatestTime)
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
	userMap, err := NewUserDaoInstance().QueryUserByIds(userIds)
	if err != nil {
		return err
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
		favoriteMap, err := NewFavoriteDaoInstance().QueryFavoriteByIds(f.CurrentId, videoIds)
		if err != nil {
			favoriteErr = err
			return
		}
		f.FavoriteMap = favoriteMap
	}()
	//获取关注信息
	go func() {
		defer wg.Done()
		relationMap, err := NewRelationDaoInstance().QueryRelationByIds(f.CurrentId, userIds)
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
	videoList := make([]Video, 0)
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
		videoList = append(videoList, Video{
			Id: video.Id,
			Author: models.User{
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

//----------------------------handler-------------------------------------
//该层功能包括处理传入参数，向service层获取视频信息，封装成响应信息
type FeedResponse struct {
	Response
	// VideoData VideoData `json:"video_data,omitempty"`
	VideoList []Video `json:"video_list,omitempty"`
	NextTime  int64   `json:"next_time,omitempty"`
}

func QueryVideoFeed(token string, latestTimeStr string) FeedResponse {
	//处理传入参数
	latestTime, err := strconv.ParseInt(latestTimeStr, 10, 64)
	if err != nil {
		return FeedResponse{
			Response: Response{StatusCode: -1, StatusMsg: err.Error()},
		}
	}

	//获取视频
	videoList, nextTime, err := QueryVideoData(latestTime, token)
	if err != nil {
		return FeedResponse{
			Response: Response{StatusCode: -1, StatusMsg: err.Error()},
		}
	}

	// fmt.Println(videoData)

	return FeedResponse{
		Response:  Response{StatusCode: 0, StatusMsg: "success"},
		VideoList: videoList,
		NextTime:  nextTime,
	}
}

//------------------controller--------------------------------------------
// Feed same demo video list for every request
// 该层功能包括获取传入参数，向handler获取视频信息，返回响应信息
func Feed(c *gin.Context) {
	//获取传入参数
	token := c.DefaultQuery("token", "defaultToken")
	defaultTime := time.Now().Unix()
	defaultTimeStr := strconv.Itoa(int(defaultTime))
	latestTimeStr := c.DefaultQuery("latest_time", defaultTimeStr)

	//获取视频
	feedResponse := QueryVideoFeed(token, latestTimeStr)

	c.JSON(http.StatusOK, feedResponse)

}

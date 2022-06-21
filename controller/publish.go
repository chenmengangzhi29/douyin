package controller

import (
	"douyin/models"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func FileServer(c *gin.Context) {
	dir, err := os.Getwd()
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}
	fileName := c.Param("name")
	path := dir + "public/" + fileName
	c.File(path)
}

// Publish check token then save upload file to public directory
func Publish(c *gin.Context) {
	token := c.PostForm("token")

	userId, err := NewTokenDaoInstance().QueryUserIdByToken(token)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}

	data, err := c.FormFile("data")
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}

	filename := filepath.Base(data.Filename)
	finalName := fmt.Sprintf("%d_%s", userId, filename)

	saveFile := filepath.Join("./public/", finalName)
	if err := c.SaveUploadedFile(data, saveFile); err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}

	title := c.PostForm("title")

	// config, err := ini.Load("./models/app.ini")
	// if err != nil {
	// 	c.JSON(http.StatusOK, Response{
	// 		StatusCode: 1,
	// 		StatusMsg:  err.Error(),
	// 	})
	// 	return
	// }

	// ip := config.Section("mysql").Key("ip").String()

	playUrl := "http://localhost:8080/douyin/" + saveFile

	video := VideoRaw{
		Id:         time.Now().Unix(),
		UserId:     userId,
		Title:      title,
		PlayUrl:    playUrl,
		CreateTime: time.Now().Unix(),
	}

	if err := models.DB.Create(&video).Error; err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		StatusCode: 0,
		StatusMsg:  finalName + " uploaded successfully",
	})
}

// PublishList all users have same publish video list
//--------------------repository---------------------------------
//向数据库取视频数据
//在feed.go中已经实现一些repository功能，可以直接复用

//新增通过用户id获取视频数据的功能
func (*VideoDao) QueryVideoByUserId(userId int64) ([]*VideoRaw, error) {
	var videos []*VideoRaw
	err := models.DB.Order("create_time desc").Where("user_id = ?", userId).Find(&videos).Error
	if err == gorm.ErrRecordNotFound {
		return nil, errors.New("do not find video")
	}
	if err != nil {
		return nil, errors.New("find video error")
	}
	return videos, nil
}

//--------------------service-----------------------------------
//鉴权操作，通过repository层获取需要的各种信息（包括视频，用户，点赞，关注信息），将所有信息转换成视频列表。
func QueryUserVideoList(token string, userId int64) ([]Video, error) {
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
	VideoList []Video

	CurrentId   int64
	VideoData   []*VideoRaw
	UserMap     map[int64]*UserRaw
	FavoriteMap map[int64]*FavoriteRaw
	RelationMap map[int64]*RelationRaw
}

func (f *QueryUserVideoListFlow) Do() ([]Video, error) {
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
	currentId, err := NewTokenDaoInstance().QueryUserIdByToken(f.Token)
	if err != nil {
		return err
	}
	f.CurrentId = currentId
	return nil
}

func (f *QueryUserVideoListFlow) prepareVideoInfo() error {
	videoData, err := NewVideoDaoInstance().QueryVideoByUserId(f.UserId)
	if err != nil {
		return err
	}
	f.VideoData = videoData

	videoIds := make([]int64, 0)
	userIds := []int64{f.UserId}
	for _, video := range f.VideoData {
		videoIds = append(videoIds, video.Id)
	}

	userMap, err := NewUserDaoInstance().QueryUserByIds(userIds)
	if err != nil {
		return err
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

func (f *QueryUserVideoListFlow) packVideoInfo() error {
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
	return nil
}

//---------------------handler--------------------------------
//处理传入参数，调用service层函数查询视频列表，封装响应
type VideoListResponse struct {
	Response
	VideoList []Video `json:"video_list,omitempty"`
}

func QueryVideoList(token string, userIdStr string) VideoListResponse {
	userId, err := strconv.ParseInt(userIdStr, 10, 64)
	if err != nil {
		return VideoListResponse{
			Response: Response{StatusCode: -1, StatusMsg: err.Error()},
		}
	}

	videoList, err := QueryUserVideoList(token, userId)
	if err != nil {
		return VideoListResponse{
			Response: Response{StatusCode: -1, StatusMsg: err.Error()},
		}
	}

	return VideoListResponse{
		Response:  Response{StatusCode: 0, StatusMsg: "success"},
		VideoList: videoList,
	}
}

//-------------------controller---------------------------------
//获取传入参数，调用handler层函数查询视频列表
func PublishList(c *gin.Context) {
	token := c.Query("token")
	userIdStr := c.Query("user_id")

	videoListResponse := QueryVideoList(token, userIdStr)

	c.JSON(http.StatusOK, videoListResponse)

}

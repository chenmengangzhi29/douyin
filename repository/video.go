package repository

import (
	"douyin/model"
	"errors"
	"mime/multipart"
	"sync"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

//--------------------------------repository------------------------------------
//该层功能是直接向数据库获取信息

//-------------视频--------------------

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
func (*VideoDao) QueryVideoByLatestTime(latestTime int64) ([]*model.VideoRaw, error) {
	var videos []*model.VideoRaw
	err := model.DB.Table("video").Limit(30).Order("create_time desc").Where("create_time < ?", latestTime).Find(&videos).Error
	if err == gorm.ErrRecordNotFound {
		return nil, errors.New("do not find video")
	}
	if err != nil {
		return nil, errors.New("find video error")
	}
	return videos, nil
}

//新增通过用户id获取视频数据的功能
func (*VideoDao) QueryVideoByUserId(userId int64) ([]*model.VideoRaw, error) {
	var videos []*model.VideoRaw
	err := model.DB.Table("video").Order("create_time desc").Where("user_id = ?", userId).Find(&videos).Error
	if err == gorm.ErrRecordNotFound {
		return nil, errors.New("do not find video")
	}
	if err != nil {
		return nil, errors.New("find video error")
	}
	return videos, nil
}

//将视频保存到本地文件夹中
func (*VideoDao) PublishVideoToPublic(video *multipart.FileHeader, saveFile string, c *gin.Context) error {
	if err := c.SaveUploadedFile(video, saveFile); err != nil {
		return err
	}
	return nil
}

//将本地文件夹中的视频上传到Oss
func (*VideoDao) PublishVideoToOss(object string, saveFile string) error {
	err := model.Bucket.PutObjectFromFile(object, saveFile)
	if err != nil {
		return err
	}
	return nil
}

//向video表添加一条记录
func (*VideoDao) PublishVideoData(videoData model.VideoRaw) error {
	if err := model.DB.Table("video").Create(&videoData).Error; err != nil {
		return err
	}
	return nil
}

//更新视频的点赞数量
func (*VideoDao) UpdateFavoriteCount(videoId int64, actionType int64) error {
	if actionType == 1 {
		err := model.DB.Table("video").Where("id = ?", videoId).Update("favorite_count", gorm.Expr("favorite_count + ?", 1)).Error
		if err != nil {
			return err
		}
	}
	if actionType == 2 {
		err := model.DB.Table("video").Where("id = ?", videoId).Update("favorite_count", gorm.Expr("favorite_count - ?", 1)).Error
		if err != nil {
			return err
		}
	}
	return nil
}

//通过一系列视频id号获取一系列视频信息
func (*VideoDao) QueryVideoByVideoIds(videoIds []int64) ([]*model.VideoRaw, error) {
	var videos []*model.VideoRaw
	err := model.DB.Table("video").Where("id in (?)", videoIds).Find(&videos).Error
	if err != nil {
		return nil, err
	}
	return videos, nil
}

//通过视频id增加视频的评论数
func (*VideoDao) AddCommentCount(videoId int64) error {
	err := model.DB.Table("video").Where("id = ?", videoId).Update("comment_count", gorm.Expr("comment_count + ?", 1)).Error
	if err != nil {
		return err
	}
	return nil
}

func (*VideoDao) SubCommentCount(videoId int64) error {
	err := model.DB.Table("video").Where("id = ?", videoId).Update("comment_count", gorm.Expr("comment_count - ?", 1)).Error
	if err != nil {
		return err
	}
	return nil
}

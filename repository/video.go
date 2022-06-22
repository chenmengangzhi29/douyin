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
	err := model.DB.Limit(30).Order("create_time desc").Where("create_time < ?", latestTime).Find(&videos).Error
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
	err := model.DB.Order("create_time desc").Where("user_id = ?", userId).Find(&videos).Error
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

func (*VideoDao) PublishVideoData(videoData model.VideoRaw) error {
	if err := model.DB.Table("video").Create(&videoData).Error; err != nil {
		return err
	}
	return nil
}

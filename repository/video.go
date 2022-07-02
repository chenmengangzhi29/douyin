package repository

import (
	"bytes"
	"douyin/model"
	"douyin/util/logger"
	"errors"
	"sync"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
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
		logger.Error("QueryVideoByLatestTime do not find video " + err.Error())
		return nil, errors.New("do not find video")
	}
	if err != nil {
		logger.Error("QueryVideoByLatestTime find video error " + err.Error())
		return nil, errors.New("find video error")
	}
	return videos, nil
}

//新增通过用户id获取视频数据的功能
func (*VideoDao) QueryVideoByUserId(userId int64) ([]*model.VideoRaw, error) {
	var videos []*model.VideoRaw
	err := model.DB.Table("video").Order("create_time desc").Where("user_id = ?", userId).Find(&videos).Error
	if err != nil {
		logger.Error("QueryVideoByUserId find video error " + err.Error())
		return nil, errors.New("find video error")
	}
	return videos, nil
}

// //将视频保存到本地文件夹中
// func (*VideoDao) PublishVideoToPublic(video *multipart.FileHeader, saveFile string, c *gin.Context) error {
// 	if err := c.SaveUploadedFile(video, saveFile); err != nil {
// 		logger.Error("PublishVideoToPublic error " + err.Error())
// 		return err
// 	}
// 	return nil
// }

//将视频上传到Oss
func (*VideoDao) PublishVideoToOss(objectKey string, video *bytes.Reader) error {
	err := model.Bucket.PutObject(objectKey, video)
	if err != nil {
		logger.Error("PublishVideoToOss error " + err.Error())
		return err
	}
	return nil
}

//从oss上获取播放地址
func (*VideoDao) QueryOssVideoURL(objectKey string) (string, error) {
	signedURL, err := model.Bucket.SignURL(objectKey, oss.HTTPPut, 60)
	if err != nil {
		logger.Errorf("Query %v URL fail, %v", objectKey, err.Error())
		return "", err
	}
	return signedURL, nil
}

//向video表添加一条记录
func (*VideoDao) PublishVideoData(videoData *model.VideoRaw) error {
	if err := model.DB.Table("video").Create(&videoData).Error; err != nil {
		logger.Error("PublishVideoData error " + err.Error())
		return err
	}
	return nil
}

//通过一系列视频id号获取一系列视频信息
func (*VideoDao) QueryVideoByVideoIds(videoIds []int64) ([]*model.VideoRaw, error) {
	var videos []*model.VideoRaw
	err := model.DB.Table("video").Where("id in (?)", videoIds).Find(&videos).Error
	if err != nil {
		logger.Error("QueryVideoByVideoIds error " + err.Error())
		return nil, err
	}
	return videos, nil
}

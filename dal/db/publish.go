package db

import (
	"context"
	"errors"

	"github.com/chenmengangzhi29/douyin/pkg/logger"
)

//向video表添加一条记录
func PublishVideoData(ctx context.Context, videoData *VideoRaw) error {
	if err := DB.WithContext(ctx).Create(&videoData).Error; err != nil {
		logger.Error("PublishVideoData error " + err.Error())
		return err
	}
	return nil
}

//新增通过用户id获取视频数据的功能
func QueryVideoByUserId(ctx context.Context, userId int64) ([]*VideoRaw, error) {
	var videos []*VideoRaw
	err := DB.Table("video").WithContext(ctx).Order("create_time desc").Where("user_id = ?", userId).Find(&videos).Error
	if err != nil {
		logger.Error("QueryVideoByUserId find video error " + err.Error())
		return nil, errors.New("find video error")
	}
	return videos, nil
}

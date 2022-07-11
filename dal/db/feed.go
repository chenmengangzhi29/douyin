package db

import (
	"context"
	"time"

	"github.com/chenmengangzhi29/douyin/pkg/constants"
	"gorm.io/gorm"
)

// Video Gorm Data Structures
type VideoRaw struct {
	gorm.Model
	UserId        int64     `gorm:"column:user_id;not null;index:idx_userid"`
	Title         string    `gorm:"column:title;type:varchar(128);not null"`
	PlayUrl       string    `gorm:"column:play_url;varchar(128);not null"`
	CoverUrl      string    `gorm:"column:cover_url;varchar(128);not null"`
	FavoriteCount int64     `gorm:"column:favorite_count;default:0"`
	CommentCount  int64     `gorm:"column:comment_count;default:0"`
	UpdatedAt     time.Time `gorm:"column:update_time;not null;index:idx_update"`
}

func (v *VideoRaw) TableName() string {
	return constants.VideoTableName
}

//QueryVideoByLatestTime query video info by latest create Time
func QueryVideoByLatestTime(ctx context.Context, latestTime int64) ([]*VideoRaw, error) {
	var videos []*VideoRaw
	time := time.UnixMilli(latestTime)
	err := DB.WithContext(ctx).Limit(30).Order("update_time desc").Where("update_time < ?", time).Find(&videos).Error
	if err != nil {
		return videos, err
	}
	return videos, nil
}

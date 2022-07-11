package db

import (
	"context"

	"github.com/chenmengangzhi29/douyin/model"
	"gorm.io/gorm"
)

// Favorite Gorm Data Structures
type FavoriteRaw struct {
	gorm.Model
	UserId  int64 `gorm:"column:user_id;not null;index:idx_userid"`
	VideoId int64 `gorm:"column:video_id;not null;index:idx_videoid"`
}

func (FavoriteRaw) TableName() string {
	return "favorite"
}

//根据当前用户id和视频id获取点赞信息
func QueryFavoriteByIds(ctx context.Context, currentId int64, videoIds []int64) (map[int64]*FavoriteRaw, error) {
	var favorites []*FavoriteRaw
	err := model.DB.WithContext(ctx).Where("user_id = ? AND video_id IN ?", currentId, videoIds).Find(&favorites).Error
	if err != nil {
		return nil, err
	}
	favoriteMap := make(map[int64]*FavoriteRaw)
	for _, favorite := range favorites {
		favoriteMap[favorite.VideoId] = favorite
	}
	return favoriteMap, nil
}

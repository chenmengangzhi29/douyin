package repository

import (
	"douyin/model"
	"errors"
	"sync"

	"gorm.io/gorm"
)

//-------------点赞--------------

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
func (*FavoriteDao) QueryFavoriteByIds(currentId int64, videoIds []int64) (map[int64]*model.FavoriteRaw, error) {
	var favorites []*model.FavoriteRaw
	err := model.DB.Where("user_id = ? AND video_id IN ?", currentId, videoIds).Find(&favorites).Error
	if err == gorm.ErrRecordNotFound {
		return nil, errors.New("favorite record not found")
	}
	if err != nil {
		return nil, errors.New("query favorite record fail")
	}
	favoriteMap := make(map[int64]*model.FavoriteRaw)
	for _, favorite := range favorites {
		favoriteMap[favorite.VideoId] = favorite
	}
	return favoriteMap, nil
}

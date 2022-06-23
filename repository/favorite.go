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

//向favorite表添加一条记录
func (*FavoriteDao) CreateFavorite(favorite *model.FavoriteRaw) error {
	err := model.DB.Table("favorite").Create(favorite).Error
	if err != nil {
		return err
	}
	return nil
}

//删除favorite表的一条记录
func (*FavoriteDao) DeleteFavorite(currentId int64, videoId int64) error {
	var favorite model.FavoriteRaw
	err := model.DB.Table("favorite").Where("user_id = ? AND video_id = ?", currentId, videoId).Delete(&favorite).Error
	if err != nil {
		return err
	}
	return nil
}

//通过一个用户id查询出该用户点赞的所有视频id号
func (*FavoriteDao) QueryFavoriteById(userId int64) ([]int64, error) {
	var favorites []*model.FavoriteRaw
	err := model.DB.Table("favorite").Where("user_id = ?", userId).Find(&favorites).Error
	if err == gorm.ErrRecordNotFound {
		return nil, errors.New("favorite record not found")
	}
	if err != nil {
		return nil, errors.New("query favorite record fail")
	}
	videoIds := make([]int64, 0)
	for _, favorite := range favorites {
		videoIds = append(videoIds, favorite.VideoId)
	}
	return videoIds, nil
}

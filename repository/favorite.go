package repository

import (
	"douyin/model"
	"douyin/util/logger"
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
	err := model.DB.Table("favorite").Where("user_id = ? AND video_id IN ?", currentId, videoIds).Find(&favorites).Error
	if err == gorm.ErrRecordNotFound {
		logger.Error("found favorite record fail " + err.Error())
		return nil, err
	}
	if err != nil {
		logger.Error("quert favorite record fail " + err.Error())
		return nil, err
	}
	favoriteMap := make(map[int64]*model.FavoriteRaw)
	for _, favorite := range favorites {
		favoriteMap[favorite.VideoId] = favorite
	}
	return favoriteMap, nil
}

//通过事务向favorite表添加一条记录,同时添加视频点赞数
func (*FavoriteDao) CreateFavorite(favorite *model.FavoriteRaw, videoId int64) error {
	model.DB.Transaction(func(tx *gorm.DB) error {
		err := tx.Table("video").Where("id = ?", videoId).Update("favorite_count", gorm.Expr("favorite_count + ?", 1)).Error
		if err != nil {
			logger.Error("AddFavoriteCount error " + err.Error())
			return err
		}

		err = tx.Table("favorite").Create(favorite).Error
		if err != nil {
			logger.Error("create favorite record fail " + err.Error())
			return err
		}

		return nil
	})
	return nil
}

//删除favorite表的一条记录,同时减少视频点赞数
func (*FavoriteDao) DeleteFavorite(currentId int64, videoId int64) error {
	model.DB.Transaction(func(tx *gorm.DB) error {
		var favorite *model.FavoriteRaw
		err := tx.Table("favorite").Where("user_id = ? AND video_id = ?", currentId, videoId).Delete(&favorite).Error
		if err != nil {
			logger.Error("delete favorite record fail " + err.Error())
			return err
		}

		err = tx.Table("video").Where("id = ?", videoId).Update("favorite_count", gorm.Expr("favorite_count - ?", 1)).Error
		if err != nil {
			logger.Error("SubFavoriteCount error " + err.Error())
			return err
		}
		return nil
	})
	return nil
}

//通过一个用户id查询出该用户点赞的所有视频id号
func (*FavoriteDao) QueryFavoriteById(userId int64) ([]int64, error) {
	var favorites []*model.FavoriteRaw
	err := model.DB.Table("favorite").Where("user_id = ?", userId).Find(&favorites).Error
	if err != nil {
		logger.Error("query favorite record fail " + err.Error())
		return nil, err
	}
	videoIds := make([]int64, 0)
	for _, favorite := range favorites {
		videoIds = append(videoIds, favorite.VideoId)
	}
	return videoIds, nil
}

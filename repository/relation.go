package repository

import (
	"douyin/model"
	"douyin/util/logger"
	"sync"
	"time"

	"gorm.io/gorm"
)

//-------------关注--------------

type RelationDao struct {
}

var relationDao *RelationDao
var relationOnce sync.Once

func NewRelationDaoInstance() *RelationDao {
	relationOnce.Do(
		func() {
			relationDao = &RelationDao{}
		})
	return relationDao
}

//根据当前用户id和目标用户id获取关注信息
func (*RelationDao) QueryRelationByIds(currentId int64, userIds []int64) (map[int64]model.RelationRaw, error) {
	var relations []model.RelationRaw
	err := model.DB.Table("relation").Where("user_id = ? AND to_user_id IN ?", currentId, userIds).Find(&relations).Error
	if err != nil {
		logger.Error("query relation by ids " + err.Error())
		return nil, err
	}
	relationMap := make(map[int64]model.RelationRaw)
	for _, relation := range relations {
		relationMap[relation.ToUserId] = relation
	}
	return relationMap, nil
}

//增加当前用户的关注总数，增加其他用户的粉丝总数，创建关注记录
func (*RelationDao) Create(currentId int64, toUserId int64) error {
	relationRaw := &model.RelationRaw{
		Id:       time.Now().Unix(),
		UserId:   currentId,
		ToUserId: toUserId,
	}
	model.DB.Transaction(func(tx *gorm.DB) error {
		err := tx.Table("user").Where("id = ?", currentId).Update("follow_count", gorm.Expr("follow_count + ?", 1)).Error
		if err != nil {
			logger.Error("add user follow_count fail " + err.Error())
			return err
		}

		err = tx.Table("user").Where("id = ?", toUserId).Update("follower_count", gorm.Expr("follower_count + ?", 1)).Error
		if err != nil {
			logger.Error("add user follower_count fail " + err.Error())
			return err
		}

		err = tx.Table("relation").Create(&relationRaw).Error
		if err != nil {
			logger.Error("create relation record fail " + err.Error())
			return err
		}

		return nil
	})
	return nil
}

//减少当前用户的关注总数，减少其他用户的粉丝总数，删除关注记录
func (*RelationDao) Delete(currentId int64, toUserId int64) error {
	var relationRaw *model.RelationRaw
	model.DB.Transaction(func(tx *gorm.DB) error {
		err := tx.Table("user").Where("id = ?", currentId).Update("follow_count", gorm.Expr("follow_count - ?", 1)).Error
		if err != nil {
			logger.Error("sub user follow_count fail " + err.Error())
			return err
		}

		err = tx.Table("user").Where("id = ?", toUserId).Update("follower_count", gorm.Expr("follower_count - ?", 1)).Error
		if err != nil {
			logger.Error("sub user follower_count fail " + err.Error())
			return err
		}

		err = tx.Table("relation").Where("user_id = ? AND to_user_id = ?", currentId, toUserId).Delete(&relationRaw).Error
		if err != nil {
			logger.Error("delete relation record fali " + err.Error())
			return err
		}
		return nil
	})
	return nil
}

//通过用户id，查询该用户关注的用户，返回两者之间的关注记录
func (*RelationDao) QueryFollowById(userId int64) ([]model.RelationRaw, error) {
	var relations []model.RelationRaw
	err := model.DB.Table("relation").Where("user_id = ?", userId).Find(&relations).Error
	if err != nil {
		logger.Error("query follow by id fail " + err.Error())
		return nil, err
	}

	return relations, nil
}

//通过用户id，查询该用户的粉丝， 返回两者之间的关注记录
func (*RelationDao) QueryFollowerById(userId int64) ([]model.RelationRaw, error) {
	var relations []model.RelationRaw
	err := model.DB.Table("relation").Where("to_user_id = ?", userId).Find(&relations).Error
	if err != nil {
		logger.Error("query follower by id fail " + err.Error())
		return nil, err
	}
	return relations, nil
}

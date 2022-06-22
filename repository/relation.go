package repository

import (
	"douyin/model"
	"errors"
	"sync"

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

//根据当前用户id和视频拥有者id获取关注信息
func (*RelationDao) QueryRelationByIds(currentId int64, userIds []int64) (map[int64]*model.RelationRaw, error) {
	var relations []*model.RelationRaw
	err := model.DB.Where("user_id = ? AND to_user_id IN ? AND status IN ?", currentId, userIds, []int64{0, -1}).Or("user_id IN ? AND to_user_id = ? AND status = ?", userIds, currentId, 1).Find(&relations).Error
	if err == gorm.ErrRecordNotFound {
		return nil, errors.New("relation record not found")
	}
	if err != nil {
		return nil, errors.New("query relation record fail")
	}
	relationMap := make(map[int64]*model.RelationRaw)
	for _, relation := range relations {
		if relation.Status == 1 {
			relationMap[relation.UserId] = relation
		} else {
			relationMap[relation.ToUserId] = relation
		}
	}
	return relationMap, nil
}

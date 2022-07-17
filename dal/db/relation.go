package db

import (
	"context"

	"github.com/cloudwego/kitex/pkg/klog"
	"gorm.io/gorm"
)

// Relation Gorm Data Structures
type RelationRaw struct {
	gorm.Model
	UserId   int64 `gorm:"column:user_id;not null;index:idx_userid"`
	ToUserId int64 `gorm:"column:to_user_id;not null;index:idx_touserid"`
}

func (RelationRaw) TableName() string {
	return "relation"
}

//根据当前用户id和目标用户id获取关注信息
func QueryRelationByIds(ctx context.Context, currentId int64, userIds []int64) (map[int64]*RelationRaw, error) {
	var relations []*RelationRaw
	err := DB.WithContext(ctx).Where("user_id = ? AND to_user_id IN ?", currentId, userIds).Find(&relations).Error
	if err != nil {
		klog.Error("query relation by ids " + err.Error())
		return nil, err
	}
	relationMap := make(map[int64]*RelationRaw)
	for _, relation := range relations {
		relationMap[relation.ToUserId] = relation
	}
	return relationMap, nil
}

//增加当前用户的关注总数，增加其他用户的粉丝总数，创建关注记录
func Create(ctx context.Context, currentId int64, toUserId int64) error {
	relationRaw := &RelationRaw{
		UserId:   currentId,
		ToUserId: toUserId,
	}
	DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err := tx.Table("user").Where("id = ?", currentId).Update("follow_count", gorm.Expr("follow_count + ?", 1)).Error
		if err != nil {
			klog.Error("add user follow_count fail " + err.Error())
			return err
		}

		err = tx.Table("user").Where("id = ?", toUserId).Update("follower_count", gorm.Expr("follower_count + ?", 1)).Error
		if err != nil {
			klog.Error("add user follower_count fail " + err.Error())
			return err
		}

		err = tx.Table("relation").Create(&relationRaw).Error
		if err != nil {
			klog.Error("create relation record fail " + err.Error())
			return err
		}

		return nil
	})
	return nil
}

//减少当前用户的关注总数，减少其他用户的粉丝总数，删除关注记录
func Delete(ctx context.Context, currentId int64, toUserId int64) error {
	var relationRaw *RelationRaw
	DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err := tx.Table("user").Where("id = ?", currentId).Update("follow_count", gorm.Expr("follow_count - ?", 1)).Error
		if err != nil {
			klog.Error("sub user follow_count fail " + err.Error())
			return err
		}

		err = tx.Table("user").Where("id = ?", toUserId).Update("follower_count", gorm.Expr("follower_count - ?", 1)).Error
		if err != nil {
			klog.Error("sub user follower_count fail " + err.Error())
			return err
		}

		err = tx.Table("relation").Where("user_id = ? AND to_user_id = ?", currentId, toUserId).Delete(&relationRaw).Error
		if err != nil {
			klog.Error("delete relation record fali " + err.Error())
			return err
		}
		return nil
	})
	return nil
}

//通过用户id，查询该用户关注的用户，返回两者之间的关注记录
func QueryFollowById(ctx context.Context, userId int64) ([]*RelationRaw, error) {
	var relations []*RelationRaw
	err := DB.WithContext(ctx).Table("relation").Where("user_id = ?", userId).Find(&relations).Error
	if err != nil {
		klog.Error("query follow by id fail " + err.Error())
		return nil, err
	}
	return relations, nil
}

//通过用户id，查询该用户的粉丝， 返回两者之间的关注记录
func QueryFollowerById(ctx context.Context, userId int64) ([]*RelationRaw, error) {
	var relations []*RelationRaw
	err := DB.WithContext(ctx).Table("relation").Where("to_user_id = ?", userId).Find(&relations).Error
	if err != nil {
		klog.Error("query follower by id fail " + err.Error())
		return nil, err
	}
	return relations, nil
}

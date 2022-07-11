package db

import (
	"context"

	"gorm.io/gorm"
)

// Comment Gorm Data Structures
type CommentRaw struct {
	gorm.Model
	UserId   int64  `gorm:"column:user_id;not null;index:idx_userid"`
	VideoId  int64  `gorm:"column:video_id;not null;index:idx_videoid"`
	Contents string `grom:"column:contents;type:varchar(255);not null"`
}

func (CommentRaw) TableName() string {
	return "comment"
}

//根据当前用户id和目标用户id获取关注信息
func QueryRelationByIds(ctx context.Context, currentId int64, userIds []int64) (map[int64]*RelationRaw, error) {
	var relations []*RelationRaw
	err := DB.WithContext(ctx).Where("user_id = ? AND to_user_id IN ?", currentId, userIds).Find(&relations).Error
	if err != nil {
		return nil, err
	}
	relationMap := make(map[int64]*RelationRaw)
	for _, relation := range relations {
		relationMap[relation.ToUserId] = relation
	}
	return relationMap, nil
}

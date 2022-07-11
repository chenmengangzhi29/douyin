package db

import (
	"context"

	"gorm.io/gorm"
)

// User Gorm Data structures
type UserRaw struct {
	gorm.Model
	Name          string `gorm:"column:name;index:idx_username,unique;type:varchar(32);not null"`
	Password      string `gorm:"column:password;type:varchar(32);not null"`
	FollowCount   int64  `gorm:"column:follow_count;default:0"`
	FollowerCount int64  `gorm:"column:follower_count;default:0"`
}

func (UserRaw) TableName() string {
	return "user"
}

//根据用户id获取用户信息
func QueryUserByIds(ctx context.Context, userIds []int64) ([]*UserRaw, error) {
	var users []*UserRaw
	err := DB.WithContext(ctx).Where("id in (?)", userIds).Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

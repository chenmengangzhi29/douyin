package db

import "gorm.io/gorm"

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

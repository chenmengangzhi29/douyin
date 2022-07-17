package db

import (
	"context"

	"github.com/cloudwego/kitex/pkg/klog"
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
		klog.Error("query user by ids fail " + err.Error())
		return nil, err
	}
	return users, nil
}

//根据用户名获取用户信息
func QueryUserByName(ctx context.Context, userName string) ([]*UserRaw, error) {
	var users []*UserRaw
	err := DB.WithContext(ctx).Where("name = ?", userName).Find(&users).Error
	if err != nil {
		klog.Error("query user by name fail " + err.Error())
		return nil, err
	}
	return users, nil
}

//上传用户信息到数据库
func UploadUserData(ctx context.Context, username string, password string) (int64, error) {
	user := &UserRaw{
		Name:          username,
		Password:      password,
		FollowCount:   0,
		FollowerCount: 0,
	}
	err := DB.WithContext(ctx).Create(&user).Error
	if err != nil {
		klog.Error("upload user data fail " + err.Error())
		return 0, err
	}
	return int64(user.ID), nil
}

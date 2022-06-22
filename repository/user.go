package repository

import (
	"douyin/model"
	"errors"
	"sync"
)

type UserDao struct {
}

var userDao *UserDao
var userOnce sync.Once

func NewUserDaoInstance() *UserDao {
	userOnce.Do(
		func() {
			userDao = &UserDao{}
		})
	return userDao
}

//根据用户id获取用户信息
func (*UserDao) QueryUserByIds(userIds []int64) (map[int64]*model.UserRaw, error) {
	var users []*model.UserRaw
	err := model.DB.Where("id in (?)", userIds).Find(&users).Error
	if err != nil {
		return nil, errors.New("query user fail")
	}
	userMap := make(map[int64]*model.UserRaw)
	for _, user := range users {
		userMap[user.Id] = user
	}
	return userMap, nil
}

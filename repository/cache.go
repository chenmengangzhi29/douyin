package repository

import (
	"douyin/model"
)

//缓存的用户信息表，存储token到用户的映射
//该缓存数据在服务重新启动自动清除
var usersLoginInfo = map[string]model.UserRaw{
	"JerryJerry123": {
		Id:            1,
		Name:          "Jerry",
		Password:      "Jerry123",
		FollowCount:   0,
		FollowerCount: 0,
		Token:         "JerryJerry123",
	},
}

var userIdSequence = int64(10)

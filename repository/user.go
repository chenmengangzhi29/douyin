package repository

import (
	"douyin/model"
	"douyin/util"
	"errors"
	"sync"
	"sync/atomic"
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
func (*UserDao) QueryUserByIds(userIds []int64) ([]model.UserRaw, error) {
	var users []model.UserRaw
	err := model.DB.Table("user").Where("id in (?)", userIds).Find(&users).Error
	if err != nil {
		util.Logger.Error("query user by ids fail" + err.Error())
		return nil, err
	}
	return users, nil
}

// userMap := make(map[int64]*model.UserRaw)
// for _, user := range users {
// 	userMap[user.Id] = user
// }

//检查用户名是否不存在
func (*UserDao) CheckUserNotExist(username string, password string) error {
	token := username + password
	if _, exist := usersLoginInfo[token]; exist {
		return errors.New("user already exists cache")
	}
	var user *model.UserRaw
	err := model.DB.Table("user").Where("name = ?", username).Find(&user).Error
	if err != nil {
		util.Logger.Error("check user not exist fail" + err.Error())
		return err
	}
	if user.Name == username {
		return errors.New("user already exists database")
	}

	return nil
}

//检查用户名是否存在
func (*UserDao) CheckUserExist(username string, password string) error {
	token := username + password
	if _, exist := usersLoginInfo[token]; exist {
		return nil
	}
	var user *model.UserRaw
	err := model.DB.Table("user").Where("token = ?", token).First(&user).Error
	if err != nil {
		util.Logger.Error("check user exist fail" + err.Error())
		return err
	}
	return nil
}

//上传用户信息到缓存的用户信息表和数据库
func (*UserDao) UploadUserData(username string, password string) (int64, string, error) {
	atomic.AddInt64(&userIdSequence, 1)
	token := username + password
	user := &model.UserRaw{
		Id:            userIdSequence,
		Name:          username,
		Password:      password,
		FollowCount:   0,
		FollowerCount: 0,
		Token:         token,
	}
	usersLoginInfo[token] = *user
	err := model.DB.Table("user").Create(&user).Error
	if err != nil {
		util.Logger.Error("upload user data fail" + err.Error())
		return 0, "", err
	}
	return user.Id, user.Token, nil
}

//通过token获取用户id和用户
func (*UserDao) QueryUserByToken(token string) (int64, string, error) {
	if userInfo, exist := usersLoginInfo[token]; exist {
		return userInfo.Id, userInfo.Token, nil
	}

	var user *model.UserRaw
	err := model.DB.Table("user").Where("token = ?", token).First(&user).Error
	if err != nil {
		util.Logger.Error("query user by token fail" + err.Error())
		return 0, "", err
	}
	return user.Id, user.Token, nil

}

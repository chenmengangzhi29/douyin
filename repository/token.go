package repository

import (
	"douyin/model"
	"douyin/util/logger"
	"errors"
	"sync"

	"gorm.io/gorm"
)

type TokenDao struct {
}

var tokenDao *TokenDao
var tokenOnce sync.Once

func NewTokenDaoInstance() *TokenDao {
	tokenOnce.Do(
		func() {
			tokenDao = &TokenDao{}
		})
	return tokenDao
}

//根据token获取用户id
func (*TokenDao) QueryUserIdByToken(token string) (int64, error) {
	var user *model.UserRaw
	err := model.DB.Table("user").Where("token = ?", token).First(&user).Error
	if err == gorm.ErrRecordNotFound {
		logger.Error("QueryUserIdByToken token not found " + err.Error())
		return -1, errors.New("token not found")
	}
	if err != nil {
		logger.Error("query user id by token fail " + err.Error())
		return -1, err
	}
	return user.Id, nil
}

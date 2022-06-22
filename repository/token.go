package repository

import (
	"douyin/model"
	"errors"
	"sync"
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
	var user model.UserRaw
	err := model.DB.Where("token = ?", token).Find(&user).Error
	if err != nil {
		return -1, errors.New("check token fail")
	}
	return user.Id, nil
}

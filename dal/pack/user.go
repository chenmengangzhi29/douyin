package pack

import (
	"github.com/chenmengangzhi29/douyin/kitex_gen/user"

	"github.com/chenmengangzhi29/douyin/dal/db"
)

func UserInfo(userRaw *db.UserRaw, isFollow bool) *user.User {
	userInfo := &user.User{
		Id:            int64(userRaw.ID),
		Name:          userRaw.Name,
		FollowCount:   userRaw.FollowCount,
		FollowerCount: userRaw.FollowerCount,
		IsFollow:      isFollow,
	}
	return userInfo
}

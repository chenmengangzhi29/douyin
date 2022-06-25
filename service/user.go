package service

import (
	"douyin/model"
	"douyin/repository"
)

//用户注册信息流，通过repository层检查用户名是否存在，若否则上传用户数据，返回用户id和token给handler层
func RegisterUserData(username string, password string) (int64, string, error) {
	return NewRegisterUserDataFlow(username, password).Do()
}

func NewRegisterUserDataFlow(username string, password string) *RegisterUserDataFlow {
	return &RegisterUserDataFlow{
		Username: username,
		Password: password,
	}
}

type RegisterUserDataFlow struct {
	Username string
	Password string

	UserId int64
	Token  string
}

func (f *RegisterUserDataFlow) Do() (int64, string, error) {
	if err := f.checkUserExist(); err != nil {
		return 0, "", err
	}
	if err := f.uploadUserData(); err != nil {
		return 0, "", err
	}
	return f.UserId, f.Token, nil
}

func (f *RegisterUserDataFlow) checkUserExist() error {
	err := repository.NewUserDaoInstance().CheckUserNotExist(f.Username, f.Password)
	if err != nil {
		return err
	}
	return nil
}

func (f *RegisterUserDataFlow) uploadUserData() error {
	userId, token, err := repository.NewUserDaoInstance().UploadUserData(f.Username, f.Password)
	if err != nil {
		return err
	}
	f.UserId = userId
	f.Token = token
	return nil
}

//用户登陆信息流，通过repository层检查用户是否存在，若存在获取用户id和token，返回handler层用户id和token
func LoginUserData(username string, password string) (int64, string, error) {
	return NewLoginUserDataFlow(username, password).Do()
}

func NewLoginUserDataFlow(username string, password string) *LoginUserDataFlow {
	return &LoginUserDataFlow{
		Username: username,
		Password: password,
	}
}

type LoginUserDataFlow struct {
	Username string
	Password string

	UserId int64
	Token  string
}

func (f *LoginUserDataFlow) Do() (int64, string, error) {
	if err := f.checkUserExist(); err != nil {
		return 0, "", err
	}
	if err := f.queryUserData(); err != nil {
		return 0, "", err
	}
	return f.UserId, f.Token, nil
}

func (f *LoginUserDataFlow) checkUserExist() error {
	err := repository.NewUserDaoInstance().CheckUserExist(f.Username, f.Password)
	if err != nil {
		return err
	}
	return nil
}

func (f *LoginUserDataFlow) queryUserData() error {
	token := f.Username + f.Password
	userId, token, err := repository.NewUserDaoInstance().QueryUserByToken(token)
	if err != nil {
		return err
	}
	f.UserId = userId
	f.Token = token
	return nil
}

//获取用户信息流，通过repository层鉴权，查询用户信息，查询关注关系，最后打包返回给handler层
func GetUserInfo(userId int64, token string) (*model.User, error) {
	return NewGetUserInfoFlow(userId, token).Do()
}

func NewGetUserInfoFlow(userId int64, token string) *GetUserInfoFlow {
	return &GetUserInfoFlow{
		UserId: userId,
		Token:  token,
	}
}

type GetUserInfoFlow struct {
	UserId int64
	Token  string
	User   model.User

	IsFollow  bool
	CurrentId int64
	UserRaw   model.UserRaw
}

func (f *GetUserInfoFlow) Do() (*model.User, error) {
	if err := f.checkToken(); err != nil {
		return nil, err
	}
	if err := f.prepareUserInfo(); err != nil {
		return nil, err
	}
	if err := f.packUserInfo(); err != nil {
		return nil, err
	}
	return &f.User, nil
}

func (f *GetUserInfoFlow) checkToken() error {
	if f.Token == "defaultToken" {
		f.CurrentId = -1
		return nil
	}
	currentId, err := repository.NewTokenDaoInstance().QueryUserIdByToken(f.Token)
	if err != nil {
		return err
	}
	f.CurrentId = currentId
	return nil
}

func (f *GetUserInfoFlow) prepareUserInfo() error {
	userIds := []int64{f.UserId}
	users, err := repository.NewUserDaoInstance().QueryUserByIds(userIds)
	if err != nil {
		return err
	}
	f.UserRaw = users[0]

	relationMap, err := repository.NewRelationDaoInstance().QueryRelationByIds(f.CurrentId, userIds)
	// if err == gorm.ErrRecordNotFound {
	// 	f.IsFollow = false
	// } else if err != nil {
	// 	return err
	// } else {
	// 	f.IsFollow = false
	// }
	_, ok := relationMap[f.UserId]
	if ok {
		f.IsFollow = true
	} else {
		f.IsFollow = false
	}

	return nil
}

func (f *GetUserInfoFlow) packUserInfo() error {
	user := model.User{
		Id:            f.UserRaw.Id,
		Name:          f.UserRaw.Name,
		FollowCount:   f.UserRaw.FollowCount,
		FollowerCount: f.UserRaw.FollowerCount,
		IsFollow:      f.IsFollow,
	}
	f.User = user
	return nil
}

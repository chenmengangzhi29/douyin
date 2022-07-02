package service

import (
	"douyin/model"
	"douyin/repository"
	"errors"
)

//关注操作信息流
//
//如果actionType等于1，表示当前用户关注其他用户，
//当前用户的关注总数增加，其他用户的粉丝总数增加，
//新建一条关注记录
//
//如果actionType等于2，表示当前用户取消关注其他用户
//当前用户的关注总数减少，其他用户的粉丝总数减少，
//删除该关注记录
func RelationActionData(token string, toUserId int64, actionType int64) error {
	return NewRelationActionDataFlow(token, toUserId, actionType).Do()
}

func NewRelationActionDataFlow(token string, toUserId int64, actionType int64) *RelationActionDataFlow {
	return &RelationActionDataFlow{
		Token:      token,
		ToUserId:   toUserId,
		ActionType: actionType,
	}
}

type RelationActionDataFlow struct {
	Token      string
	ToUserId   int64
	ActionType int64

	CurrentId int64
}

//如果ActionType等于1，则创建关注记录
//如果ActionType等于2，则删除关注记录
func (f *RelationActionDataFlow) Do() error {
	if err := f.checkToken(); err != nil {
		return err
	}
	if err := f.checkToUserId(); err != nil {
		return err
	}
	if f.ActionType == 1 {
		if err := f.CreateInfo(); err != nil {
			return err
		}
	}
	if f.ActionType == 2 {
		if err := f.DeleteInfo(); err != nil {
			return err
		}
	}
	return nil
}

//鉴权
func (f *RelationActionDataFlow) checkToken() error {
	currentId, err := repository.NewTokenDaoInstance().QueryUserIdByToken(f.Token)
	if err != nil {
		return err
	}
	f.CurrentId = currentId
	return nil
}

func (f *RelationActionDataFlow) checkToUserId() error {
	users, err := repository.NewUserDaoInstance().QueryUserByIds([]int64{f.ToUserId})
	if err != nil {
		return err
	}
	if len(users) == 0 {
		return errors.New("toUserId not exist")
	}
	return nil
}

//增加当前用户的关注总数，增加其他用户的粉丝总数，创建关注记录
func (f *RelationActionDataFlow) CreateInfo() error {
	err := repository.NewRelationDaoInstance().Create(f.CurrentId, f.ToUserId)
	if err != nil {
		return err
	}
	return nil
}

//减少当前用户的关注总数，减少其他用户的粉丝总数，更新或删除关注记录
func (f *RelationActionDataFlow) DeleteInfo() error {
	err := repository.NewRelationDaoInstance().Delete(f.CurrentId, f.ToUserId)
	if err != nil {
		return err
	}
	return nil
}

//关注列表信息流，包括鉴权，获取目标用户关注的用户id号，获取用户id号对应的用户信息，获取当前用户和这些用户的关注信息
func FollowListData(token string, userId int64) ([]*model.User, error) {
	return NewFollowListDataFlow(token, userId).Do()
}

func NewFollowListDataFlow(token string, userId int64) *FollowListDataFlow {
	return &FollowListDataFlow{
		Token:  token,
		UserId: userId,
	}
}

type FollowListDataFlow struct {
	Token    string
	UserId   int64
	UserList []*model.User

	CurrentId   int64
	Users       []*model.UserRaw
	RelationMap map[int64]*model.RelationRaw
}

func (f *FollowListDataFlow) Do() ([]*model.User, error) {
	if err := f.checkToken(); err != nil {
		return nil, err
	}
	if err := f.checkUserId(); err != nil {
		return nil, err
	}
	if err := f.prepareFollowInfo(); err != nil {
		return nil, err
	}
	if err := f.packFollowInfo(); err != nil {
		return nil, err
	}
	return f.UserList, nil
}

//鉴权
func (f *FollowListDataFlow) checkToken() error {
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

//检查用户是否存在
func (f *FollowListDataFlow) checkUserId() error {
	users, err := repository.NewUserDaoInstance().QueryUserByIds([]int64{f.UserId})
	if err != nil {
		return err
	}
	if len(users) == 0 {
		return errors.New("userId not exist")
	}
	return nil
}

func (f *FollowListDataFlow) prepareFollowInfo() error {
	//获取目标用户关注的用户id号
	relations, err := repository.NewRelationDaoInstance().QueryFollowById(f.UserId)
	if err != nil {
		return err
	}
	userIds := make([]int64, 0)
	for _, relation := range relations {
		userIds = append(userIds, relation.ToUserId)
	}

	//获取用户id号对应的用户信息
	users, err := repository.NewUserDaoInstance().QueryUserByIds(userIds)
	if err != nil {
		return err
	}
	f.Users = users

	if f.CurrentId == -1 {
		return nil
	}

	//获取当前用户和这些用户的关注信息
	relationMap, err := repository.NewRelationDaoInstance().QueryRelationByIds(f.CurrentId, userIds)
	if err != nil {
		return err
	}
	f.RelationMap = relationMap

	return nil
}

//打包从repository层获取的数据，返回
func (f *FollowListDataFlow) packFollowInfo() error {
	userList := make([]*model.User, 0)
	for _, user := range f.Users {
		var isFollow bool = false

		if f.CurrentId != -1 {
			_, ok := f.RelationMap[user.Id]
			if ok {
				isFollow = true
			}
		}
		userList = append(userList, &model.User{
			Id:            user.Id,
			Name:          user.Name,
			FollowCount:   user.FollowCount,
			FollowerCount: user.FollowerCount,
			IsFollow:      isFollow,
		})
	}

	f.UserList = userList
	return nil
}

//粉丝列表信息流，包括鉴权，查询目标用户的被关注记录，获取这些记录的关注方id， 获取关注方的信息，获取当前用户与关注方的关注记录
func FollowerListData(token string, userId int64) ([]*model.User, error) {
	return NewFollowerListDataFlow(token, userId).Do()
}

func NewFollowerListDataFlow(token string, userId int64) *FollowerListDataFlow {
	return &FollowerListDataFlow{
		Token:  token,
		UserId: userId,
	}
}

type FollowerListDataFlow struct {
	Token    string
	UserId   int64
	UserList []*model.User

	CurrentId   int64
	Users       []*model.UserRaw
	RelationMap map[int64]*model.RelationRaw
}

func (f *FollowerListDataFlow) Do() ([]*model.User, error) {
	if err := f.checkToken(); err != nil {
		return nil, err
	}
	if err := f.checkUserId(); err != nil {
		return nil, err
	}
	if err := f.prepareFollowerInfo(); err != nil {
		return nil, err
	}
	if err := f.packFollowerInfo(); err != nil {
		return nil, err
	}
	return f.UserList, nil
}

//鉴权
func (f *FollowerListDataFlow) checkToken() error {
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

//检查用户id是否存在
func (f *FollowerListDataFlow) checkUserId() error {
	users, err := repository.NewUserDaoInstance().QueryUserByIds([]int64{f.UserId})
	if err != nil {
		return err
	}
	if len(users) == 0 {
		return errors.New("userId not exist")
	}
	return nil
}

//查询目标用户的被关注记录，获取这些记录的关注方id， 获取关注方的信息，获取当前用户与关注方的关注记录
func (f *FollowerListDataFlow) prepareFollowerInfo() error {
	//查询目标用户的被关注记录
	relations, err := repository.NewRelationDaoInstance().QueryFollowerById(f.UserId)
	if err != nil {
		return err
	}

	//获取这些记录的关注方id
	userIds := make([]int64, 0)
	for _, relation := range relations {
		userIds = append(userIds, relation.UserId)
	}

	//获取关注方的信息
	users, err := repository.NewUserDaoInstance().QueryUserByIds(userIds)
	if err != nil {
		return err
	}
	f.Users = users

	//如果是默认用户，直接返回
	if f.CurrentId == -1 {
		return nil
	}

	//获取当前用户与关注方的关注记录
	relationMap, err := repository.NewRelationDaoInstance().QueryRelationByIds(f.CurrentId, userIds)
	if err != nil {
		return err
	}
	f.RelationMap = relationMap

	return nil
}

//打包repository层返回的数据，返回
func (f *FollowerListDataFlow) packFollowerInfo() error {
	userList := make([]*model.User, 0)
	for _, user := range f.Users {
		var isFollow bool = false

		if f.CurrentId != -1 {
			_, ok := f.RelationMap[user.Id]
			if ok {
				isFollow = true
			}
		}
		userList = append(userList, &model.User{
			Id:            user.Id,
			Name:          user.Name,
			FollowCount:   user.FollowCount,
			FollowerCount: user.FollowerCount,
			IsFollow:      isFollow,
		})
	}

	f.UserList = userList
	return nil
}

package service

import (
	"douyin/model"
	"douyin/repository"
	"errors"
	"fmt"
	"sync"
	"time"
)

//创建评论流，包括鉴权，创建评论记录，增加视频评论数，获取评论用户信息并打包返回
func CreateCommentData(token string, videoId int64, commentText string) (*model.Comment, error) {
	return NewCreateCommentDataFlow(token, videoId, commentText).Do()
}

func NewCreateCommentDataFlow(token string, videoId int64, commentText string) *CreateCommentDataFlow {
	return &CreateCommentDataFlow{
		Token:       token,
		VideoId:     videoId,
		CommentText: commentText,
	}
}

type CreateCommentDataFlow struct {
	Token       string
	VideoId     int64
	CommentText string
	Comment     model.Comment

	CurrentId  int64
	CommentRaw model.CommentRaw
	User       model.UserRaw
}

func (f *CreateCommentDataFlow) Do() (*model.Comment, error) {
	if err := f.checkToken(); err != nil {
		return nil, err
	}

	if err := f.prepareCommentInfo(); err != nil {
		return nil, err
	}

	if err := f.packCommentInfo(); err != nil {
		return nil, err
	}

	return &f.Comment, nil
}

//鉴权
func (f *CreateCommentDataFlow) checkToken() error {
	currentId, err := repository.NewTokenDaoInstance().QueryUserIdByToken(f.Token)
	if err != nil {
		return err
	}
	f.CurrentId = currentId
	return nil
}

func (f *CreateCommentDataFlow) prepareCommentInfo() error {
	commentRaw := model.CommentRaw{
		Id:         time.Now().Unix(),
		UserId:     f.CurrentId,
		VideoId:    f.VideoId,
		Contents:   f.CommentText,
		CreateDate: time.Now().Format("01-01"),
	}
	f.CommentRaw = commentRaw

	var wg sync.WaitGroup
	wg.Add(3)
	var commentErr, videoErr, userErr error
	//创建评论记录
	go func() {
		defer wg.Done()
		err := repository.NewCommentDaoInstance().CreateComment(&f.CommentRaw)
		if err != nil {
			commentErr = err
			return
		}
	}()
	//更改视频评论数
	go func() {
		defer wg.Done()
		err := repository.NewVideoDaoInstance().AddCommentCount(f.VideoId)
		if err != nil {
			videoErr = err
			return
		}
	}()
	//获取当前用户信息
	go func() {
		defer wg.Done()
		userMap, err := repository.NewUserDaoInstance().QueryUserByIds([]int64{f.CurrentId})
		if err != nil {
			userErr = err
			return
		}
		f.User = *userMap[f.CurrentId]
	}()
	wg.Wait()
	if commentErr != nil {
		return commentErr
	}
	if videoErr != nil {
		return videoErr
	}
	if userErr != nil {
		return userErr
	}

	return nil
}

//打包成可以直接返回的评论信息
func (f *CreateCommentDataFlow) packCommentInfo() error {
	comment := model.Comment{
		Id: f.CommentRaw.Id,
		User: model.User{
			Id:            f.User.Id,
			Name:          f.User.Name,
			FollowCount:   f.User.FollowCount,
			FollowerCount: f.User.FollowerCount,
			IsFollow:      false,
		},
		Content:    f.CommentText,
		CreateDate: f.CommentRaw.CreateDate,
	}
	f.Comment = comment
	return nil
}

//删除视频信息流，鉴权，删除评论记录，减少视频评论数，获取当前用户信息和评论信息返回
func DeleteCommentData(token string, videoId int64, commentId int64) (*model.Comment, error) {
	return NewDeleteCommentDataFlow(token, videoId, commentId).Do()
}

func NewDeleteCommentDataFlow(token string, videoId int64, commentId int64) *DeleteCommentDataFlow {
	return &DeleteCommentDataFlow{
		Token:     token,
		VideoId:   videoId,
		CommentId: commentId,
	}
}

type DeleteCommentDataFlow struct {
	Token     string
	VideoId   int64
	CommentId int64
	Comment   model.Comment

	CurrentId  int64
	CommentRaw model.CommentRaw
	UserRaw    model.UserRaw
}

func (f *DeleteCommentDataFlow) Do() (*model.Comment, error) {
	if err := f.checkToken(); err != nil {
		return nil, err
	}
	if err := f.prepareCommentInfo(); err != nil {
		return nil, err
	}
	if err := f.packCommentInfo(); err != nil {
		return nil, err
	}
	return &f.Comment, nil
}

//鉴权
func (f *DeleteCommentDataFlow) checkToken() error {
	currentId, err := repository.NewTokenDaoInstance().QueryUserIdByToken(f.Token)
	if err != nil {
		return err
	}
	f.CurrentId = currentId
	return nil
}

func (f *DeleteCommentDataFlow) prepareCommentInfo() error {
	var wg sync.WaitGroup
	wg.Add(3)
	var commentErr, videoErr, userErr error
	//删除评论记录
	go func() {
		commentRaw, err := repository.NewCommentDaoInstance().DeleteComment(f.CommentId)
		if err != nil {
			commentErr = err
			return
		}
		f.CommentRaw = commentRaw
	}()
	//减少视频评论数
	go func() {
		err := repository.NewVideoDaoInstance().SubCommentCount(f.VideoId)
		if err != nil {
			videoErr = err
			return
		}
	}()
	//获取用户信息
	go func() {
		userMap, err := repository.NewUserDaoInstance().QueryUserByIds([]int64{f.CurrentId})
		if err != nil {
			userErr = err
			return
		}
		f.UserRaw = *userMap[f.CurrentId]
	}()
	wg.Wait()
	if commentErr != nil {
		return commentErr
	}
	if videoErr != nil {
		return videoErr
	}
	if userErr != nil {
		return userErr
	}
	return nil
}

//打包评论信息和用户信息返回
func (f *DeleteCommentDataFlow) packCommentInfo() error {
	comment := model.Comment{
		Id: f.CommentRaw.Id,
		User: model.User{
			Id:            f.UserRaw.Id,
			Name:          f.UserRaw.Name,
			FollowCount:   f.UserRaw.FollowCount,
			FollowerCount: f.UserRaw.FollowerCount,
			IsFollow:      false,
		},
		Content:    f.CommentRaw.Contents,
		CreateDate: f.CommentRaw.CreateDate,
	}
	f.Comment = comment
	return nil
}

//评论列表信息流，包括鉴权，获取一系列评论信息，获取一系列用户信息，获取一系列关注信息，打包返回
func CommentListData(token string, videoId int64) ([]model.Comment, error) {
	return NewCommentListDataFlow(token, videoId).Do()
}

func NewCommentListDataFlow(token string, videoId int64) *CommentListDataFlow {
	return &CommentListDataFlow{
		Token:   token,
		VideoId: videoId,
	}
}

type CommentListDataFlow struct {
	Token       string
	VideoId     int64
	CommentList []model.Comment

	CurrentId   int64
	Comments    []model.CommentRaw
	UserMap     map[int64]*model.UserRaw
	RelationMap map[int64]*model.RelationRaw
}

func (f *CommentListDataFlow) Do() ([]model.Comment, error) {
	if err := f.checkToken(); err != nil {
		return nil, err
	}
	if err := f.prepareCommentInfo(); err != nil {
		return nil, err
	}
	if err := f.packCommentInfo(); err != nil {
		return nil, err
	}
	return f.CommentList, nil
}

//鉴权
func (f *CommentListDataFlow) checkToken() error {
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

func (f *CommentListDataFlow) prepareCommentInfo() error {
	//获取一系列评论信息
	comments, err := repository.NewCommentDaoInstance().QueryCommentByVideoId(f.VideoId)
	if err != nil {
		return err
	}
	f.Comments = comments

	//获取评论信息的用户id
	userIds := make([]int64, 0)
	for _, comment := range f.Comments {
		userIds = append(userIds, comment.UserId)
	}

	//获取一系列用户信息
	userMap, err := repository.NewUserDaoInstance().QueryUserByIds(userIds)
	if err != nil {
		return err
	}
	f.UserMap = userMap

	if f.CurrentId == -1 {
		return nil
	}

	//获取一系列关注信息
	relationMap, err := repository.NewRelationDaoInstance().QueryRelationByIds(f.CurrentId, userIds)
	if err != nil {
		return err
	}
	f.RelationMap = relationMap

	return nil
}

//打包评论信息返回
func (f *CommentListDataFlow) packCommentInfo() error {
	commentList := make([]model.Comment, 0)
	for _, comment := range f.Comments {
		commentUser, ok := f.UserMap[comment.UserId]
		if !ok {
			return errors.New("has no comment user info for " + fmt.Sprint(comment.UserId))
		}

		var isFollow bool = false

		if f.CurrentId != -1 {
			_, ok := f.RelationMap[comment.UserId]
			if ok {
				isFollow = true
			}
		}

		commentList = append(commentList, model.Comment{
			Id: comment.Id,
			User: model.User{
				Id:            commentUser.Id,
				Name:          commentUser.Name,
				FollowCount:   commentUser.FollowCount,
				FollowerCount: commentUser.FollowerCount,
				IsFollow:      isFollow,
			},
			Content:    comment.Contents,
			CreateDate: comment.CreateDate,
		})
	}
	f.CommentList = commentList
	return nil
}

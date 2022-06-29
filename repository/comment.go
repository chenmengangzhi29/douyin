package repository

import (
	"douyin/model"
	"douyin/util/logger"
	"sync"
)

type CommentDao struct {
}

var commentDao *CommentDao
var commentOnce sync.Once

func NewCommentDaoInstance() *CommentDao {
	commentOnce.Do(
		func() {
			commentDao = &CommentDao{}
		})
	return commentDao
}

//通过一条评论创建一条评论记录
func (*CommentDao) CreateComment(comment *model.CommentRaw) error {
	err := model.DB.Table("comment").Create(comment).Error
	if err != nil {
		logger.Error("create comment fail " + err.Error())
		return err
	}
	return nil
}

//通过评论id号删除一条评论，返回该评论
func (*CommentDao) DeleteComment(commentId int64) (model.CommentRaw, error) {
	var commentRaw model.CommentRaw
	err := model.DB.Table("comment").Where("id = ?", commentId).Delete(&commentRaw).Error
	if err != nil {
		logger.Error("delete comment fail " + err.Error())
		return model.CommentRaw{}, err
	}
	return commentRaw, nil
}

//通过视频id号倒序返回一组评论信息
func (*CommentDao) QueryCommentByVideoId(videoId int64) ([]model.CommentRaw, error) {
	var comments []model.CommentRaw
	err := model.DB.Table("comment").Order("create_date desc").Where("video_id = ?", videoId).Find(&comments).Error
	if err != nil {
		logger.Error("query comment by video id fail " + err.Error())
		return nil, err
	}
	return comments, nil
}

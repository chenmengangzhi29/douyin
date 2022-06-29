package repository

import (
	"douyin/model"
	"douyin/util/logger"
	"sync"

	"gorm.io/gorm"
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

//通过一条评论创建一条评论记录并增加视频评论数
func (*CommentDao) CreateComment(comment *model.CommentRaw) error {
	model.DB.Transaction(func(tx *gorm.DB) error {
		err := tx.Table("comment").Create(comment).Error
		if err != nil {
			logger.Error("create comment fail " + err.Error())
			return err
		}
		err = tx.Table("video").Where("id = ?", comment.VideoId).Update("comment_count", gorm.Expr("comment_count + ?", 1)).Error
		if err != nil {
			logger.Error("AddCommentCount error " + err.Error())
			return err
		}
		return nil
	})
	return nil
}

//通过评论id号删除一条评论并减少视频评论数，返回该评论
func (*CommentDao) DeleteComment(commentId int64) (model.CommentRaw, error) {
	var commentRaw model.CommentRaw
	model.DB.Transaction(func(tx *gorm.DB) error {
		err := tx.Table("comment").Where("id = ?", commentId).Delete(&commentRaw).Error
		if err != nil {
			logger.Error("delete comment fail " + err.Error())
			return err
		}
		err = tx.Table("video").Where("id = ?", commentRaw.VideoId).Update("comment_count", gorm.Expr("comment_count - ?", 1)).Error
		if err != nil {
			logger.Error("AddCommentCount error " + err.Error())
			return err
		}
		return nil
	})
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

//通过评论id查询一组评论信息
func (*CommentDao) QueryCommentByCommentIds(commentIds []int64) ([]model.CommentRaw, error) {
	var comments []model.CommentRaw
	err := model.DB.Table("comment").Where("id In ?", commentIds).Find(&comments).Error
	if err != nil {
		logger.Error("query comment by comment id fail " + err.Error())
		return nil, err
	}
	return comments, nil
}

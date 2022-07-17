package db

import (
	"context"

	"github.com/cloudwego/kitex/pkg/klog"
	"gorm.io/gorm"
)

// Comment Gorm Data Structures
type CommentRaw struct {
	gorm.Model
	UserId   int64  `gorm:"column:user_id;not null;index:idx_userid"`
	VideoId  int64  `gorm:"column:video_id;not null;index:idx_videoid"`
	Contents string `grom:"column:contents;type:varchar(255);not null"`
}

func (CommentRaw) TableName() string {
	return "comment"
}

//通过一条评论创建一条评论记录并增加视频评论数
func CreateComment(ctx context.Context, comment *CommentRaw) error {
	DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err := tx.Table("comment").Create(&comment).Error
		if err != nil {
			klog.Error("create comment fail " + err.Error())
			return err
		}
		err = tx.Table("video").Where("id = ?", comment.VideoId).Update("comment_count", gorm.Expr("comment_count + ?", 1)).Error
		if err != nil {
			klog.Error("AddCommentCount error " + err.Error())
			return err
		}
		err = tx.Table("comment").First(&comment).Error
		if err != nil {
			klog.Errorf("find comment %v fail, %v", comment, err.Error())
			return err
		}
		return nil
	})
	return nil
}

//通过评论id号删除一条评论并减少视频评论数，返回该评论
func DeleteComment(ctx context.Context, commentId int64) (*CommentRaw, error) {
	var commentRaw *CommentRaw
	DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err := tx.Table("comment").Where("id = ?", commentId).First(&commentRaw).Error
		if err == gorm.ErrRecordNotFound {
			klog.Errorf("not find comment %v, %v", commentRaw, err.Error())
			return err
		}
		if err != nil {
			klog.Errorf("find comment %v fail, %v", commentRaw, err.Error())
			return err
		}
		err = tx.Table("comment").Where("id = ?", commentId).Delete(&CommentRaw{}).Error
		if err != nil {
			klog.Error("delete comment fail " + err.Error())
			return err
		}
		err = tx.Table("video").Where("id = ?", commentRaw.VideoId).Update("comment_count", gorm.Expr("comment_count - ?", 1)).Error
		if err != nil {
			klog.Error("AddCommentCount error " + err.Error())
			return err
		}
		return nil
	})
	return commentRaw, nil
}

//通过评论id查询一组评论信息
func QueryCommentByCommentIds(ctx context.Context, commentIds []int64) ([]*CommentRaw, error) {
	var comments []*CommentRaw
	err := DB.WithContext(ctx).Table("comment").Where("id In ?", commentIds).Find(&comments).Error
	if err != nil {
		klog.Error("query comment by comment id fail " + err.Error())
		return nil, err
	}
	return comments, nil
}

//通过视频id号倒序返回一组评论信息
func QueryCommentByVideoId(ctx context.Context, videoId int64) ([]*CommentRaw, error) {
	var comments []*CommentRaw
	err := DB.WithContext(ctx).Table("comment").Order("updated_at desc").Where("video_id = ?", videoId).Find(&comments).Error
	if err != nil {
		klog.Error("query comment by video id fail " + err.Error())
		return nil, err
	}
	return comments, nil
}

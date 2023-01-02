package repository

import (
	"aweme_kitex/models"
	"aweme_kitex/utils"
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

// 通过一条评论创建一条评论记录并增加视频评论数
func (*CommentDao) CreateComment(comment *models.CommentRaw) error {
	DB.Transaction(func(tx *gorm.DB) error {
		err := tx.Table("comment").Create(comment).Error
		if err != nil {
			utils.Error("create comment fail " + err.Error())
			return err
		}
		err = tx.Table("video").Where("video_id = ?", comment.VideoId).Update("comment_count", gorm.Expr("comment_count + ?", 1)).Error
		if err != nil {
			utils.Error("AddCommentCount error " + err.Error())
			return err
		}
		return nil
	})
	return nil
}

// 通过评论id号删除一条评论，返回该评论
func (*CommentDao) DeleteComment(commentId string) (*models.CommentRaw, error) {
	var commentRaw *models.CommentRaw
	DB.Transaction(func(tx *gorm.DB) error {
		err := tx.Table("comment").Where("id = ?", commentId).First(&commentRaw).Error
		if err == gorm.ErrRecordNotFound {
			utils.Errorf("not find comment %v, %v", commentRaw, err.Error())
			return err
		}
		if err != nil {
			utils.Error("delete comment fail " + err.Error())
			return err
		}
		err = tx.Table("comment").Where("comment_id = ?", commentId).Delete(&commentRaw).Error
		if err != nil {
			utils.Error("delete comment fail " + err.Error())
			return err
		}
		err = tx.Table("video").Where("video_id = ?", commentRaw.VideoId).Update("comment_count", gorm.Expr("comment_count - ?", 1)).Error
		if err != nil {
			utils.Error("DelCommentCount error " + err.Error())
			return err
		}
		return nil
	})
	return commentRaw, nil
}

// 通过评论id查询一组评论信息
func (*CommentDao) QueryCommentByCommentIds(commentIds []string) ([]*models.CommentRaw, error) {
	var comments []*models.CommentRaw
	err := DB.Table("comment").Where("comment_id In ?", commentIds).Find(&comments).Error
	if err != nil {
		utils.Error("query comment by comment id fail " + err.Error())
		return nil, err
	}
	return comments, nil
}

// 通过视频id号倒序返回一组评论信息
func (*CommentDao) QueryCommentByVideoId(videoId string) ([]*models.CommentRaw, error) {
	var comments []*models.CommentRaw
	err := DB.Debug().Table("comment").Order("created_at desc").Where("video_id = ?", videoId).Find(&comments).Error
	if err != nil {
		utils.Error("query comment err: " + err.Error())
		return nil, err
	}
	return comments, nil
}

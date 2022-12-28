package repository

import (
	"aweme_kitex/models"
	"aweme_kitex/utils"
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

// 通过一条评论创建一条评论记录
func (*CommentDao) CreateComment(comment *models.CommentRaw) error {
	err := DB.Table("comment").Create(comment).Error
	if err != nil {
		utils.Error("create comment err: " + err.Error())
		return err
	}
	return nil
}

// 通过评论id号删除一条评论，返回该评论
func (*CommentDao) DeleteComment(commentId string) (*models.CommentRaw, error) {
	var commentRaw *models.CommentRaw
	err := DB.Debug().Table("comment").Where("comment_id = ?", commentId).Delete(&commentRaw).Error
	if err != nil {
		utils.Error("delete comment err: " + err.Error())
		return &models.CommentRaw{}, err
	}
	return commentRaw, nil
}

// 通过视频id号倒序返回一组评论信息
func (*CommentDao) QueryCommentByVideoId(videoId string) ([]models.CommentRaw, error) {
	var comments []models.CommentRaw
	err := DB.Debug().Table("comment").Order("created_at desc").Where("video_id = ?", videoId).Find(&comments).Error
	if err != nil {
		utils.Error("query comment err: " + err.Error())
		return nil, err
	}
	return comments, nil
}

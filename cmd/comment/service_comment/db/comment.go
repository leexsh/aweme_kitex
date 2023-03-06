package db

import (
	"aweme_kitex/cfg"
	"aweme_kitex/models"
	"aweme_kitex/pkg/logger"
	"context"
	"errors"
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

// 通过一条评论创建一条评论记录并增加视频评论数
func (*CommentDao) CreateComment(ctx context.Context, comment *models.CommentRaw) error {
	return cfg.DB.WithContext(ctx).Table("comment").Create(comment).Error
}

// 通过评论id号删除一条评论，返回该评论
func (*CommentDao) DeleteComment(ctx context.Context, commentId string) (*models.CommentRaw, error) {
	var commentRaw *models.CommentRaw
	err := cfg.DB.WithContext(ctx).Table("comment").Where("id = ?", commentId).First(&commentRaw).Error
	if err != nil {
		return nil, err
	}
	err = cfg.DB.WithContext(ctx).Table("comment").Delete(&commentRaw).Error
	return commentRaw, err
}

// 通过评论id查询一组评论信息
func (*CommentDao) QueryCommentByCommentIds(ctx context.Context, commentIds []string) ([]*models.CommentRaw, error) {
	var comments []*models.CommentRaw
	err := cfg.DB.WithContext(ctx).Table("comment").Where("comment_id In ?", commentIds).Find(&comments).Error
	if err != nil {
		logger.Error("query comment by comment id fail " + err.Error())
		return nil, err
	}
	return comments, nil
}

// 通过视频id号倒序返回一组评论信息
func (*CommentDao) QueryCommentByVideoId(ctx context.Context, videoId string) ([]*models.CommentRaw, error) {
	var comments []*models.CommentRaw
	err := cfg.DB.Debug().WithContext(ctx).Table("comment").Order("created_at desc").Where("video_id = ?", videoId).Find(&comments).Error
	if err != nil {
		logger.Error("query comment err: " + err.Error())
		return nil, err
	}
	return comments, nil
}

// 检查commentid
func (c *CommentDao) CheckCommentId(ctx context.Context, commentIds []string) ([]*models.CommentRaw, error) {
	comments, err := c.QueryCommentByCommentIds(ctx, commentIds)
	if err != nil {
		return nil, err
	}
	if len(comments) == 0 {
		return nil, errors.New("commentId not exist")
	}
	return comments, nil
}

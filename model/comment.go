package model

import "sync"

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
func (*CommentDao) CreateComment(comment *CommentRaw) error {
	err := db.Table("comment").Create(comment).Error
	if err != nil {
		return err
	}
	return nil
}

// 通过评论id号删除一条评论，返回该评论
func (*CommentDao) DeleteComment(commentId string) (*CommentRaw, error) {
	var commentRaw *CommentRaw
	err := db.Debug().Table("comment").Where("comment_id = ?", commentId).Delete(&commentRaw).Error
	if err != nil {
		return &CommentRaw{}, err
	}
	return commentRaw, nil
}

// 通过视频id号倒序返回一组评论信息
func (*CommentDao) QueryCommentByVideoId(videoId string) ([]CommentRaw, error) {
	var comments []CommentRaw
	err := db.Debug().Table("comment").Order("created_at desc").Where("video_id = ?", videoId).Find(&comments).Error
	if err != nil {
		return nil, err
	}
	return comments, nil
}

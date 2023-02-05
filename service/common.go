package service

import (
	"aweme_kitex/models"
	"aweme_kitex/models/dal"
	"context"
	"errors"
)

// 检查视频id是否存在
func checkVideoId(videoId []string) ([]*models.VideoRawData, error) {
	videos, err := dal.NewVideoDaoInstance().QueryVideosByIs(context.Background(), videoId)
	if err != nil {
		return nil, err
	}
	if len(videos) == 0 {
		return nil, errors.New("video not exist")
	}
	return videos, nil
}

// 检查commentid
func checkCommentId(commentIds []string) ([]*models.CommentRaw, error) {
	comments, err := dal.NewCommentDaoInstance().QueryCommentByCommentIds(context.Background(), commentIds)
	if err != nil {
		return nil, err
	}
	if len(comments) == 0 {
		return nil, errors.New("commentId not exist")
	}
	return comments, nil
}

// 检查用户是否存在
func checkUserId(uids []string) ([]*models.UserRawData, error) {
	users, err := dal.NewUserDaoInstance().QueryUserByIds(context.Background(), uids)
	if err != nil {
		return nil, err
	}
	if len(users) == 0 {
		return nil, errors.New("userId not exist")
	}
	return users, nil
}

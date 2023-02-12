package service_comment

import (
	"aweme_kitex/cmd/comment/kitex_gen/comment"
	"aweme_kitex/cmd/comment/kitex_gen/user"
	"aweme_kitex/models"
	"aweme_kitex/models/dal"
	"aweme_kitex/utils"
	"context"
	"sync"
)

func createComment(ctx context.Context, userId, videoId, content, commitId string) (*comment.Comment, error) {
	return newCommentDataFlow(ctx, userId, videoId, content, commitId).createComment()
}

func delComment(ctx context.Context, userId, commentId string) (*comment.Comment, error) {
	return newCommentDataFlow(ctx, userId, "", "", commentId).delComment()
}

type commentDataFlow struct {
	ctx        context.Context
	currentUid string

	content    string
	videoId    string
	commentId  string
	comment    *comment.Comment
	user       *models.UserRawData
	commentRaw *models.CommentRaw
}

func newCommentDataFlow(ctx context.Context, id, videoId, content, commentId string) *commentDataFlow {
	return &commentDataFlow{
		ctx:        ctx,
		currentUid: id,
		videoId:    videoId,
		content:    content,
		commentId:  commentId,
	}
}

func (c *commentDataFlow) createComment() (*comment.Comment, error) {
	if _, err := dal.NewVideoDaoInstance().CheckVideoId(c.ctx, []string{c.videoId}); err != nil {
		return nil, err
	}
	if err := c.prepareComment("1"); err != nil {
		return nil, err
	}
	if err := c.packageComment(); err != nil {
		return nil, err
	}
	return c.comment, nil
}

func (c *commentDataFlow) delComment() (*comment.Comment, error) {
	if _, err := dal.NewCommentDaoInstance().CheckCommentId(c.ctx, []string{c.commentId}); err != nil {
		return nil, err
	}
	if _, err := dal.NewVideoDaoInstance().CheckVideoId(c.ctx, []string{c.videoId}); err != nil {
		return nil, err
	}
	if err := c.prepareComment("2"); err != nil {
		return nil, err
	}
	commet := &comment.Comment{
		CommentId: c.commentRaw.Id,
		User: &user.User{
			UserId:        c.user.UserId,
			Name:          c.user.Name,
			FollowerCount: c.user.FollowerCount,
			FollowCount:   c.user.FollowCount,
			IsFollow:      false,
		},
		Content:    c.commentRaw.Content,
		CreateTime: utils.TimeToString(c.commentRaw.CreatedAt),
	}
	c.comment = commet
	return c.comment, nil
}

func (c *commentDataFlow) prepareComment(action string) error {
	commentRaw := &models.CommentRaw{
		Id:      utils.GenerateUUID(),
		UserId:  c.currentUid,
		VideoId: c.videoId,
		Content: c.content,
	}
	c.commentRaw = commentRaw
	var wg sync.WaitGroup
	wg.Add(2)
	var commentErr, userErr error
	go func() {
		defer wg.Done()
		if action == "1" {
			err := dal.NewCommentDaoInstance().CreateComment(c.ctx, commentRaw)
			if err != nil {
				commentErr = err
			}
		} else if action == "2" {
			comment, err := dal.NewCommentDaoInstance().DeleteComment(c.ctx, c.commentId)
			if err != nil {
				commentErr = err
			}
			c.commentRaw = comment
		}
	}()
	go func() {
		defer wg.Done()
		user, err := dal.NewUserDaoInstance().QueryUserByUserId(c.ctx, c.currentUid)
		if err != nil {
			userErr = err
		}
		c.user = user
	}()
	wg.Wait()
	if commentErr != nil {
		return commentErr
	}
	if userErr != nil {
		return userErr
	}
	return nil
}

func (c *commentDataFlow) packageComment() error {
	commet := &comment.Comment{
		CommentId: c.commentRaw.Id,
		User: &user.User{
			UserId:        c.user.UserId,
			Name:          c.user.Name,
			FollowerCount: c.user.FollowerCount,
			FollowCount:   c.user.FollowCount,
			IsFollow:      false,
		},
		Content:    c.content,
		CreateTime: utils.TimeToString(c.commentRaw.CreatedAt),
	}
	c.comment = commet
	return nil
}

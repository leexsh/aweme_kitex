package service

import (
	"aweme_kitex/models"
	"aweme_kitex/models/dal"
	"aweme_kitex/utils"
	"context"
	"errors"
	"fmt"
	"sync"
)

func CreateComment(userId, videoId, content, commitId string) (*models.Comment, error) {
	return newCommentDataFlow(userId, videoId, content, commitId).createComment()
}

func DelComment(userId, commentId string) (*models.Comment, error) {
	return newCommentDataFlow(userId, "", "", commentId).delComment()
}

type commentDataFlow struct {
	currentUid string

	content    string
	videoId    string
	commentId  string
	comment    *models.Comment
	user       *models.UserRawData
	commentRaw *models.CommentRaw
}

func newCommentDataFlow(id, videoId, content, commentId string) *commentDataFlow {
	return &commentDataFlow{
		currentUid: id,
		videoId:    videoId,
		content:    content,
		commentId:  commentId,
	}
}

func (c *commentDataFlow) createComment() (*models.Comment, error) {
	if _, err := checkVideoId([]string{c.videoId}); err != nil {
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
func (c *commentDataFlow) delComment() (*models.Comment, error) {
	if _, err := checkCommentId([]string{c.commentId}); err != nil {
		return nil, err
	}
	if _, err := checkVideoId([]string{c.videoId}); err != nil {
		return nil, err
	}
	if err := c.prepareComment("2"); err != nil {
		return nil, err
	}
	commet := &models.Comment{
		Id: c.commentRaw.Id,
		User: &models.User{
			UserId:        c.user.UserId,
			Name:          c.user.Name,
			FollowerCount: c.user.FollowerCount,
			FollowCount:   c.user.FollowCount,
			IsFollow:      false,
		},
		VideoId:    c.videoId,
		Content:    c.commentRaw.Content,
		CreateDate: utils.TimeToString(c.commentRaw.CreatedAt),
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
			err := dal.NewCommentDaoInstance().CreateComment(context.Background(), commentRaw)
			if err != nil {
				commentErr = err
			}
		} else if action == "2" {
			comment, err := dal.NewCommentDaoInstance().DeleteComment(context.Background(), c.commentId)
			if err != nil {
				commentErr = err
			}
			c.commentRaw = comment
		}
	}()
	go func() {
		defer wg.Done()
		user, err := dal.NewUserDaoInstance().QueryUserByUserId(context.Background(), c.currentUid)
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
	commet := &models.Comment{
		Id: c.commentRaw.Id,
		User: &models.User{
			UserId:        c.user.UserId,
			Name:          c.user.Name,
			FollowerCount: c.user.FollowerCount,
			FollowCount:   c.user.FollowCount,
			IsFollow:      false,
		},
		VideoId:    c.videoId,
		Content:    c.content,
		CreateDate: utils.TimeToString(c.commentRaw.CreatedAt),
	}
	c.comment = commet
	return nil
}

// ---------------------------
func ShowCommentList(uid, videoId string) ([]*models.Comment, error) {
	return newCommentListDataFlow(uid, videoId).do()
}

type commentListDataFlow struct {
	VideoId     string
	CommentList []*models.Comment

	userId      string
	Comments    []*models.CommentRaw
	UserMap     map[string]*models.UserRawData
	RelationMap map[string]*models.RelationRaw
}

func newCommentListDataFlow(uid string, videoId string) *commentListDataFlow {
	return &commentListDataFlow{
		userId:  uid,
		VideoId: videoId,
	}
}

func (c *commentListDataFlow) do() ([]*models.Comment, error) {
	if err := c.checkVideoId(); err != nil {
		return nil, err
	}
	if err := c.prepareListCommentInfo(); err != nil {
		return nil, err
	}
	if err := c.packCommentListInfo(); err != nil {
		return nil, err
	}
	return c.CommentList, nil
}

// 检查视频id是否正确
func (f *commentListDataFlow) checkVideoId() error {
	videos, err := dal.NewVideoDaoInstance().QueryVideosByIs(context.Background(), []string{f.VideoId})
	if err != nil {
		return err
	}
	if len(videos) == 0 {
		return errors.New("videoId not exist")
	}
	return nil
}

func (c *commentListDataFlow) prepareListCommentInfo() error {
	// 获取一系列评论信息
	comments, err := dal.NewCommentDaoInstance().QueryCommentByVideoId(context.Background(), c.VideoId)
	if err != nil {
		return err
	}
	c.Comments = comments

	// 获取评论信息的用户id
	userIds := make([]string, 0)
	for _, comment := range c.Comments {
		userIds = append(userIds, comment.UserId)
	}

	// 获取一系列用户信息
	users, err := dal.NewUserDaoInstance().QueryUserByIds(context.Background(), userIds)
	if err != nil {
		return err
	}
	userMap := make(map[string]*models.UserRawData)
	for _, user := range users {
		userMap[user.UserId] = user
	}
	c.UserMap = userMap

	// 获取一系列关注信息
	relationMap, err := dal.NewRelationDaoInstance().QueryRelationByIds(context.Background(), c.userId, userIds)
	if err != nil {
		return err
	}
	c.RelationMap = relationMap

	return nil
}

// 打包评论信息返回
func (c *commentListDataFlow) packCommentListInfo() error {
	commentList := make([]*models.Comment, 0)
	for _, comment := range c.Comments {
		commentUser, ok := c.UserMap[comment.UserId]
		if !ok {
			return errors.New("has no comment user info for " + fmt.Sprint(comment.UserId))
		}

		var isFollow bool = false
		_, ok = c.RelationMap[comment.UserId]
		if ok {
			isFollow = true
		}

		commentList = append(commentList, &models.Comment{
			Id: comment.Id,
			User: &models.User{
				UserId:        commentUser.UserId,
				Name:          commentUser.Name,
				FollowCount:   commentUser.FollowCount,
				FollowerCount: commentUser.FollowerCount,
				IsFollow:      isFollow,
			},
			VideoId:    comment.VideoId,
			Content:    comment.Content,
			CreateDate: utils.TimeToString(comment.CreatedAt),
		})
	}
	c.CommentList = commentList
	return nil
}

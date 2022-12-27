package service

import (
	"aweme_kitex/models"
	"aweme_kitex/models/repository"
	"aweme_kitex/utils"
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
	if err := c.prepareComment("1"); err != nil {
		return nil, err
	}
	if err := c.packageComment(); err != nil {
		return nil, err
	}
	return c.comment, nil
}
func (c *commentDataFlow) delComment() (*models.Comment, error) {
	if err := c.prepareComment("2"); err != nil {
		return nil, err
	}
	commet := &models.Comment{
		Id: c.commentRaw.Id,
		User: models.User{
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
	wg.Add(3)
	var commentErr, videoErr, userErr error
	go func() {
		defer wg.Done()
		if action == "1" {
			err := repository.NewCommentDaoInstance().CreateComment(commentRaw)
			if err != nil {
				commentErr = err
			}
		} else if action == "2" {
			comment, err := repository.NewCommentDaoInstance().DeleteComment(c.commentId)
			if err != nil {
				commentErr = err
			}
			c.commentRaw = comment
		}
	}()
	go func() {
		defer wg.Done()
		err := repository.NewVideoDaoInstance().UpdateCommentCount(c.videoId, action)
		if err != nil {
			videoErr = err
		}
	}()
	go func() {
		defer wg.Done()
		user, err := repository.NewUserDaoInstance().QueryUserByUserId(c.currentUid)
		if err != nil {
			userErr = err
		}
		c.user = user
	}()
	wg.Wait()
	if commentErr != nil {
		return commentErr
	}
	if videoErr != nil {
		return videoErr
	}
	if userErr != nil {
		return userErr
	}
	return nil
}

func (c *commentDataFlow) packageComment() error {
	commet := &models.Comment{
		Id: c.commentRaw.Id,
		User: models.User{
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

func ShowCommentList(uid, videoId string) ([]models.Comment, error) {
	return newCommentListDataFlow(uid, videoId).do()
}

type commentListDataFlow struct {
	VideoId     string
	CommentList []models.Comment

	userId      string
	Comments    []models.CommentRaw
	UserMap     map[string]*models.UserRawData
	RelationMap map[string]*models.RelationRaw
}

func newCommentListDataFlow(uid string, videoId string) *commentListDataFlow {
	return &commentListDataFlow{
		userId:  uid,
		VideoId: videoId,
	}
}

func (c *commentListDataFlow) do() ([]models.Comment, error) {
	if err := c.prepareListCommentInfo(); err != nil {
		return nil, err
	}
	if err := c.packCommentListInfo(); err != nil {
		return nil, err
	}
	return c.CommentList, nil
}

func (c *commentListDataFlow) prepareListCommentInfo() error {
	// 获取一系列评论信息
	comments, err := repository.NewCommentDaoInstance().QueryCommentByVideoId(c.VideoId)
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
	users, err := repository.NewUserDaoInstance().QueryUserByIds(userIds)
	if err != nil {
		return err
	}
	userMap := make(map[string]*models.UserRawData)
	for _, user := range users {
		userMap[user.UserId] = user
	}
	c.UserMap = userMap

	// 获取一系列关注信息
	relationMap, err := repository.NewRelationDaoInstance().QueryRelationByIds(c.userId, userIds)
	if err != nil {
		return err
	}
	c.RelationMap = relationMap

	return nil
}

// 打包评论信息返回
func (c *commentListDataFlow) packCommentListInfo() error {
	commentList := make([]models.Comment, 0)
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

		commentList = append(commentList, models.Comment{
			Id: comment.Id,
			User: models.User{
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

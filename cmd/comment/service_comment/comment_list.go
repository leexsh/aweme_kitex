package service_comment

import (
	"aweme_kitex/cmd/comment/kitex_gen/comment"
	"aweme_kitex/cmd/comment/kitex_gen/user"
	"aweme_kitex/models"
	"aweme_kitex/models/dal"
	"aweme_kitex/pkg/jwt"
	"aweme_kitex/utils"
	"context"
	"errors"
	"fmt"
)

type CommentListService struct {
	ctx context.Context
}

// NewCommentListService new CommentListService
func NewCommentListService(ctx context.Context) *CommentListService {
	return &CommentListService{ctx: ctx}
}
func (s *CommentListService) CommentList(req *comment.CommentListRequest) ([]*comment.Comment, error) {
	uc, err := jwt.AnalyzeToken(req.Token)
	if err != nil {
		return nil, err
	}
	comments, err := showCommentList(s.ctx, uc.Id, req.VideoId)
	if err != nil {
		return nil, err
	}
	return comments, nil
}

// ---------------------------
func showCommentList(ctx context.Context, uid, videoId string) ([]*comment.Comment, error) {
	return newCommentListDataFlow(ctx, uid, videoId).do()
}

type commentListDataFlow struct {
	ctx         context.Context
	VideoId     string
	CommentList []*comment.Comment

	userId      string
	Comments    []*models.CommentRaw
	UserMap     map[string]*models.UserRawData
	RelationMap map[string]*models.RelationRaw
}

func newCommentListDataFlow(ctx context.Context, uid string, videoId string) *commentListDataFlow {
	return &commentListDataFlow{
		ctx:     ctx,
		VideoId: videoId,
		userId:  uid,
	}
}

func (c *commentListDataFlow) do() ([]*comment.Comment, error) {
	if _, err := dal.NewVideoDaoInstance().CheckVideoId(c.ctx, []string{c.VideoId}); err != nil {
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

func (c *commentListDataFlow) prepareListCommentInfo() error {
	// 获取一系列评论信息
	comments, err := dal.NewCommentDaoInstance().QueryCommentByVideoId(c.ctx, c.VideoId)
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
	users, err := dal.NewUserDaoInstance().QueryUserByIds(c.ctx, userIds)
	if err != nil {
		return err
	}
	userMap := make(map[string]*models.UserRawData)
	for _, user := range users {
		userMap[user.UserId] = user
	}
	c.UserMap = userMap

	// 获取一系列关注信息
	relationMap, err := dal.NewRelationDaoInstance().QueryRelationByIds(c.ctx, c.userId, userIds)
	if err != nil {
		return err
	}
	c.RelationMap = relationMap

	return nil
}

// 打包评论信息返回
func (c *commentListDataFlow) packCommentListInfo() error {
	commentList := make([]*comment.Comment, 0)
	for _, comm := range c.Comments {
		commentUser, ok := c.UserMap[comm.UserId]
		if !ok {
			return errors.New("has no comment user info for " + fmt.Sprint(comm.UserId))
		}

		var isFollow bool = false
		_, ok = c.RelationMap[comm.UserId]
		if ok {
			isFollow = true
		}
		curComment := &comment.Comment{
			CommentId: comm.Id,
			User: &user.User{
				UserId:        commentUser.UserId,
				Name:          commentUser.Name,
				FollowCount:   commentUser.FollowCount,
				FollowerCount: commentUser.FollowerCount,
				IsFollow:      isFollow,
			},
			Content:    comm.Content,
			CreateTime: utils.TimeToString(comm.CreatedAt),
		}
		commentList = append(commentList, curComment)
	}
	c.CommentList = commentList
	return nil
}

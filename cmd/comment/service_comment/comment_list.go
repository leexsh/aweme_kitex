package service_comment

import (
	"aweme_kitex/cmd/comment/kitex_gen/comment"
	"aweme_kitex/cmd/comment/kitex_gen/user"
	commentRPC "aweme_kitex/cmd/comment/rpc"
	"aweme_kitex/cmd/comment/service_comment/db"
	"aweme_kitex/cmd/feed/kitex_gen/feed"
	"aweme_kitex/cmd/relation/kitex_gen/relation"
	user2 "aweme_kitex/cmd/user/kitex_gen/user"
	"aweme_kitex/pkg/jwt"
	"aweme_kitex/pkg/types"
	"aweme_kitex/pkg/utils"
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
	Comments    []*types.CommentRaw
	UserMap     map[string]*user2.User
	RelationMap map[string]bool
}

func newCommentListDataFlow(ctx context.Context, uid string, videoId string) *commentListDataFlow {
	return &commentListDataFlow{
		ctx:     ctx,
		VideoId: videoId,
		userId:  uid,
	}
}

func (c *commentListDataFlow) do() ([]*comment.Comment, error) {
	err := commentRPC.CheckVideoInvalid(c.ctx, &feed.CheckVideoInvalidRequest{VideoId: []string{c.VideoId}})
	if err != nil {
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
	comments, err := db.NewCommentDaoInstance().QueryCommentByVideoId(c.ctx, c.VideoId)
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
	users, err := commentRPC.GetUserInfo(c.ctx, &user2.SingleUserInfoRequest{UserIds: userIds})
	if err != nil {
		return err
	}
	userMap := make(map[string]*user2.User)
	for _, user := range users {
		userMap[user.UserId] = user
	}
	c.UserMap = userMap

	// 获取一系列关注信息
	relationMap := make(map[string]bool, len(userIds))
	for _, id := range userIds {
		res, err := commentRPC.QueryRelation(c.ctx, &relation.QueryRelationRequest{
			UserId:   c.userId,
			ToUserId: id,
			IsFollow: false,
		})
		if err != nil {
			return err
		}
		relationMap[id] = res
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

package service_comment

import (
	"aweme_kitex/cmd/comment/kitex_gen/comment"
	"aweme_kitex/cmd/comment/kitex_gen/user"
	commentRPC "aweme_kitex/cmd/comment/rpc"
	"aweme_kitex/cmd/comment/service_comment/db"
	"aweme_kitex/cmd/feed/kitex_gen/feed"
	relation2 "aweme_kitex/cmd/relation/kitex_gen/relation"
	user2 "aweme_kitex/cmd/user/kitex_gen/user"
	"aweme_kitex/pkg/jwt"
	"aweme_kitex/pkg/utils"
	"context"
	"errors"
)

type DeleteCommentService struct {
	ctx context.Context
}

// NewDeleteCommentService new DeleteCommentService
func NewDeleteCommentService(ctx context.Context) *DeleteCommentService {
	return &DeleteCommentService{ctx: ctx}
}

func (s *DeleteCommentService) DelComment(req *comment.CommentActionRequest) (*comment.Comment, error) {
	uc, err := jwt.AnalyzeToken(req.Token)
	if err != nil {
		return nil, err
	}
	return s.do(uc.Id, req.CommentId, req.VideoId)
}

func (s *DeleteCommentService) do(uid, cid, vid string) (*comment.Comment, error) {
	// 1. delete comment
	delComm, err := db.NewCommentDaoInstance().DeleteComment(s.ctx, cid)
	if err != nil {
		return nil, err
	}
	// 2. use video rpc to sub comment count
	err = commentRPC.ChangeCommentCount(s.ctx, &feed.ChangeCommentCountRequest{
		VideoId: vid,
		Action:  1,
	})
	if err != nil {
		return nil, err
	}
	// 3. use relationRPC get relation
	videos, err := commentRPC.GetVideosById(s.ctx, &feed.CheckVideoInvalidRequest{VideoId: []string{vid}})
	if err != nil || len(videos) <= 0 {
		return nil, err
	}
	relation, err := commentRPC.QueryRelation(s.ctx, &relation2.QueryRelationRequest{
		UserId:   uid,
		ToUserId: videos[0].Author.UserId,
		IsFollow: false,
	})
	// 4.user rpc get user info
	resp, err := commentRPC.GetUserInfo(s.ctx, &user2.SingleUserInfoRequest{UserIds: []string{uid}})
	if err != nil {
		return nil, err
	}
	if len(resp) == 0 {
		return nil, errors.New("not found this user")
	}
	us := resp[uid]
	comm := &comment.Comment{
		CommentId: delComm.Id,
		User: &user.User{
			UserId:        us.UserId,
			Name:          us.Name,
			FollowCount:   us.FollowCount,
			FollowerCount: us.FollowerCount,
			IsFollow:      relation,
		},
		Content:    delComm.Content,
		CreateTime: utils.TimeToString(delComm.CreatedAt),
	}
	return comm, nil
}

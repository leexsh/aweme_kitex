package service_comment

import (
	"aweme_kitex/cmd/comment/kitex_gen/comment"
	"aweme_kitex/cmd/comment/kitex_gen/user"
	commentRPC "aweme_kitex/cmd/comment/rpc"
	"aweme_kitex/cmd/comment/service_comment/db"
	"aweme_kitex/cmd/feed/kitex_gen/feed"
	user2 "aweme_kitex/cmd/user/kitex_gen/user"
	"aweme_kitex/models"
	"aweme_kitex/pkg/jwt"
	"aweme_kitex/pkg/utils"
	"context"
	"errors"
)

type CreateCommentService struct {
	ctx     context.Context
	uid     string
	vid     string
	content string
}

// NewCreateCommentService new CreateCommentService
func NewCreateCommentService(ctx context.Context) *CreateCommentService {
	return &CreateCommentService{ctx: ctx}
}

// CreateComment add comment
func (s *CreateCommentService) CreateComment(req *comment.CommentActionRequest) (*comment.Comment, error) {
	uc, err := jwt.AnalyzeToken(req.Token)
	if err != nil {
		return nil, err
	}
	s.vid = req.VideoId
	s.uid = uc.Id
	s.content = *req.CommentContent
	return s.do()
}

func (s *CreateCommentService) do() (*comment.Comment, error) {
	// 1. create comment
	commentRaw := &models.CommentRaw{
		Id:      utils.GenerateUUID(),
		UserId:  s.uid,
		VideoId: s.vid,
		Content: s.content,
	}
	err := db.NewCommentDaoInstance().CreateComment(s.ctx, commentRaw)
	// 2.use video rpc add comment count
	err = commentRPC.ChangeCommentCount(s.ctx, &feed.ChangeCommentCountRequest{
		s.vid,
		2,
	})
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	// 3.user rpc get user info
	resp, err := commentRPC.GetUserInfo(s.ctx, &user2.SingleUserInfoRequest{UserIds: []string{s.uid}})
	if err != nil {
		return nil, err
	}
	if len(resp) == 0 {
		return nil, errors.New("not found this user")
	}
	us := resp[0]
	comm := &comment.Comment{
		CommentId: commentRaw.Id,
		User: &user.User{
			UserId:        us.UserId,
			Name:          us.Name,
			FollowCount:   us.FollowCount,
			FollowerCount: us.FollowerCount,
			IsFollow:      false,
		},
		Content:    commentRaw.Content,
		CreateTime: utils.TimeToString(commentRaw.CreatedAt),
	}
	return comm, nil
}

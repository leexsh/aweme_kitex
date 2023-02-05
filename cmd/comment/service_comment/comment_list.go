package service_comment

import (
	"aweme_kitex/cmd/comment/kitex_gen/comment"
	"aweme_kitex/cmd/comment/kitex_gen/user"
	"aweme_kitex/pkg/jwt"
	"aweme_kitex/service"
	"context"
)

type CommentListService struct {
	ctx context.Context
}

// NewCommentListService new CommentListService
func NewCommentListService(ctx context.Context) *CommentListService {
	return &CommentListService{ctx: ctx}
}
func (s *CommentListService) CommentList(req *comment.CommentListRequest) ([]*comment.Comment, error) {
	uc, err := s.CheckToken(req.Token)
	if err != nil {
		return nil, err
	}
	comments, err := service.ShowCommentList(uc.Id, req.VideoId)
	commentList := make([]*comment.Comment, 0)
	for _, comm := range comments {
		curComment := &comment.Comment{
			CommentId: comm.Id,
			User: &user.User{
				UserId:        comm.User.UserId,
				Name:          comm.User.Name,
				FollowCount:   comm.User.FollowCount,
				FollowerCount: comm.User.FollowerCount,
				IsFollow:      comm.User.IsFollow,
			},
			Content:    comm.Content,
			CreateTime: comm.CreateDate,
		}
		commentList = append(commentList, curComment)
	}

	return commentList, nil
}

func (s *CommentListService) CheckToken(token string) (*jwt.UserClaim, error) {
	uc, err := jwt.AnalyzeToken(token)
	if err != nil {
		return nil, err
	}
	return uc, nil
}

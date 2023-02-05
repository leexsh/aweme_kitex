package service_comment

import (
	"aweme_kitex/cmd/comment/kitex_gen/comment"
	"aweme_kitex/cmd/comment/kitex_gen/user"
	"aweme_kitex/pkg/jwt"
	"aweme_kitex/service"
	"context"
)

type DeleteCommentService struct {
	ctx context.Context
}

// NewDeleteCommentService new DeleteCommentService
func NewDeleteCommentService(ctx context.Context) *DeleteCommentService {
	return &DeleteCommentService{ctx: ctx}
}

func (s *DeleteCommentService) DelComment(req *comment.CommentActionRequest) (*comment.Comment, error) {
	uc, err := s.CheckToken(req.Token)
	if err != nil {
		return nil, err
	}
	commet, err := service.DelComment(uc.Id, req.CommentId)
	curComment := &comment.Comment{
		CommentId: commet.Id,
		User: &user.User{
			UserId:        commet.User.UserId,
			Name:          commet.User.Name,
			FollowCount:   commet.User.FollowCount,
			FollowerCount: commet.User.FollowerCount,
			IsFollow:      commet.User.IsFollow,
		},
		Content:    commet.Content,
		CreateTime: commet.CreateDate,
	}
	return curComment, nil

}

func (s *DeleteCommentService) CheckToken(token string) (*jwt.UserClaim, error) {
	uc, err := jwt.AnalyzeToken(token)
	if err != nil {
		return nil, err
	}
	return uc, nil
}

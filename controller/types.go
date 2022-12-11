package controller

import (
	"aweme_kitex/models"
	"aweme_kitex/utils"
	"errors"
)

type (
	Video    models.Video
	User     models.User
	Response models.Response
	Comment  models.Comment
)

// 鉴权
func CheckToken(token string) (*utils.UserClaim, error) {
	if token == defaultToken {
		return nil, errors.New("error: check token failed, please update Token")
	}
	uc, err := utils.AnalyzeToke(token)
	if err != nil {
		return nil, err
	}
	return uc, nil
}

type UserLoginResponse struct {
	Response
	UserId   string `json:"user_id,omitempty"`
	UserName string `json:"user_name,omitempty"`
	Token    string `json:"token,omitempty"`
}

type UserRegisterResponse struct {
	Response
	UserId string `json:"user_id,omitempty"`
	Token  string `json:"token"`
}

type UserResponse struct {
	Response
	User User `json:"user"`
}

type UserListResponse struct {
	Response
	UserList []User `json:"user_list"`
}

type VideoListResponse struct {
	Response
	VideoList []Video `json:"video_list,omitempty"`
}

// comment
type CommentListResponse struct {
	Response
	CommentList []Comment `json:"comment_list,omitempty"`
}

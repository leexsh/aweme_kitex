package models

import (
	"aweme_kitex/cmd/feed/kitex_gen/feed"
	"aweme_kitex/pkg/errno"
	"errors"
	"time"
)

func BuildBaseResp(err error) *feed.BaseResp {
	if err == nil {
		return baseResp(errno.Success)
	}
	e := errno.ErrMsg{}
	if errors.As(err, &e) {
		return baseResp(e)
	}
	s := errno.ServiceErr.WithMessage(err.Error())
	return baseResp(s)
}

func baseResp(err errno.ErrMsg) *feed.BaseResp {
	return &feed.BaseResp{
		StatusCode:  err.ErrCode,
		StatusMsg:   err.ErrMsg,
		ServiceTime: time.Now().Unix(),
	}
}

package pack

import (
	"aweme_kitex/cmd/relation/kitex_gen/relation"
	"aweme_kitex/cmd/relation/kitex_gen/user"
	"time"
)

func RelationActionResponse(code int64, msg string) (resp *relation.RelationActionResponse) {
	resp = new(relation.RelationActionResponse)
	resp.BaseResp.StatusCode = code
	resp.BaseResp.StatusMsg = msg
	resp.BaseResp.ServiceTime = time.Now().Unix()
	return
}

func FollowListResponse(code int64, msg string, users []*user.User) (resp *relation.FollowListResponse) {
	resp = new(relation.FollowListResponse)
	resp.BaseResp.ServiceTime = time.Now().Unix()
	resp.BaseResp.StatusCode = code
	resp.BaseResp.StatusMsg = msg
	if users != nil {
		resp.UserList = users
	}
	return
}

func FollowerListResponse(code int64, msg string, users []*user.User) (resp *relation.FollowerListResponse) {
	resp = new(relation.FollowerListResponse)
	resp.BaseResp.ServiceTime = time.Now().Unix()
	resp.BaseResp.StatusCode = code
	resp.BaseResp.StatusMsg = msg
	if users != nil {
		resp.UserList = users
	}
	return
}

func RelationResponse(code int64, msg string, isFollow bool) (resp *relation.QueryRelationResponse) {
	resp = new(relation.QueryRelationResponse)
	resp.BaseResp.ServiceTime = time.Now().Unix()
	resp.BaseResp.StatusCode = code
	resp.BaseResp.StatusMsg = msg
	resp.IsFollow = isFollow
	return
}

package handlers

import (
	"aweme_kitex/cmd/api/rpc"
	"aweme_kitex/cmd/relation/kitex_gen/relation"
	"aweme_kitex/pkg/errno"
	"aweme_kitex/pkg/jwt"
	"context"

	"github.com/gin-gonic/gin"
)

func RelationAction(c *gin.Context) {
	token := c.Query("token")
	_, err := jwt.AnalyzeToken(token)
	if err != nil {
		SendResponse(c, errno.TokenInvalidErr)
		return
	}
	toUserIdStr := c.Query("to_user_id")
	actionTypeStr := c.Query("action")

	if len(token) == 0 || len(toUserIdStr) == 0 || len(actionTypeStr) == 0 {
		SendResponse(c, errno.ParamErr)
		return
	}

	if actionTypeStr != "1" && actionTypeStr != "2" {
		SendResponse(c, errno.ActionTypeErr)
		return
	}
	req := &relation.RelationActionRequest{
		Token:      token,
		ToUserId:   toUserIdStr,
		ActionType: actionTypeStr,
	}
	err = rpc.RelationAction(context.Background(), req)
	if err != nil {
		SendResponse(c, err)
		return
	}
	SendRelationActionResponse(c, errno.Success)
}

// Followlist get user follow list info
func FollowList(c *gin.Context) {
	token := c.Query("token")
	_, err := jwt.AnalyzeToken(token)
	if err != nil {
		SendResponse(c, errno.TokenInvalidErr)
		return
	}

	req := &relation.FollowListRequest{Token: token}
	userList, err := rpc.FollowList(context.Background(), req)
	if err != nil {
		SendResponse(c, err)
		return
	}
	SendRelationListResponse(c, errno.Success, userList)
}

// FollowerList get user follower list info
func FollowerList(c *gin.Context) {
	token := c.Query("token")
	_, err := jwt.AnalyzeToken(token)
	if err != nil {
		SendResponse(c, errno.TokenInvalidErr)
		return
	}

	req := &relation.FollowerListRequest{Token: token}
	userList, err := rpc.FollowerList(context.Background(), req)
	if err != nil {
		SendResponse(c, err)
		return
	}
	SendRelationListResponse(c, errno.Success, userList)
}

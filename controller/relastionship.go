package controller

import (
	"aweme_kitex/models"
	"aweme_kitex/utils"
	"errors"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

/*
好友关系
*/

func errorResponse(c *gin.Context, code int32, err error) {
	c.JSON(200, Response{
		code,
		err.Error(),
	})
}

// check
func RelationAction(c *gin.Context) {
	user, err := CheckToken(c.Query("token"))
	toUserId := c.Query("to_user_id")
	action, _ := strconv.Atoi(c.Query("action"))
	if err != nil || (action != 1 && action != 2) {
		errorResponse(c, -1, errors.New("params invalid"))
		return
	}
	userId := user.Id
	var relation *RelationRaw = nil
	err = db.Debug().Where("user_id=? AND to_user_id=?", userId, toUserId).Find(&relation).Error
	if err != nil {
		errorResponse(c, -1, err)
		return
	}
	if action == 1 {
		// follow
		var err error
		newRelation := &models.RelationRaw{
			Id:       utils.GenerateUUID(),
			UserId:   user.Id,
			ToUserId: toUserId,
			Status:   1,
		}
		err = models.NewRelationDaoInstance().InsertRaw(newRelation)
		if err != nil {
			errorResponse(c, -1, err)
			return
		}
	} else if action == 2 {
		// unfollow
		var err error
		err = db.Delete(&relation).Error
		if err != nil {
			errorResponse(c, -1, err)
			return
		}
	}
	c.JSON(200, Response{
		200,
		"update relationShip success",
	})
}

// 获取关注
func FollowList(c *gin.Context) {
	var err error
	user, err := CheckToken(c.Query("token"))
	fmt.Println(user)
	relationsUId := make([]string, 0)
	db.Table("relation").Debug().Select("to_user_id").Where("user_id=?", user.Id).Find(&relationsUId)
	usersList, err := accordRelationGetUserInfo(relationsUId)
	if err != nil {
		errorResponse(c, -1, err)
		return
	}
	c.JSON(200, UserListResponse{
		Response{
			StatusCode: 200,
			StatusMsg:  "Success",
		},
		usersList,
	})
}

func accordRelationGetUserInfo(uIds []string) ([]User, error) {
	users := make([]models.UserRawData, 0)
	err := db.Table("user").Debug().Where("user_id IN ?", uIds).Find(&users).Error
	if err != nil {
		return nil, err
	}
	response_user := make([]User, len(users))
	for i, user := range users {
		response_user[i].UserId = user.UserId
		response_user[i].Name = user.Name
		response_user[i].FollowCount = user.FollowCount
		response_user[i].FollowCount = user.FollowerCount
		response_user[i].IsFollow = true
	}
	return response_user, nil
}

// 获取粉丝
func FollowerList(c *gin.Context) {
	var err error
	user, err := CheckToken(c.Query("token"))
	relationsUId := make([]string, 0)
	err = db.Table("relation").Select("user_id").Where("to_user_id=?", user.Id).Find(&relationsUId).Error
	usersList, err := accordRelationGetUserInfo(relationsUId)
	if err != nil {
		errorResponse(c, -1, err)
		return
	}
	c.JSON(200, UserListResponse{
		Response{
			StatusCode: 200,
			StatusMsg:  "Success",
		},
		usersList,
	})
}

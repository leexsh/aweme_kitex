package service

import (
	"aweme_kitex/models"
	"aweme_kitex/models/dal"
	"aweme_kitex/pkg/jwt"
	"aweme_kitex/utils"
	"context"
	"errors"
	"fmt"
)

// user service_user

// register
func RegisterUser(name, password string) (string, string, error) {
	return newRegisterUserDataFlow(name, password).do()
}

type registerUserDataFlow struct {
	userId   string
	userName string
	token    string
	password string
}

func newRegisterUserDataFlow(name, password string) *registerUserDataFlow {
	return &registerUserDataFlow{
		userName: name,
		password: password,
	}
}

func (r *registerUserDataFlow) do() (string, string, error) {
	// insert to data
	userId := utils.GenerateUUID()
	token, _ := jwt.GenerateToken(userId, r.userName)
	newUser := &models.UserRawData{
		UserId:        userId,
		Name:          r.userName,
		Password:      utils.Md5(r.password),
		Token:         token,
		FollowCount:   0,
		FollowerCount: 0,
	}
	// todo 写数据库和缓存
	fmt.Println(newUser)
	err := dal.NewUserDaoInstance().UploadUserData(context.Background(), newUser)
	return userId, token, err
}

// login
func LoginUser(name, password string) (string, string, error) {
	return newLoginUserDataFlow(name, password).do()
}

type loginUserDataFlow struct {
	userId   string
	userName string
	token    string
	password string
}

func newLoginUserDataFlow(name, password string) *loginUserDataFlow {
	return &loginUserDataFlow{
		userName: name,
		password: password,
	}
}
func (l *loginUserDataFlow) do() (uid string, token string, err error) {
	user, err := dal.NewUserDaoInstance().QueryUserByPassword(context.Background(), l.userName, utils.Md5(l.password))
	if user != nil {
		l.token = user.Token
		l.userId = user.UserId
	}
	return l.userId, l.token, err
}

// user info
func QueryUserInfo(user *jwt.UserClaim, remoteUid string) (*models.User, error) {
	return newUserInfoDataFlow(user, remoteUid).do()
}

type userInfoDataFlow struct {
	RemoteUser   *models.UserRawData
	isfollow     bool
	CurrentUName string
	CurrentUId   string
}

func newUserInfoDataFlow(user *jwt.UserClaim, remoteUid string) *userInfoDataFlow {
	return &userInfoDataFlow{
		RemoteUser: &models.UserRawData{
			UserId: remoteUid,
		},
		CurrentUId:   user.Id,
		CurrentUName: user.Name,
	}
}

func (u *userInfoDataFlow) do() (*models.User, error) {
	if err := u.prepareInfo(); err != nil {
		return nil, err
	}
	user, err := u.packUserInfo()
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *userInfoDataFlow) prepareInfo() error {
	// uids := []string{u.RemoteUser.UserId}
	user, err := dal.NewUserDaoInstance().QueryUserByUserId(context.Background(), u.RemoteUser.UserId)
	if err != nil {
		return err
	}
	// if len(users) == 0 {
	// 	return errors.New("user not exist")
	// }
	// u.RemoteUser = users[0]
	u.RemoteUser = user
	relationMap, err := dal.NewRelationDaoInstance().QueryRelationByIds(context.Background(), u.CurrentUId, []string{u.RemoteUser.UserId})
	_, ok := relationMap[u.CurrentUId]
	if !ok {
		u.isfollow = false
	} else {
		u.isfollow = true
	}
	return nil

}

func (u *userInfoDataFlow) packUserInfo() (*models.User, error) {
	if u.RemoteUser == nil {
		return nil, errors.New("NOT FOUND this user")
	}
	user := &models.User{
		UserId:        u.RemoteUser.UserId,
		Name:          u.RemoteUser.Name,
		FollowCount:   u.RemoteUser.FollowCount,
		FollowerCount: u.RemoteUser.FollowerCount,
		IsFollow:      u.isfollow,
	}
	return user, nil
}

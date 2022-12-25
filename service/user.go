package service

import (
	"aweme_kitex/model"
	"aweme_kitex/utils"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

// user service

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
	token, _ := model.GenerateToken(userId, r.userName)
	newUser := &model.UserRawData{
		UserId:        userId,
		Name:          r.userName,
		Password:      utils.Md5(r.password),
		Token:         token,
		FollowCount:   0,
		FollowerCount: 0,
	}
	// todo 写数据库和缓存
	fmt.Println(newUser)
	err := model.NewUserDaoInstance().UploadUserData(newUser)
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
	user, err := model.NewUserDaoInstance().QueryUserByPassword(l.userName, utils.Md5(l.password))
	if user != nil {
		l.token = user.Token
		l.userId = user.UserId
	}
	return l.userId, l.token, err
}

// user info
func QueryUserInfo(user *model.UserClaim, remoteUid string) (*model.User, error) {
	return newUserInfoDataFlow(user, remoteUid).do()
}

type userInfoDataFlow struct {
	RemoteUser   *model.UserRawData
	isfollow     bool
	CurrentUName string
	CurrentUId   string
}

func newUserInfoDataFlow(user *model.UserClaim, remoteUid string) *userInfoDataFlow {
	return &userInfoDataFlow{
		RemoteUser: &model.UserRawData{
			UserId: remoteUid,
		},
		CurrentUId:   user.Id,
		CurrentUName: user.Name,
	}
}

func (u *userInfoDataFlow) do() (*model.User, error) {
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
	uids := []string{u.RemoteUser.UserId}
	userMap, err := model.NewUserDaoInstance().QueryUserByIds(uids)
	if err != nil {
		return err
	}
	u.RemoteUser = userMap[u.RemoteUser.UserId]
	_, err = model.NewRelationDaoInstance().QueryRelationByIds(u.CurrentUId, uids)
	if err == gorm.ErrRecordNotFound {
		u.isfollow = false
	} else if err != nil {
		return err
	} else {
		u.isfollow = true
	}
	return nil

}

func (u *userInfoDataFlow) packUserInfo() (*model.User, error) {
	if u.RemoteUser == nil {
		return nil, errors.New("NOT FOUND this user")
	}
	user := &model.User{
		UserId:        u.RemoteUser.UserId,
		Name:          u.RemoteUser.Name,
		FollowCount:   u.RemoteUser.FollowCount,
		FollowerCount: u.RemoteUser.FollowerCount,
		IsFollow:      u.isfollow,
	}
	return user, nil
}

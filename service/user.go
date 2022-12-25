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
	return NewRegisterUserDataFlow(name, password).do()
}

type RegisterUserDataFlow struct {
	userId   string
	userName string
	token    string
	password string
}

func NewRegisterUserDataFlow(name, password string) *RegisterUserDataFlow {
	return &RegisterUserDataFlow{
		userName: name,
		password: password,
	}
}

func (r *RegisterUserDataFlow) do() (string, string, error) {
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
	return NewLoginUserDataFlow(name, password).do()
}

type LoginUserDataFlow struct {
	userId   string
	userName string
	token    string
	password string
}

func NewLoginUserDataFlow(name, password string) *LoginUserDataFlow {
	return &LoginUserDataFlow{
		userName: name,
		password: password,
	}
}
func (l *LoginUserDataFlow) do() (uid string, token string, err error) {
	user, err := model.NewUserDaoInstance().QueryUserByPassword(l.userName, utils.Md5(l.password))
	if user != nil {
		l.token = user.Token
		l.userId = user.UserId
	}
	return l.userId, l.token, err
}

// user info
func QueryUserInfo(user *model.UserClaim, remoteUid string) (*model.User, error) {
	return NewUserInfoDataFlow(user, remoteUid).do()
}

type UserInfoDataFlow struct {
	RemoteUser   *model.UserRawData
	isfollow     bool
	CurrentUName string
	CurrentUId   string
}

func NewUserInfoDataFlow(user *model.UserClaim, remoteUid string) *UserInfoDataFlow {
	return &UserInfoDataFlow{
		RemoteUser: &model.UserRawData{
			UserId: remoteUid,
		},
		CurrentUId:   user.Id,
		CurrentUName: user.Name,
	}
}

func (u *UserInfoDataFlow) do() (*model.User, error) {
	if err := u.prepareInfo(); err != nil {
		return nil, err
	}
	user, err := u.packUserInfo()
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *UserInfoDataFlow) prepareInfo() error {
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

func (u *UserInfoDataFlow) packUserInfo() (*model.User, error) {
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

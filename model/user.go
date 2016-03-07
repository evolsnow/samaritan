package model

import (
	"github.com/evolsnow/samaritan/base"
	"github.com/evolsnow/samaritan/common/log"
)

type User struct {
	Id         int    `json:"-" redis:"id"` //private id
	Pid        string `json:"id,omitempty" redis:"pid"`
	SamId      string `json:"samId,omitempty" redis:"samId"` //unique id in samaritan
	createTime int64  `json:"createTime,omitempty" redis:"createTime"`
	Alias      string `json:"alias,omitempty" redis:"alias"` //nick name
	Name       string `json:"name,omitempty" redis:"name"`   //real name
	Phone      string `json:"phone,omitempty" redis:"phone"`
	Password   string `json:"passwd" redis:"passwd"`
	Email      string `json:"email,omitempty" redis:"email"`
	Avatar     string `json:"avatar,omitempty" redis:"avatar"` //avatar url
	School     string `json:"school,omitempty" redis:"school"`
	Department string `json:"depart,omitempty" redis:"depart"`
	Grade      int    `json:"grade,omitempty" redis:"grade"`
	Class      string `json:"class,omitempty" redis:"class"`
	StudentNum string `json:"stuNum,omitempty" redis:"stuNum"` //1218404001...
}

//read user's password
func (u *User) GetPassword() (pwd string) {
	pwd, err := readPassword(u.Id)
	if err != nil {
		log.Error("Error get user's password:", err)
		return ""
	}
	return
}

//generate avatar url for user
func (u *User) CreateAvatar() {
	path, err := base.GenerateAvatar(u.Phone)
	if err != nil {
		log.Error("Error generate user's avatar:", err)
	}
	u.Avatar = path
	err = createUserAvatar(u.Id, u.Avatar)
	if err != nil {
		log.Error("Error create user's avatar:", err)
	}
}

//full url of avatar img
func (u *User) FullAvatarUrl() string {
	prefix := "https://samaritan.tech:10000"
	log.Debug(prefix + u.Avatar)
	return prefix + u.Avatar
}

//save a new user
func (u *User) Save() {
	//return user id for jwt token use
	log.Debug("create user:", u)
	createUser(u)
}

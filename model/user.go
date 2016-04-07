package model

import (
	"github.com/evolsnow/samaritan/common/base"
	"github.com/evolsnow/samaritan/common/log"
)

type User struct {
	Id         int    `json:"-" redis:"id"` //private id
	Pid        string `json:"id,omitempty" redis:"pid"`
	SamId      string `json:"samId,omitempty" redis:"samId"` //unique id in samaritan
	createTime int64  `redis:"createTime"`
	Alias      string `json:"alias,omitempty" redis:"alias"` //nick name
	Name       string `json:"name,omitempty" redis:"name"`   //real name
	Phone      string `json:"phone,omitempty" redis:"phone"`
	Password   string `json:"password" redis:"password"`
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

//user created projects
func (u *User) GetCreatedProjects() []Project {
	ret, err := readCreatedProjects(u.Id)
	if err != nil {
		log.Error("Err get created projects:", err)
		return nil
	}
	return ret
}

//user joined projects
func (u *User) GetJoinedProjects() []Project {
	ret, err := readJoinedProjects(u.Id)
	if err != nil {
		log.Error("Err get joined projects:", err)
		return nil
	}
	return ret
}

//all projects
func (u *User) GetAllProjects() []Project {
	pjs := u.GetCreatedProjects()
	tmp := u.GetJoinedProjects()
	in := func(a Project, list []Project) bool {
		for _, b := range list {
			if b.Id == a.Id {
				return true
			}
		}
		return false
	}
	for _, p := range tmp {
		if !in(p, pjs) {
			pjs = append(pjs, p)
		}
	}
	return pjs
}

//generate avatar url for user
func (u *User) CreateAvatar() {
	path, err := base.GenerateAvatar(u.Email)
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
	prefix := "https://img.samaritan.tech"
	log.Debug(prefix + u.Avatar)
	return prefix + u.Avatar
}

//create or update a new user
func (u *User) Save() {
	if u.Id == 0 {
		//new user
		log.DebugJson("create user:", u)
		createUser(u)
	} else {
		kvMap := prepareToUpdate(u)
		log.DebugJson("update user with: ", kvMap)
		updateUser(u.Id, kvMap)
	}
}

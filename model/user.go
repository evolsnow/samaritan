package model

import (
	"github.com/evolsnow/samaritan/common/base"
	"github.com/evolsnow/samaritan/common/log"
	"strings"
)

type User struct {
	Id         int    `json:"-" redis:"id"` //private id
	Pid        string `json:"id,omitempty" redis:"pid"`
	SamId      string `json:"samId,omitempty" redis:"samId"`           //unique id in samaritan
	CreateTime int64  `json:"createTime,omitempty" redis:"createTime"` //create time timestamp of this user
	Alias      string `json:"alias,omitempty" redis:"alias"`           //nick name
	Name       string `json:"name,omitempty" redis:"name"`             //real name
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

// InitedUser returns a full loaded user object by id
func InitedUser(id int) (u *User) {
	u = &User{Id: id}
	u.load()
	if u.Pid == "" {
		return nil
	}
	u.Avatar = u.FullAvatarUrl()
	return
}

// Sync reloads user from db
func (u *User) Sync() {
	*u = *InitedUser(u.Id)
}

//// GetPassword reads user's password
//func (u *User) GetPassword() (pwd string) {
//	pwd, err := readPassword(u.Id)
//	if err != nil {
//		log.Error("Error get user's password:", err)
//		return ""
//	}
//	return
//}

// GetName reads user's real name
func (u *User) GetName() (name string) {
	name, err := readName(u.Id)
	if err != nil {
		log.Error("Error get user's name:", err)
		return ""
	}
	return
}

// GetCreatedProjects gets user created projects
func (u *User) GetCreatedProjects() []Project {
	ret, err := readCreatedProjects(u.Id)
	if err != nil {
		log.Error("Error get created projects:", err)
		return nil
	}
	return ret
}

// GetJoinedProjects gets user joined projects
func (u *User) GetJoinedProjects() []Project {
	ret, err := readJoinedProjects(u.Id)
	if err != nil {
		log.Error("Error get joined projects:", err)
		return nil
	}
	return ret
}

// GetAllProjects gets all projects
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

// GetAllCompletedMissionsId gets all completed missions id
func (u *User) GetAllCompletedMissionsId() []int {
	var ids []int
	ids, err := readCompletedMissionsId(u.Id)
	if err != nil {
		log.Error("Error get completed missions:", err)
	}
	return ids
}

// GetAllAcceptedMissionsId gets all accepted missions id
func (u *User) GetAllAcceptedMissionsId() []int {
	var ids []int
	ids, err := readAcceptedMissionsId(u.Id)
	if err != nil {
		log.Error("Error get accepted missions:", err)
	}
	return ids
}

func (u *User) GetAllOfflineMsg() []*Chat {
	chs, err := readOfflineMsg(u.Id)
	if err != nil {
		log.Error("Error get user offline msgs:", err)
	}
	return chs
}

// CompleteMission adds mission to user completed mission
func (u *User) CompleteMission(mid int) {
	if err := updateCompletedMission(u.Id, mid); err != nil {
		log.Error("Error update completed mission:", err)
	}
}

// UnCompleteMission adds mission to user uncompleted mission
func (u *User) UnCompleteMission(mid int) {
	if err := updateUnCompletedMission(u.Id, mid); err != nil {
		log.Error("Error update uncompleted mission:", err)
	}
}

// AcceptMission add mission to user accepted mission set
func (u *User) AcceptMission(mid int) {
	updateAcceptedMission(u.Id, mid)
}

// JoinProject add project to user joined project set
func (u *User) JoinProject(pid int) {
	updateJoinedProject(u.Id, pid)
}

// CreateAvatar generates avatar for user
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

// FullAvatarUrl return full url of avatar img
func (u *User) FullAvatarUrl() string {
	prefix := "https://img.samaritan.tech/"
	if strings.Contains(u.Avatar, "/avatar") {
		//default or qi niu
		return prefix + u.Avatar
	}
	return base.QiNiuDownloadUrl(u.Avatar)
}

// Save creates or updates a new user
func (u *User) Save() {
	if u.Id == 0 {
		//new user
		log.DebugJson("create user:", u)
		createUser(u)
		//go CreateUserMysql(*u)

	} else {
		kvMap := prepareToUpdate(u)
		log.DebugJson("update user with: ", kvMap)
		updateUser(u.Id, kvMap)
	}
}

// Load full read from redis
func (u *User) load() (err error) {
	uPtr, err := readUserWithId(u.Id)
	if err != nil {
		log.Debug(err)
	}
	*u = *uPtr
	return
}

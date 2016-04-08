package model

import (
	"github.com/evolsnow/samaritan/common/dbms"
	"github.com/evolsnow/samaritan/common/log"
)

type Project struct {
	Id         int    `json:"-" redis:"id"`             //private id
	Pid        string `json:"id,omitempty" redis:"pid"` //public id
	Name       string `json:"name,omitempty" redis:"name"`
	createTime int64  `redis:"createTime"`                           //create time timestamp of this project
	Desc       string `json:"desc,omitempty" redis:"desc"`           //description for the project
	CreatorId  int    `json:"creatorId,omitempty" redis:"creatorId"` //who created the project
	Private    bool   `json:"private,omitempty" redis:"private"`
	MembersId  []int  `json:"membersId,omitempty" redis:"-"` //user list who in this project
}

func (p *Project) GetCreator() (creator *User) {
	creator, err := readCreator(p.Id)
	if err != nil {
		log.Error("Error get creator:", err)
		return nil
	}
	log.DebugJson("get creator:", creator)
	return
}

func (p *Project) AddMember(uid int) (err error) {
	err = updateProjectMember(p.Id, uid, 1)
	if err != nil {
		log.Error("Error add Member:", err)
		return err
	}
	return
}

func (p *Project) RemoveMember(uid int) (err error) {
	err = updateProjectMember(p.Id, uid, -1)
	if err != nil {
		log.Error("Error remove Member:", err)
		return err
	}
	return
}

func (p *Project) GetMembers() (members []*User) {
	members, err := readProjectMembers(p.Id)
	if err != nil {
		log.Error("Error get project members:", err)
		return nil
	}
	log.Debug("proj members:", members)
	return
}

func (p *Project) GetMembersId() []int {
	if p.Id == 0 {
		p.Id = dbms.ReadProjectId(p.Pid)
	}
	ids, err := readProjectMembersId(p.Id)
	if err != nil {
		log.Error("Error get project members", err)
		return nil
	}
	log.Debug("proj members id:", ids)
	return ids
}

func (p *Project) GetMissions() (missions []*Mission) {
	missions, err := readProjectMissions(p.Id)
	if err != nil {
		log.Error("Error get project missions:", err)
		return nil
	}
	log.Debug("proj missions:", missions)
	return
}

func (p *Project) Save() {
	if p.Id == 0 {
		//new project
		log.DebugJson("create project:", p)
		createProject(p)
	} else {
		kvMap := prepareToUpdate(p)
		log.Debug("update project with: ", kvMap)
		updateProject(p.Id, kvMap)
	}
}

//delete a project
func (p *Project) Remove() (err error) {
	if err = deleteProject(p.Id); err != nil {
		log.Error("Error delete project:", err)
	}
	if p.CreatorId == 0 {
		p.CreatorId = p.GetCreator().Id
	}
	if err = updateUserProjectSet(p.Id, p.CreatorId); err != nil {
		log.Error("Err update pj set")
	}
	return
}

//full read from redis
func (p *Project) Load() (err error) {
	err = readFullProject(p)
	if err != nil {
		log.Debug(err)
	}
	return
}

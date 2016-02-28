package model

import (
	"log"
)

type Project struct {
	Id         int    `json:"-" redis:"id"`             //private id
	Pid        string `json:"id,omitempty" redis:"pid"` //public id
	Name       string `json:"name,omitempty" redis:"name"`
	createTime int64  `json:"createTime,omitempty" redis:"createTime"`   //create time timestamp of this project
	Desc       string `json:"desc,omitempty" redis:"desc"`               //description for the project
	CreatorId  int    `json:"publisherId,omitempty" redis:"publisherId"` //who created the project
	Private    bool   `json:"private,omitempty" redis:"private"`
	MembersId  []int  `json:"membersId,omitempty" redis:"-"` //user list who in this project
}

func (p *Project) GetCreator() (creator *User) {
	creator, err := readUser(p.CreatorId)
	if err != nil {
		log.Println("Error get creator:", err)
		return nil
	}
	return
}

func (p *Project) AddMember(uid int) (err error) {
	err = updateProjectMember(p.Id, uid, 1)
	if err != nil {
		log.Println("Error add Member:", err)
		return err
	}
	return
}

func (p *Project) RemoveMember(uid int) (err error) {
	err = updateProjectMember(p.Id, uid, -1)
	if err != nil {
		log.Println("Error remove Member:", err)
		return err
	}
	return
}

func (p *Project) GetMembers() (participants []*User) {
	if len(p.MembersId) == 0 {
		return
	}
	participants, err := readProjectMembers(p.Id)
	if err != nil {
		log.Println("Error get project members:", err)
		return nil
	}
	return
}

func (p *Project) GetMembersId() []int {
	if p.Id == 0 {
		p.Id = ReadProjectId(p.Pid)
	}
	ids, err := readProjectMembersId(p.Id)
	if err != nil {
		log.Println("Error get project members with name", err)
		return nil
	}
	return ids
}

func (p *Project) Save() {
	createProject(p)
}

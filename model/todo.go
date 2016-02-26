package model

import (
	"log"
)

type Todo struct {
	Id           int    `json:"id, omitempty" redis:"id"`
	StartTime    uint64 `json:"startTime, omitempty" redis:"startTime"` //start timestamp of this action
	Deadline     uint64 `json:"deadline, omitempty" redis:"dealine"`    //end time
	Desc         string `json:"desc, omitempty" redis:"desc"`           //description for the action
	OwnerId      int    `json:"ownerId, omitempty" redis:"ownerId"`     //whose
	Accomplished bool   `json:"accomplished, omitempty" redis:"accomplished"`
	MissionId    int    `json:"missionId, omitempty" redis:"missionId"` //belong to which mission
}

//get user from to-do's owner id
func (td *Todo) GetOwner() (owner *User) {
	owner, err := readUser(td.OwnerId)
	if err != nil {
		log.Println("Error get user with todo:", err)
		return nil
	}
	return
}

//get mission from to-do's mission id
func (td *Todo) GetMission() (m *Mission) {
	m, err := readMission(td.MissionId)
	if err != nil {
		log.Println("Error get mission with todo:", err)
		return nil
	}
	return
}

//update to-do accomplish status
func (td *Todo) Done() (err error) {
	err = updateTodoStatus(td.OwnerId, td.Id)
	if err != nil {
		log.Println("Error update to-do status:", err)
		return
	}
	return
}

//save a new to-do
func (td *Todo) Save() (err error) {
	err = createTodo(td)
	if err != nil {
		log.Println("Error save to-do:", err)
		return
	}
	return
}

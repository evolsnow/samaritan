package model

import (
	"github.com/garyburd/redigo/redis"
	"log"
)

type Todo struct {
	Id           int    `json:"id" redis:"id"`
	StartTime    uint64 `json:"startTime" redis:"startTime"` //start timestamp of this action
	DeadLine     uint64 `json:"deadLine" redis:"deadLine"`   //end time
	Desc         string `json:"desc" redis:"desc"`           //description for the action
	OwnerId      int    `json:"ownerId" redis:"ownerId"`     //whose
	Accomplished bool   `json:"accomplished" redis:"accomplished"`
	MissionId    int    `json:"missionId" redis:"missionId"` //belong to which mission
}

//get user from to-do's owner id
func (td *Todo) GetOwner() (owner *User) {
	reply, err := readUser(td.OwnerId)
	if err != nil {
		log.Println("Error get user with todo:", err)
		return
	}
	redis.ScanStruct(reply, owner)
	return
}

//get mission from to-do's mission id
func (td *Todo) GetMission() (m *Mission) {
	reply, err := readMission(td.MissionId)
	if err != nil {
		log.Println("Error get mission with todo:", err)
		return
	}
	redis.ScanStruct(reply, m)
	return
}

//update to-do accomplish status
func (td *Todo) Accomplished() (err error) {
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

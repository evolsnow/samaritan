package model

import (
	"github.com/garyburd/redigo/redis"
)

type Todo struct {
	Id        int    `json:"id" redis:"id"`
	StartTime uint64 `json:"startTime" redis:"startTime"` //start timestamp of this action
	DeadLine  uint64 `json:"deadLine" redis:"deadLine"`   //end time
	Desc      string `json:"desc" redis:"desc"`           //description for the action
	OwnerId   int    `json:"ownerId" redis:"ownerId"`     //whose
	Status    int    `json:"status" redis:"status"`       //0:not begin, 1:ongoing, 2: overdue, 3:accomplished
	MissionId int    `json:"missionId" redis:"missionId"` //belong to which mission
}

func (td *Todo) GetOwner() (owner *User) {
	reply, err := readUser(td.OwnerId)
	if err != nil {
		return
	}
	redis.ScanStruct(reply, owner)
	return
}

func (td *Todo) GetMission() (m *Mission) {
	reply, err := readMission(td.MissionId)
	if err != nil {
		return
	}
	redis.ScanStruct(reply, m)
	return
}

package model

import (
	"log"
)

type Todo struct {
	Id           int    `json:"-" redis:"id"` //private id
	Pid          string `json:"id,omitempty" redis:"pid"`
	StartTime    int64  `json:"startTime,omitempty" redis:"startTime"` //start timestamp of this action
	TaskTime     int64  `json:"taskTime,omitempty" redis:"taskTime"`
	Place        string `json:"place,omitempty" redis:"place"`
	Repeat       bool   `json:"repeat,omitempty" redis:"repeat"`
	RepeatPeriod int64  `json:"repeatPeriod,omitempty" redis:"repeatPeriod"`
	//Deadline     int64  `json:"deadline,omitempty" redis:"dealine"`    //end time
	Desc       string `json:"desc,omitempty" redis:"desc"` //description for the action
	Remark     string `json:"remark,omitempty" redis:"remark"`
	OwnerId    int    `json:"ownerId,omitempty" redis:"ownerId"` //whose
	Done       bool   `json:"done,omitempty" redis:"done"`
	FinishTime int64  `json:"finishTime,omitempty" redis:"finishTime"`
	MissionId  int    `json:"missionId,omitempty" redis:"missionId"` //belong to which mission
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

//get project from to-do's project id
func (td *Todo) GetMission() (m *Mission) {
	m, err := readMission(td.MissionId)
	if err != nil {
		log.Println("Error get mission with todo:", err)
		return nil
	}
	return
}

//update to-do done status
func (td *Todo) Finish() (err error) {
	err = updateTodoStatus(td.OwnerId, td.Id)
	if err != nil {
		log.Println("Error update to-do status:", err)
		return
	}
	return
}

//save a new to-do
func (td *Todo) Save() {
	createTodo(td)
}

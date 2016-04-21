package model

import (
	"github.com/evolsnow/samaritan/common/base"
	"github.com/evolsnow/samaritan/common/log"
)

type Todo struct {
	Id         int      `json:"-" redis:"id"` //private id
	Pid        string   `json:"id,omitempty" redis:"pid"`
	CreateTime int64    `json:"createTime,omitempty" redis:"createTime"` //create time timestamp of this todo
	StartTime  int64    `json:"startTime,omitempty" redis:"startTime"`   //start timestamp of this action
	Place      string   `json:"place,omitempty" redis:"place"`
	Pictures   []string `json:"pictures,omitempty" redis:"-"`
	Repeat     bool     `json:"repeat,omitempty" redis:"repeat"`
	RepeatMode int      `json:"repeatMode,omitempty" redis:"repeatMode"`
	AllDay     bool     `json:"allDay,omitempty" redis:"allDay"`
	Desc       string   `json:"desc,omitempty" redis:"desc"` //description for the action
	Remark     string   `json:"remark,omitempty" redis:"remark"`
	OwnerId    int      `json:"ownerId,omitempty" redis:"ownerId"` //whose
	Done       bool     `json:"done,omitempty" redis:"done"`
	FinishTime int64    `json:"finishTime,omitempty" redis:"finishTime"`
	MissionId  int      `json:"missionId,omitempty" redis:"missionId"` //belong to which mission
}

// InitedTodo returns a full loaded to-do object by id
func InitedTodo(id int) (td *Todo) {
	td = &Todo{Id: id}
	td.load()
	td.Pictures = td.GetPics()
	return
}

// Sync reloads to-do from db
func (td *Todo) Sync() {
	*td = *InitedTodo(td.Id)
}

// GetOwner gets user from to-do's owner id
func (td *Todo) GetOwner() (owner *User) {
	var err error
	if td.OwnerId == 0 {
		owner, err = readOwner(td.Id)
	} else {
		owner, err = readUserWithId(td.OwnerId)
	}
	if err != nil {
		log.Error("Error get user with todo:", err)
		return nil
	}
	log.DebugJson("get owner:", owner)
	return
}

// GetMission gets project from to-do's project id
func (td *Todo) GetMission() (m *Mission) {
	var err error
	if td.MissionId == 0 {
		m, err = readBelongedMission(td.Id)
	} else {
		m, err = readMissionWithId(td.MissionId)
	}
	if err != nil {
		log.Error("Error get mission with todo:", err)
		return nil
	}
	log.Debug("get mission:", m)
	return
}

// Finish updates to-do done status
func (td *Todo) Finish() (err error) {
	if td.OwnerId == 0 {
		td.OwnerId = td.GetOwner().Id
	}
	err = updateTodoStatus(td.OwnerId, td.Id)
	if err != nil {
		log.Error("Error update to-do status:", err)
	}
	return
}

//GetPics get to-do's all pics
func (td *Todo) GetPics() (pics []string) {
	raw, err := readTodoPics(td.Id)
	if err != nil {
		log.Error("Error get to-do's pics")
		return
	}
	pics = make([]string, len(raw))
	for i, v := range raw {
		pics[i] = base.QiNiuDownloadUrl(v)
	}
	return
}

// UpdatePics update to-do's pictures list
func (td *Todo) UpdatePics(pics []string) (err error) {
	err = updateTodoPics(td.Id, pics)
	if err != nil {
		log.Error("Error update to-do pics:", err)
	}
	return
}

// ave a new to-do
func (td *Todo) Save() {
	if td.Id == 0 {
		//new to-do
		log.DebugJson("create todo:", td)
		createTodo(td)
		//go CreateTodoMysql(*td)

	} else {
		kvMap := prepareToUpdate(td)
		log.Debug("update todo with: ", kvMap)
		updateTodo(td.Id, kvMap)
	}
}

// Remove deletes a to-do
func (td *Todo) Remove() (err error) {
	if err = deleteTodo(td.Id); err != nil {
		log.Error("Error delete todo:", err)
	}
	//go DeleteTodoMysql(td.Id)
	return
}

// Load full read from redis
func (td *Todo) load() (err error) {
	tPtr, err := readTodoWithId(td.Id)
	if err != nil {
		log.Debug(err)
	}
	*td = *tPtr
	return
}

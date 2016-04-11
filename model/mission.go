package model

import (
	"github.com/evolsnow/samaritan/common/dbms"
	"github.com/evolsnow/samaritan/common/log"
)

type Mission struct {
	Id            int       `json:"-" redis:"id"`             //private id
	Pid           string    `json:"id,omitempty" redis:"pid"` //public id
	createTime    int64     `redis:"createTime"`              //create time timestamp of this mission
	Name          string    `json:"name,omitempty" redis:"name"`
	Desc          string    `json:"desc,omitempty" redis:"desc"`                   //description for the project
	PublisherId   int       `json:"publisherId,omitempty" redis:"publisherId"`     //who published the mission
	ReceiversId   []int     `json:"receiversId,omitempty" redis:"-"`               //user list who accepted the mission
	CompletionNum int       `json:"completionNum,omitempty" redis:"completionNum"` //completed number
	CompletedTime int64     `json:"completedTime,omitempty" redis:"completedTime"`
	Comments      []Comment `json:"comments,omitempty" redis:"-"`
	ProjectId     int       `json:"projectId,omitempty" redis:"projectId"` //belong to which mission
}

type Comment struct {
	Id         int    `json:"-" redis:"id"`
	Pid        string `json:"id,omitempty" redis:"pid"`
	When       int64  `json:"when,omitempty" redis:"when"`
	MissionPid string `json:"-" redis:"-"`
	CriticPid  string `json:"uid,omitempty" redis:"criticPid"`
	CriticName string `json:"uName,omitempty" redis:"criticName"`
}

func (m *Mission) GetReceiversId() []int {
	if m.Id == 0 {
		m.Id = dbms.ReadMissionId(m.Pid)
	}
	ids, err := readMissionReceiversId(m.Id)
	if err != nil {
		log.Error("Error get mission receivers", err)
		return nil
	}
	log.Debug("receivers id:", ids)
	return ids
}

func (m *Mission) AddReceiver(uid int) (err error) {
	err = updateMissionReceiver(m.Id, uid, 1)
	if err != nil {
		log.Error("Error add receiver:", err)
		return err
	}
	return
}

func (m *Mission) GetComments() (comments []*Comment) {
	comments, err := readMissionComments(m.Id)
	if err != nil {
		log.Error("Error get mission comments:", err)
		return nil
	}
	log.Debug("mission comments:", comments)
	return
}

//save a new mission
func (m *Mission) Save() {
	if m.Id == 0 {
		//new mission
		log.DebugJson("create mission:", m)
		createMission(m)
	} else {
		kvMap := prepareToUpdate(m)
		log.Debug("update mission with: ", kvMap)
		updateMission(m.Id, kvMap)
	}
}

//save a new mission
func (m *Mission) ForceSave() {
	if m.Id == 0 {
		//new mission
		log.DebugJson("force create mission:", m)
		createMission(m)
	} else {
		kvMap := prepareToForceUpdate(m)
		log.Debug("force update mission with: ", kvMap)
		updateMission(m.Id, kvMap)
	}
}

func (cm *Comment) Save() {
	log.Debug("create comment:", cm)
	createMissionComment(cm)
}

//full read from redis
func (m *Mission) Load() (err error) {
	err = readFullMission(m)
	if err != nil {
		log.Debug(err)
	}
	return
}

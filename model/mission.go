package model

import "github.com/evolsnow/samaritan/common/log"

type Mission struct {
	Id            int       `json:"-" redis:"id"`                            //private id
	Pid           string    `json:"id,omitempty" redis:"pid"`                //public id
	createTime    int64     `json:"createTime,omitempty" redis:"createTime"` //create time timestamp of this mission
	Name          string    `json:"name,omitempty" redis:"name"`
	Desc          string    `json:"desc,omitempty" redis:"desc"`                   //description for the project
	PublisherId   int       `json:"publisherId,omitempty" redis:"publisherId"`     //who published the mission
	ReceiversId   []int     `json:"receiversId,omitempty" redis:"-"`               //user list who accepted the mission
	CompletionNum int       `json:"completionNum,omitempty" redis:"completionNum"` //completed number
	CompletedTime int64     `json:"completedTime,omitempty" redis:"completedTime"`
	Comments      []Comment `json:"comments,omitempty" redis:"-"`
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
		m.Id = ReadMissionId(m.Pid)
	}
	ids, err := readMissionReceiversId(m.Id)
	if err != nil {
		log.Error("Error get mission receivers", err)
		return nil
	}
	log.Debug("receivers id:", ids)
	return ids
}

func (cm *Comment) Save() {
	log.Debug("create comment:", cm)
	createMissionComment(cm)
}

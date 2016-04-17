package model

import (
	"github.com/evolsnow/samaritan/common/base"
	"github.com/evolsnow/samaritan/common/dbms"
	"github.com/evolsnow/samaritan/common/log"
)

type Mission struct {
	Id            int       `json:"-" redis:"id"`                            //private id
	Pid           string    `json:"id,omitempty" redis:"pid"`                //public id
	CreateTime    int64     `json:"createTime,omitempty" redis:"createTime"` //create time timestamp of this mission
	Name          string    `json:"name,omitempty" redis:"name"`
	Desc          string    `json:"desc,omitempty" redis:"desc"` //description for the mission
	Pictures      []string  `json:"pictures,omitempty" redis:"-"`
	PublisherId   int       `json:"publisherId,omitempty" redis:"publisherId"`     //who published the mission
	ReceiversId   []int     `json:"receiversId,omitempty" redis:"-"`               //user list who accepted the mission
	CompletionNum int       `json:"completionNum,omitempty" redis:"completionNum"` //completed number
	CompletedTime int64     `json:"completedTime,omitempty" redis:"completedTime"`
	Deadline      int64     `json:"deadline,omitempty" redis:"deadline"`
	Comments      []Comment `json:"comments,omitempty" redis:"-"`
	ProjectId     int       `json:"projectId,omitempty" redis:"projectId"` //belong to which mission
}

type Comment struct {
	Id         int    `json:"-" redis:"id"`
	Pid        string `json:"id,omitempty" redis:"pid"`
	CreateTime int64  `json:"createTime,omitempty" redis:"createTime"`
	MissionPid string `json:"-" redis:"-"`
	CriticPid  string `json:"uid,omitempty" redis:"criticPid"`
	CriticName string `json:"uName,omitempty" redis:"criticName"`
}

// GetReceiversId gets mission's receivers' id slice
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

// AddReceiver adds a user to mission's receiver set
func (m *Mission) AddReceiver(uid int) (err error) {
	err = updateMissionReceiver(m.Id, uid, 1)
	if err != nil {
		log.Error("Error add receiver:", err)
		return err
	}
	m.UpdateCompleteNum()
	return
}

// GetComments gets mission's comments
func (m *Mission) GetComments() (comments []*Comment) {
	comments, err := readMissionComments(m.Id)
	if err != nil {
		log.Error("Error get mission comments:", err)
		return nil
	}
	log.Debug("mission comments:", comments)
	return
}

// GetCompletedUsersId gets users id slice who have completed the mission
func (m *Mission) GetCompletedUsersId() []int {
	if m.Id == 0 {
		m.Id = dbms.ReadMissionId(m.Pid)
	}
	ids, err := readMissionCompletedUsersId(m.Id)
	if err != nil {
		log.Error("Error get mission completed user", err)
		return nil
	}
	log.Debug("users id:", ids)
	return ids
}

// UpdateCompleteNum when user-completed num changed
func (m *Mission) UpdateCompleteNum() {
	lg := len(m.GetReceiversId())
	if lg == 0 {
		return
	}
	m.CompletionNum = 100 * len(m.GetCompletedUsersId()) / len(m.GetReceiversId())
}

// GetPictures return mission's pics
func (m *Mission) GetPictures() (pics []string) {
	raw, err := readMissionPics(m.Id)
	if err != nil {
		log.Error("Error get mission's pics")
		return
	}
	pics = make([]string, len(raw))
	for i, v := range pics {
		pics[i] = base.QiNiuDownloadUrl(v)
	}
	return
}

// UpdatePics update to-do's pictures list
func (m *Mission) UpdatePics(pics []string) (err error) {
	err = updateMissionPics(m.Id, pics)
	if err != nil {
		log.Error("Error update mission pics:", err)
	}
	return
}

// Save a new mission
func (m *Mission) Save() {
	if m.Id == 0 {
		//new mission
		log.DebugJson("create mission:", m)
		createMission(m)
		//go CreateMissionMysql(*m)

	} else {
		kvMap := prepareToUpdate(m)
		log.Debug("update mission with: ", kvMap)
		updateMission(m.Id, kvMap)
	}
}

// ForceSave a new mission
func (m *Mission) ForceSave() {
	if m.Id == 0 {
		//new mission
		log.DebugJson("force create mission:", m)
		createMission(m)
		//go CreateMissionMysql(*m)

	} else {
		kvMap := prepareToForceUpdate(m)
		log.Debug("force update mission with: ", kvMap)
		updateMission(m.Id, kvMap)
	}
}

// Save a comment
func (cm *Comment) Save() {
	log.Debug("create comment:", cm)
	createMissionComment(cm)
}

// Load full read from redis
func (m *Mission) Load() (err error) {
	err = readFullMission(m)
	if err != nil {
		log.Debug(err)
	}
	return
}

package model

import (
	"log"
)

type Mission struct {
	Id          int    `json:"id, omitempty" redis:"id"`
	StartTime   uint64 `json:"startTime, omitempty" redis:"startTime"`     //start timestamp of this action
	Desc        string `json:"desc, omitempty" redis:"desc"`               //description for the action
	PublisherId int    `json:"publisherId, omitempty" redis:"publisherId"` //who published the mission
	ReceiversId []int  `json:"receiversId, omitempty" redis:"-"`           //user list who received the mission
}

func (m *Mission) GetPublisher() (publisher *User) {
	publisher, err := readUser(m.PublisherId)
	if err != nil {
		log.Println("Error get publisher:", err)
		return nil
	}
	return
}

func (m *Mission) AddReceiver(userId int) (err error) {
	err = createMissionRcv(m.Id, userId)
	if err != nil {
		log.Println("Error add receiver:", err)
		return err
	}
	return
}

func (m *Mission) GetReceivers() (receivers []*User) {
	if len(m.ReceiversId) == 0 {
		return
	}
	receivers, err := readMissionRcv(m.Id)
	if err != nil {
		log.Println("Error get mission receivers:", err)
		return nil
	}
	return
}

func (m *Mission) Save() (err error) {
	err = createMission(m)
	if err != nil {
		log.Println("Error save mission:", err)
		return err
	}
	return
}

package model

import (
	"github.com/garyburd/redigo/redis"
	"log"
)

type Mission struct {
	Id          int    `json:"id" redis:"id"`
	StartTime   uint64 `json:"startTime" redis:"startTime"`     //start timestamp of this action
	DeadLine    uint64 `json:"deadLine" redis:"deadLine"`       //end time
	Desc        string `json:"desc" redis:"desc"`               //description for the action
	Color       [3]int `json:"color" redis:"-"`                 //RGB mode
	PublisherId int    `json:"publisherId" redis:"publisherId"` //who published the mission
	ReceiversId []int  `json:"receiversId" redis:"-"`           //user list who received the mission
}

func (m *Mission) GetPublisher() (publisher *User) {
	publisher, err := readUser(m.PublisherId)
	if err != nil {
		log.Println("Error get publisher:", err)
		return nil, nil
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

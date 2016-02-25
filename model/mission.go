package model

import (
	"github.com/garyburd/redigo/redis"
)

type Mission struct {
	Id          int
	StartTime   uint64 //start timestamp of this action
	DeadLine    uint64 //end time
	Desc        string //description for the action
	Color       [3]int //RGB mode
	PublisherId int    //who published the mission
	ReceiversId []int  //user list who received the mission
}

func (m *Mission) GetPublisher() (owner *User) {
	reply, err := readUser(m.PublisherId)
	if err != nil {
		return
	}
	redis.ScanStruct(reply, owner)
	return
}

func (m *Mission) AddReceiver(userId int) (err error) {
	return updateMissionRcv(m.Id, userId)
}

func (m *Mission) GetReceivers() (receivers []*User) {
	if len(m.ReceiversId) == 0 {
		return
	}
	replys, err := readMissionRcv(m.Id, m.ReceiversId)
	if err != nil {
		return
	}
	for i := range replys {
		user := new(User)
		redis.ScanStruct(replys[i], user)
		receivers = append(receivers, user)

	}
	return
}

func (m *Mission) Save() (err error) {
	return createMission(m)
}

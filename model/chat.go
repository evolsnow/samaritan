package model

import (
	"encoding/json"
	"github.com/evolsnow/samaritan/common/base"
	"github.com/evolsnow/samaritan/common/dbms"
	"github.com/evolsnow/samaritan/common/log"
	"github.com/evolsnow/samaritan/common/rpc"
	"time"
)

const (
	PeerToPeer = iota //private msg server <<-->> client
	Discuss           //group chat server <<-->> client
	//system call
	UserJoined        //server -->> client
	UserLeft          //server -->> client
	InvitedToProject  //server -->> client
	KickedFromProject //server -->> client
	InvitedToMission  //server -->> client
)

type Chat struct {
	Id             int               `json:"-" redis:"id"`
	ConversationId string            `json:"convId,omitempty" redis:"convId"`
	Type           int               `json:"type" redis:"type"`
	Msg            string            `json:"msg,omitempty" redis:"msg"`
	Target         string            `json:"target,omitempty" redis:"target"`       //joined or left user
	GroupName      string            `json:"groupName,omitempty" redis:"groupName"` //as mission's name
	From           string            `json:"from,omitempty" redis:"from"`
	SenderId       int               `json:"-" redis:"-"` //server side use
	To             []string          `json:"to,omitempty" redis:"to"`
	ReceiversId    []int             `json:"-" redis:"-"` //server side use
	Timestamp      int64             `json:"timestamp,omitempty" redis:"timestamp"`
	ExtraInfo      map[string]string `json:"extraInfo" redis:"-"`
	SerializedInfo string            `json:"-" redis:"info"` //serialize from extra info
}

// Response deals with the chat message
func (ct *Chat) Response() {
	switch ct.Type {
	//notify the special user
	case InvitedToProject, KickedFromProject:
		uid := dbms.ReadUserId(ct.Target)
		ct.ReceiversId = append(ct.ReceiversId, uid)

	//notify one user
	case PeerToPeer:
		uid := dbms.ReadUserId(ct.To[0])
		ct.ReceiversId = append(ct.ReceiversId, uid)

	//notify other members in this conversation
	case UserJoined, UserLeft, Discuss:
		ids := readChatMembers(ct)
		ct.ReceiversId = ids[:0]
		for i, uid := range ids {
			if uid == ct.SenderId {
				ct.ReceiversId = append(ct.ReceiversId, ids[i+1:]...)
				break
			} else {
				ct.ReceiversId = append(ct.ReceiversId, uid)
			}
		}
	case InvitedToMission:
		ct.ReceiversId = make([]int, len(ct.To))
		for i, uPid := range ct.To {
			ct.ReceiversId[i] = dbms.ReadUserId(uPid)
		}
		log.Debug("receivers:", ct.ReceiversId)
	}
	ct.send()
}

//send with socket or apple push
func (ct *Chat) send() {
	ct.Timestamp = time.Now().Unix()
	ct.Save()
	var userTokens []string
	for _, uid := range ct.ReceiversId {
		userTokens = append(userTokens, base.MakeToken(uid))
		go ct.AddToOffline(uid)
	}
	offlineTokens := rpc.SocketPush(userTokens, ct.Msg, ct.ExtraInfo) //use webSocket push

	//if ct.Type != UserJoined && ct.Type != UserLeft {
	applePush(offlineTokens, ct)
	//}
}

//read user's device token and push
func applePush(tokens []string, ct *Chat) {
	var deviceList []string
	//read device token from db
	for _, token := range tokens {
		uid, _ := base.ParseToken(token)
		dt := dbms.ReadDeviceToken(uid)
		if dt != "" {
			deviceList = append(deviceList, dt)
		}
	}
	log.Debug("dt:", deviceList)
	rpc.AppPush(deviceList, ct.Msg, ct.ExtraInfo)
}

// Save saves the offline chat message
func (ct *Chat) Save() {
	data, _ := json.Marshal(ct.ExtraInfo)
	ct.SerializedInfo = string(data)
	//save chat
	ct.Id = createChat(ct)
}

func (ct *Chat) AddToOffline(uid int) {
	//update user offline msg
	createOfflineMsg(uid, ct.Id)
}

// Load from redis
func (ct *Chat) Load() (err error) {
	cPtr, err := readChatWithId(ct.Id)
	if err != nil {
		log.Error(err)
	}
	*ct = *cPtr
	return
}

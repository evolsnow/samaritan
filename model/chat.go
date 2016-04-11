package model

import (
	"github.com/evolsnow/samaritan/common/base"
	"github.com/evolsnow/samaritan/common/dbms"
	"github.com/evolsnow/samaritan/common/rpc"
	"time"
)

const (
	PeerToPeer = iota //private msg server <<-->> client
	Discuss           //group chat server <<-->> client
	//system call
	UserJoined        // server -->> client
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
}

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
	}
	ct.send()
}

func (ct *Chat) send() {
	ct.Timestamp = time.Now().Unix()
	var userTokens []string
	for _, uid := range ct.ReceiversId {
		userTokens = append(userTokens, base.MakeToken(uid))
	}
	offlineTokens := rpc.SocketPush(userTokens, ct.Msg, ct.ExtraInfo) //use webSocket push
	for _, ft := range offlineTokens {
		uid, _ := base.ParseToken(ft)
		go ct.Save(uid)
	}
	if ct.Type == PeerToPeer || ct.Type == Discuss {
		applePush(offlineTokens, ct)
	}
}

func applePush(tokens []string, ct *Chat) {
	deviceList := make([]string, len(tokens), len(tokens))
	//read device token from db
	for _, token := range tokens {
		uid, _ := base.ParseToken(token)
		dt, _ := readDeviceToken(uid)
		deviceList = append(deviceList, dt)
	}
	rpc.IOSPush(deviceList, ct.Msg, ct.ExtraInfo)
}

func (ct *Chat) Save(uid int) {
	if ct.Id == 0 {
		//not saved
		ct.Id = createChat(ct)
	} else {
		//offline msg saved
		updateOfflineMsg(uid, ct.Id)
	}
}

package model

import (
	"github.com/evolsnow/samaritan/common/rpc"
	"github.com/gorilla/websocket"
	"time"
)

const (
	PeerToPeer = iota //private msg server <<-->> client
	Discuss           //group chat server <<-->> client
	//system call
	UserJoined        // server -->> client
	UserLeft          //server -->> client
	InvitedToMission  //server -->> client
	KickedFromMission //server -->> client
)

type Chat struct {
	Id             int      `json:"-" redis:"id"`
	ConversationId string   `json:"convId,omitempty" redis:"convId"`
	Type           int      `json:"type" redis:"type"`
	Msg            string   `json:"msg,omitempty" redis:"msg"`
	Target         string   `json:"target,omitempty" redis:"target"`       //joined or left user
	GroupName      string   `json:"groupName,omitempty" redis:"groupName"` //as mission's name
	From           string   `json:"from,omitempty" redis:"from"`
	SenderId       int      `json:"-" redis:"-"` //server side use
	To             []string `json:"to" redis:"to"`
	ReceiversId    []int    `json:"-" redis:"-"` //server side use
	Timestamp      int64    `json:"timestamp" redis:"timestamp"`
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

func (ct *Chat) Response(scm map[int]*websocket.Conn, dm map[int]string) {
	switch ct.Type {
	//notify the special user
	case InvitedToMission, KickedFromMission:
		uid := ReadUserId(ct.Target)
		ct.ReceiversId = append(ct.ReceiversId, uid)

	//notify one user
	case PeerToPeer:
		uid := ReadUserId(ct.To[0])
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
	}
	ct.send(scm, dm)
}

func (ct *Chat) send(scm map[int]*websocket.Conn, dm map[int]string) {
	ct.Timestamp = time.Now().Unix()
	offlineIds := make([]int, len(ct.ReceiversId), len(ct.ReceiversId))
	for _, uid := range ct.ReceiversId {
		sc, ok := scm[uid]
		if ok {
			go func(*websocket.Conn, int) {
				if err := sc.WriteJSON(ct); err != nil {
					offlineIds = append(offlineIds, uid)
				}
			}(sc, uid)
		} else {
			go ct.Save(uid)
			offlineIds = append(offlineIds, uid)
		}
	}
	if ct.Type == PeerToPeer || ct.Type == Discuss {
		applePush(offlineIds, ct, dm)
	}
}

func applePush(ids []int, ct *Chat, dm map[int]string) {
	deviceList := make([]string, len(ids), len(ids))
	for _, uid := range ids {
		token, ok := dm[uid]
		if !ok {
			//load from redis
			token, _ = readDeviceToken(uid)
			dm[uid] = token
		}
		deviceList = append(deviceList, token)
	}
	rpc.IOSPush(deviceList, ct.Msg)
}

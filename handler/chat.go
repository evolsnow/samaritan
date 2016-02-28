package handler

import (
	"encoding/json"
	"github.com/anachronistic/apns"
	"github.com/evolsnow/httprouter"
	"github.com/evolsnow/samaritan/model"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
)

const (
	PeerToPeer = iota //private msg server <<-->> client
	Discuss           //group chat server <<-->> client
	//system call
	UserJoined      // server -->> client
	UserLeft        //server -->> client
	InvitedToGroup  //server -->> client
	KickedFromGroup //server -->> client
)

type Chat struct {
	Id             int      `json:"-" redis:"id"`
	ConversationId string   `json:"convId" redis:"convId"`
	Type           int      `json:"type" redis:"type"`
	Msg            string   `json:"msg,omitempty" redis:"msg"`
	Target         string   `json:"target,omitempty" redis:"target"` //joined or left user
	GroupName      string   `json:"groupName,omitempty" redis:"groupName"`
	From           string   `json:"from,omitempty" redis:"from"`
	SenderId       int      `json:"-" redis:"-"` //server side use
	To             []string `json:"to" redis:"to"`
	ReceiversId    []int    `json:"-" redis:"-"` //server side use
	Timestamp      int64    `json:"timestamp" redis:"timestamp"`
}

var upgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool {
	return true
}} // use default options for webSocket

var socketConnMap = make(map[int]*websocket.Conn)
var deviceMap = make(map[int]string)

//keep deviceToken and connection
func Socket(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	uid, _ := strconv.Atoi(ps.Get("userId"))
	deviceToken := ps.ByName("deviceToken")
	deviceMap[uid] = deviceToken
	establishSocketConn(w, r, uid)
}

func establishSocketConn(w http.ResponseWriter, r *http.Request, uid int) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	socketConnMap[uid] = c
	log.Println("new socket conn:", uid)
	defer c.Close()
	defer delete(socketConnMap, uid)
	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("rect: %s", message)
		go handlerMsg(message)
	}
}

func handlerMsg(msg []byte) {
	ct := new(Chat)
	if err := json.Unmarshal(msg, ct); err == nil {
		ct.SenderId = model.ReadUserId(ct.From)
		ct.Response()
	}
}

func (ct *Chat) Response() {
	switch ct.Type {
	//notify the special user
	case InvitedToGroup, KickedFromGroup:
		uid := model.ReadUserId(ct.Target)
		ct.ReceiversId = append(ct.ReceiversId, uid)

	//notify other members in this conversation
	case UserJoined, UserLeft, PeerToPeer, Discuss:
		ids := readChatMembers(ct)
		ct.ReceiversId = ids[:0]
		for i, uid := range ids {
			if uid == ct.SenderId {
				ct.ReceiversId = append(ct.ReceiversId, ids[i:]...)
				break
			} else {
				ct.ReceiversId = append(ct.ReceiversId, uid)
			}
		}
	}
	ct.send()
}

func (ct *Chat) send() {
	ct.Timestamp = time.Now().Unix()
	offlineIds := make([]int, len(ct.ReceiversId), len(ct.ReceiversId))
	for _, uid := range ct.ReceiversId {
		sc, ok := socketConnMap[uid]
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
		applePush(offlineIds, ct)
	}
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

func applePush(ids []int, ct *Chat) {
	payload := apns.NewPayload()
	payload.Alert = ct.Msg
	payload.Sound = "default"
	payload.Badge = 1
	client := apns.NewClient("gateway.sandbox.push.apple.com:2195", "static/certs/cert.pem", "static/certs/key.pem")

	var wg sync.WaitGroup
	wg.Add(len(ids))
	for _, uid := range ids {
		token, ok := deviceMap[uid]
		if !ok {
			//load from redis
			token, _ = readDeviceToken(uid)
			deviceMap[uid] = token
		}
		pn := apns.NewPushNotification()
		pn.DeviceToken = token
		pn.AddPayload(payload)
		go func(*apns.PushNotification) {
			defer wg.Done()
			resp := client.Send(pn)
			if resp.Error != nil {
				log.Println("push notification error:", resp.Error)
			} else {
				log.Println("successfully push:", pn.DeviceToken)
			}
		}(pn)
	}
	wg.Wait()
}

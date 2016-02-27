package handler

import (
	"encoding/json"
	"github.com/anachronistic/apns"
	"github.com/evolsnow/httprouter"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type Conversation struct {
	Id             int      `json:"-" redis:"id"`
	ConversationId string   `json:"convId" redis:"convId"`
	Msg            string   `json:"msg" redis:"msg"`
	GroupName      string   `json:"private" redis:"groupName"`
	From           string   `json:"from" redis:"from"`
	SenderId       int      `json:"-" redis:"-"`
	To             []string `json:"to" redis:"to"`
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
		log.Printf("recv: %s", message)
		go handlerMsg(message)
	}
}

func handlerMsg(msg []byte) {
	cv := new(Conversation)
	if err := json.Unmarshal(msg, cv); err == nil {
		cv.SenderId, _ = readUserId(cv.From)
		cv.Execute()
	}
}

func (cv *Conversation) Execute() {
	ids := readConvIds(cv)
	reply := cv
	reply.Timestamp = time.Now().Unix()
	offlineIds := make([]int, len(ids), len(ids))
	for _, uid := range ids {
		if uid == cv.SenderId {
			continue
		}
		sc, ok := socketConnMap[uid]
		if ok {
			go func(*websocket.Conn, int) {
				if err := sc.WriteJSON(reply); err != nil {
					offlineIds = append(offlineIds, uid)
				}
			}(sc, uid)
		} else {
			go reply.Save(uid)
			offlineIds = append(offlineIds, uid)
		}
	}
	applePush(offlineIds, reply)
}

func (cv *Conversation) Save(uid int) {
	if cv.Id == 0 {
		//not saved
		cv.Id = createConversation(cv)
	} else {
		//offline msg saved
		updateOfflineMsg(uid, cv.Id)
	}
}

func applePush(ids []int, cv *Conversation) {
	payload := apns.NewPayload()
	payload.Alert = cv.Msg
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

package handler

import (
	"encoding/json"
	"github.com/evolsnow/httprouter"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"strconv"
	"time"
)

type Conversation struct {
	Id             int      `json:"-" redis:"id"`
	ConversationId string   `json:"convId" redis:"convId"`
	Msg            string   `json:"msg" redis:"msg"`
	GroupName      string   `json:"private" redis:"groupName"`
	From           string   `json:"from" redis:"from"`
	To             []string `json:"to,omitempty" redis:"to"`
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
	if err := json.Unmarshal(msg, cv); err != nil {
		cv.Execute()
	}
}

func (cv *Conversation) Execute() {
	ids := readConvIds(cv)
	reply := cv
	reply.Timestamp = time.Now().Unix()
	for _, uid := range ids {
		sc, ok := socketConnMap[uid]
		if ok {
			go sc.WriteJSON(reply)
		} else {
			applePush()
			go reply.Offline(uid)
		}
	}
}

func (cv *Conversation) Offline(uid int) {
	if cv.Id == 0 {
		//not saved
		cv.Id = createConversation(cv)
	} else {
		//offline msg saved
		updateOfflineMsg(uid, cv.Id)
	}
}
func applePush() {}

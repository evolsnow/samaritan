package handler

import (
	"encoding/json"
	"github.com/evolsnow/httprouter"
	"github.com/evolsnow/samaritan/common/log"
	"github.com/evolsnow/samaritan/model"
	"github.com/gorilla/websocket"
	"net/http"
	"strconv"
)

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
		log.Warn("upgrade:", err)
		return
	}
	socketConnMap[uid] = c
	log.Info("new socket conn:", uid)
	defer c.Close()
	defer delete(socketConnMap, uid)
	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Warn("read:", err)
			break
		}
		log.Debug("rec: %s", message)
		go handlerMsg(message)
	}
}

func handlerMsg(msg []byte) {
	ct := new(model.Chat)
	if err := json.Unmarshal(msg, ct); err == nil {
		ct.SenderId = model.ReadUserId(ct.From)
		ct.Response(socketConnMap, deviceMap)
	}
}

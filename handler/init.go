package handler

import (
	"encoding/json"
	"github.com/evolsnow/samaritan/common/caches"
	"github.com/evolsnow/samaritan/common/rpc"
	"github.com/evolsnow/samaritan/model"
)

//get cache
var cache = caches.NewCache()

func init() {
	go func() {
		for {
			msg := <-rpc.Chats
			go handlerMsg([]byte(msg))
		}
	}()
}

func handlerMsg(msg []byte) {
	ct := new(model.Chat)
	if err := json.Unmarshal(msg, ct); err == nil {
		ct.Response()
	}
}

package handler

import (
	"fmt"
	"github.com/evolsnow/samaritan/conn"
	"github.com/evolsnow/samaritan/model"
	"github.com/garyburd/redigo/redis"
	"log"
)

const (
	CId        = "id"
	CConvId    = "convId"
	CMsg       = "msg"
	CGroupName = "groupName"
	CFrom      = "from"
)

const (
	//clientId = "clientId:" //index for userId, clientId:John return john's userId
	clientId       = model.ClientId
	privateChat    = "privateChat:"
	publicChat     = "publicChat:"
	offlineMsgList = "user:%d:offlineMsg" //redis type:list
)

func readUserId(user string) (uid int, err error) {
	c := conn.Pool.Get()
	defer c.Close()
	key := clientId + user
	uid, err = redis.Int(c.Do("GET", key))
	return
}

func createPrivateConvRecord(chatId string, ua, ub int) {
	c := conn.Pool.Get()
	defer c.Close()
	c.Send("SADD", privateChat+chatId, ua)
	c.Send("SADD", privateChat+chatId, ub)
	c.Flush()
}

func readConvIds(cv *Conversation) (ids []int) {
	c := conn.Pool.Get()
	defer c.Close()
	key := ""
	if cv.GroupName == "" {
		key = privateChat + cv.ConversationId
	} else {
		key = publicChat + cv.ConversationId
	}
	ids, _ = redis.Ints(c.Do("SMEMBERS", key))
	return
}

func createConversation(cv *Conversation) int {
	c := conn.Pool.Get()
	defer c.Close()
	lua := `
	local cid = redis.call("INCR", "autoIncrConv")
	redis.call("HMSET", "conv:"..cid,
					KEYS[1], cid, KEYS[3], KEYS[4], KEYS[5], KEYS[6],
					KEYS[7], KEYS[8], KEYS[9], KEYS[10],)
	rerurn cid
	`
	k1, k2 := CId, cv.Id
	k3, k4 := CConvId, cv.ConversationId
	k5, k6 := CMsg, cv.Msg
	k7, k8 := CGroupName, cv.GroupName
	k9, k10 := CFrom, cv.From
	script := redis.NewScript(10, lua)
	id, err := redis.Int(script.Do(c, k1, k2, k3, k4, k5, k6, k7, k8, k9, k10))
	if err != nil {
		log.Println("Error create offline conversation:", err)
	}
	return id
}

func updateOfflineMsg(uid, convId int) {
	c := conn.Pool.Get()
	defer c.Close()
	key := fmt.Sprintf(offlineMsgList, uid)
	c.Do("RPUSH", key, convId)
}

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
	CType      = "type"
	CTarget    = "target"
	CMsg       = "msg"
	CGroupName = "groupName"
	CFrom      = "from"
)

const (
	//clientId = "clientId:" //index for userId, clientId:John return john's userId
	clientId       = model.ClientId
	deviceToken    = "deviceToken:%d" //ios device token
	privateChat    = "p2pChat:"
	offlineMsgList = "user:%d:offlineMsg" //redis type:list
)

func readUserId(user string) (uid int, err error) {
	c := conn.Pool.Get()
	defer c.Close()
	key := clientId + user
	uid, err = redis.Int(c.Do("GET", key))
	return
}

func readDeviceToken(uid int) (token string, err error) {
	c := conn.Pool.Get()
	defer c.Close()
	key := fmt.Sprintf(deviceToken, uid)
	token, err = redis.String(c.Do("GET", key))
	return
}

func createPrivateChatRecord(chatId string, ua, ub int) {
	c := conn.Pool.Get()
	defer c.Close()
	c.Send("SADD", privateChat+chatId, ua)
	c.Send("SADD", privateChat+chatId, ub)
	c.Flush()
}

func readChatMembers(ct *Chat) (ids []int) {
	c := conn.Pool.Get()
	defer c.Close()
	key := ""
	if ct.Type == PeerToPeer {
		key = privateChat + ct.ConversationId
		ids, _ = redis.Ints(c.Do("SMEMBERS", key))
	} else {
		//read group members
		p := new(model.Project)
		p.Name = ct.GroupName
		ids = p.GetMembersId()
	}
	return
}

func createChat(ct *Chat) int {
	c := conn.Pool.Get()
	defer c.Close()
	//todo expire the msg
	lua := `
	local cid = redis.call("INCR", "autoIncrChat")
	redis.call("HMSET", "chat:"..cid,
					KEYS[1], cid, KEYS[3], KEYS[4], KEYS[5], KEYS[6],
					KEYS[7], KEYS[8], KEYS[9], KEYS[10],
					KEYS[11], KEYS[12], KEYS[13], KEYS[14])
	rerurn cid
	`
	ka := []interface{}{
		CId, ct.Id,
		CConvId, ct.ConversationId,
		CType, ct.Type,
		CTarget, ct.Target,
		CMsg, ct.Msg,

		CGroupName, ct.GroupName,
		CFrom, ct.From,
	}
	script := redis.NewScript(len(ka), lua)
	id, err := redis.Int(script.Do(c, ka...))
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

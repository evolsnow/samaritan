package model

import (
	"github.com/evolsnow/samaritan/conn"
	"github.com/garyburd/redigo/redis"
	"log"
	"strconv"
)

//user redis name

const (
	UserId       = "id"
	UserAlias    = "alias"
	UserName     = "name"
	UserPhone    = "phone"
	UserPassword = "passwd"
	UserAvatar   = "avatar"
	UserSchool   = "school"
	UserDep      = "depart"
	UserGrade    = "grade"
	UserClass    = "class"
	UserStuNum   = "stuNum"
)

//to-do thing redis name

const (
	TodoId        = "id"
	TodoStartTime = "startTime"
	TodoDeadLine  = "deadLine"
	TodoDesc      = "desc"
	TodoOwnerId   = "ownerId"
	TodoStatus    = "status"
	TodoMissionId = "missionId"
)

//mission redis name
const (
	MissionId          = "id"
	MissionStartTime   = "startTime"
	MissionDeadLine    = "deadLine"
	MissionDesc        = "desc"
	MissionPublisherId = "publisherId"
	MissionColor       = "color"     //Color saved as redis-list -> mission:{id}:color [00,00,00]
	MissionRcv         = "receivers" //Receivers saved as redis-set -> mission:{id}:receivers (1,3,2)
)

//All model's redis actions

func getPassword(id int) string {
	c := conn.Pool.Get()
	defer c.Close()
	user := "user:" + strconv.Itoa(id)
	pwd, _ := redis.String(c.Do("HGET", user, UserPassword))
	return pwd
}

func updatePassword(id int, pwd string) {
	c := conn.Pool.Get()
	defer c.Close()
	user := "user:" + strconv.Itoa(id)
	_, err := c.Do("HSET", user, UserPassword, pwd)
	if err != nil {
		log.Println("Failed to update password for user:%s", id)
	}
}

func readUser(id int) (u *User, err error) {
	c := conn.Pool.Get()
	defer c.Close()
	user := "user:" + strconv.Itoa(id)
	ret, err := redis.Values(c.Do("HGETALL", user))
	if err != nil {
		return
	}
	err = redis.ScanStruct(ret, u)
	return
}

func readMissionRcv(mid int) (reply []*User, err error) {
	c := conn.Pool.Get()
	defer c.Close()
	receiverSet := "mission" + strconv.Itoa(mid) + MissionRcv
	multiGetUserLua := `
	local data = redis.call("SMEMBERS", KEYS[1])
	local ret = {}
  	for idx=1,#data do
  		ret[idx] = redis.call("HGETALL","user:"..data[idx])
  	end
  	return ret
   	`
	multiGetUserScript := redis.NewScript(0, multiGetUserLua)
	users, err := redis.Values(multiGetUserScript.Do(c, receiverSet))
	reply = make([]*User, len(users))
	for _, v := range users {
		rcv := new(User)
		err = redis.ScanStruct(v.([]interface{}), rcv)
		reply = append(reply, rcv)
	}
	return reply, err
}

func readMission(mid int) (reply []interface{}, err error) {
	c := conn.Pool.Get()
	defer c.Close()
	mission := "mission:" + strconv.Itoa(mid)
	reply, err = redis.Values(c.Do("HGETALL", mission))
	return
}

func createMissionRcv(mid, uid int) (err error) {
	c := conn.Pool.Get()
	defer c.Close()
	receiverSet := "mission:" + strconv.Itoa(mid) + MissionRcv
	_, err = c.Do("SADD", receiverSet, uid)
	return
}

func createMission(m *Mission) {

}

package model

import (
	"fmt"
	"github.com/evolsnow/samaritan/conn"
	"github.com/garyburd/redigo/redis"
	"strconv"
)

//user redis key name
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

//to-do thing redis key name
const (
	TodoId           = "id"
	TodoStartTime    = "startTime"
	TodoDeadline     = "deadline"
	TodoDesc         = "desc"
	TodoOwnerId      = "ownerId"
	TodoAccomplished = "accomplished"
	TodoMissionId    = "missionId"
)

//mission redis key name
const (
	MissionId          = "id"
	MissionStartTime   = "startTime"
	MissionDesc        = "desc"
	MissionPublisherId = "publisherId"

//Receivers saved as redis-set -> mission:{id}:rcv (1,3,2)
)

//other useful index set key name
const (
	userBelongSet              = "%s:%s:%d:%s"          //just for further analysis-> school:department:grade:class
	userTodoList               = "user:%d:todoList"     //user's all to-do, redis-type:List
	userTodoNotAccomplishedSet = "user:%d:todoStatus:0" //to-do status, redis-type:Set
	userTodoAccomplishedSet    = "user:%d:todoStatus:1"

	userMissionJoinedSet    = "user:%d:missions:participate" //user's all missions redis-type:Set
	userMissionPublishedSet = "user:%d:missions:publish"
	userMissionColorList    = "user:%d:mission:%d:color" //user defined mission color redis-type:List

	missionRcvSet = "mission:%d:rcv" //mission's receivers redis-type:Set
)

//redis actions of model User

func createUser(u *User) error {
	c := conn.Pool.Get()
	defer c.Close()
	createUserLua := `
	local tid = KEYS[2]
	redis.call("HMSET", "user:"..uid,
					KEYS[1], KEYS[2], KEYS[3], KEYS[4], KEYS[5], KEYS[6], KEYS[7], KEYS[8],
					KEYS[9], KEYS[10], KEYS[11], KEYS[12], KEYS[13], KEYS[14], KEYS[15], KEYS[16],
					KEYS[17], KEYS[18], KEYS[19], KEYS[20], KEYS[21], KEYS[22])
	redis.call("SADD", KEYS[23], uid)
	`
	//to-do model
	u.Id, _ = redis.Int(c.Do("INCR", "autoIncrUser"))
	k1, k2 := UserId, u.Id
	k3, k4 := UserAlias, u.Alias
	k5, k6 := UserName, u.Name
	k7, k8 := UserPhone, u.Phone
	k9, k10 := UserPassword, u.Password
	k11, k12 := UserAvatar, u.Avatar
	k13, k14 := UserSchool, u.School
	k15, k16 := UserDep, u.Department
	k17, k18 := UserGrade, u.Grade
	k19, k20 := UserClass, u.Class
	k21, k22 := UserStuNum, u.StudentNum

	//redis set
	k23 := fmt.Sprintf(userBelongSet, u.School, u.Department, u.Grade, u.Class)

	createUserScript := redis.NewScript(23, createUserLua)
	_, err := createUserScript.Do(c, k1, k2, k3, k4, k5, k6, k7, k8, k9, k10, k11, k12,
		k13, k14, k15, k16, k17, k18, k19, k20, k21, k22, k23)
	return err
}

func ReadPassword(uid int) (pwd string, err error) {
	c := conn.Pool.Get()
	defer c.Close()
	user := "user:" + strconv.Itoa(uid)
	pwd, err = redis.String(c.Do("HGET", user, UserPassword))
	return
}

func updatePassword(id int, pwd string) error {
	c := conn.Pool.Get()
	defer c.Close()
	user := "user:" + strconv.Itoa(id)
	_, err := c.Do("HSET", user, UserPassword, pwd)
	return err
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

//redis actions of model to-do
func createTodo(td *Todo) error {
	c := conn.Pool.Get()
	defer c.Close()
	createTodoLua := `
	local tid = KEYS[2]
	redis.call("HMSET", "todo:"..tid,
					KEYS[1], KEYS[2], KEYS[3], KEYS[4], KEYS[5], KEYS[6], KEYS[7], KEYS[8],
					KEYS[9], KEYS[10], KEYS[11], KEYS[12], KEYS[13], KEYS[14])
	redis.call("RPUSH", KEYS[15], tid)
	redis.call("SADD", KEYS[16], tid)
	`
	//to-do model
	td.Id, _ = redis.Int(c.Do("INCR", "autoIncrTodo"))
	k1, k2 := TodoId, td.Id
	k3, k4 := TodoDesc, td.Desc
	k5, k6 := TodoStartTime, td.StartTime
	k7, k8 := TodoDeadline, td.Deadline
	k9, k10 := TodoAccomplished, td.Accomplished
	k11, k12 := TodoOwnerId, td.OwnerId
	k13, k14 := TodoMissionId, td.MissionId
	//redis list
	k15 := fmt.Sprintf(userTodoList, td.OwnerId)
	k16 := fmt.Sprintf(userTodoNotAccomplishedSet, td.OwnerId)

	createTodoScript := redis.NewScript(16, createTodoLua)
	_, err := createTodoScript.Do(c, k1, k2, k3, k4, k5, k6, k7, k8, k9, k10, k11, k12, k13, k14, k15, k16)
	return err
}

func updateTodoStatus(uid, tid int) error {
	c := conn.Pool.Get()
	defer c.Close()
	accomplished := fmt.Sprintf(userTodoAccomplishedSet, uid)
	notAccomplished := fmt.Sprintf(userTodoNotAccomplishedSet, uid)
	_, err := c.Do("SMOVE", notAccomplished, accomplished, tid)
	return err
}

//redis actions of model mission
func createMission(m *Mission) error {
	c := conn.Pool.Get()
	defer c.Close()
	createMissionLua := `
	local mid = KEYS[2]
	redis.call("HMSET", "mission:"..mid,
					   KEYS[1], KEYS[2], KEYS[3], KEYS[4], KEYS[5], KEYS[6],
					   KEYS[7], KEYS[8],KEYS[9], KEYS[10])
	redis.call("SADD", KEYS[11], mid)
	redis.call("SADD", KEYS[12], mid)
	`
	//mission models
	m.Id, _ = redis.Int(c.Do("INCR", "autoIncrMission"))
	k1, k2 := MissionId, m.Id
	k3, k4 := MissionStartTime, m.StartTime
	//k5, k6 := MissionDeadline, m.Deadline
	k7, k8 := MissionDesc, m.Desc
	k9, k10 := MissionPublisherId, m.PublisherId
	//redis set
	k11 := fmt.Sprintf(userMissionJoinedSet, m.PublisherId)
	k12 := fmt.Sprintf(userMissionPublishedSet, m.PublisherId)

	createMissionScript := redis.NewScript(10, createMissionLua)
	_, err := createMissionScript.Do(c, k1, k2, k3, k4, k7, k8, k9, k10, k11, k12)
	return err
}

func readMissionRcv(mid int) (reply []*User, err error) {
	c := conn.Pool.Get()
	defer c.Close()
	rcvSet := fmt.Sprintf(missionRcvSet, mid)
	multiGetUserLua := `
	local data = redis.call("SMEMBERS", KEYS[1])
	local ret = {}
  	for idx=1,#data do
  		ret[idx] = redis.call("HGETALL","user:"..data[idx])
  	end
  	return ret
   	`
	multiGetUserScript := redis.NewScript(0, multiGetUserLua)
	users, err := redis.Values(multiGetUserScript.Do(c, rcvSet))
	reply = make([]*User, len(users))
	for _, v := range users {
		rcv := new(User)
		err = redis.ScanStruct(v.([]interface{}), rcv)
		reply = append(reply, rcv)
	}
	return reply, err
}

func readMission(mid int) (m *Mission, err error) {
	c := conn.Pool.Get()
	defer c.Close()
	mission := "mission:" + strconv.Itoa(mid)
	ret, err := redis.Values(c.Do("HGETALL", mission))
	if err != nil {
		return
	}
	err = redis.ScanStruct(ret, m)
	return
}

func createMissionRcv(mid, uid int) (err error) {
	c := conn.Pool.Get()
	defer c.Close()
	rcvSet := fmt.Sprintf(missionRcvSet, mid)
	_, err = c.Do("SADD", rcvSet, uid)
	return
}

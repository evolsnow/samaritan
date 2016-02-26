package model

import (
	"fmt"
	"github.com/evolsnow/samaritan/conn"
	"github.com/garyburd/redigo/redis"
	"log"
	"strconv"
)

//user redis key name
const (
	UId       = "id"
	UAlias    = "alias"
	UName     = "name"
	UPhone    = "phone"
	UPassword = "passwd"
	UAvatar   = "avatar"
	USchool   = "school"
	UDep      = "depart"
	UGrade    = "grade"
	UClass    = "class"
	UStuNum   = "stuNum"
)

//to-do thing redis key name
const (
	TId        = "id"
	TStartTime = "startTime"
	TDeadline  = "deadline"
	TDesc      = "desc"
	TOwnerId   = "ownerId"
	TDone      = "done"
	TMissionId = "missionId"
)

//mission redis key name
const (
	MId        = "id"
	MStartTime = "startTime"
	MDesc      = "desc"
	MPubId     = "publisherId"

//Receivers saved as redis-set -> mission:{id}:rcv (1,3,2)
)

//other useful index set key name
const (
	userGroup        = "%s:%s:%d:%s"          //just for further analysis-> school:department:grade:class
	userTdList       = "user:%d:todoList"     //user's all to-do, redis-type:List
	userTdNotDoneSet = "user:%d:todoStatus:0" //to-do status, redis-type:Set
	userTdDoneSet    = "user:%d:todoStatus:1"

	userMsJoinedSet    = "user:%d:missions:participate" //user's all missions redis-type:Set
	userMsPublishedSet = "user:%d:missions:publish"
	userMsColorList    = "user:%d:mission:%d:color" //user defined mission color redis-type:List

	missionRcvSet = "mission:%d:rcv" //mission's receivers redis-type:Set
)

//redis actions of model User
func createUser(u *User) int {
	c := conn.Pool.Get()
	lua := `
	local uid = KEYS[2]
	redis.call("HMSET", "user:"..uid,
					KEYS[1], KEYS[2], KEYS[3], KEYS[4], KEYS[5], KEYS[6], KEYS[7], KEYS[8],
					KEYS[9], KEYS[10], KEYS[11], KEYS[12], KEYS[13], KEYS[14], KEYS[15], KEYS[16],
					KEYS[17], KEYS[18], KEYS[19], KEYS[20], KEYS[21], KEYS[22])
	redis.call("SADD", KEYS[23], uid)
	`
	u.Id, _ = redis.Int(c.Do("INCR", "autoIncrUser"))
	//return user's id asap
	go func() {
		//user model
		k1, k2 := UId, u.Id
		k3, k4 := UAlias, u.Alias
		k5, k6 := UName, u.Name
		k7, k8 := UPhone, u.Phone
		k9, k10 := UPassword, u.Password
		k11, k12 := UAvatar, u.Avatar
		k13, k14 := USchool, u.School
		k15, k16 := UDep, u.Department
		k17, k18 := UGrade, u.Grade
		k19, k20 := UClass, u.Class
		k21, k22 := UStuNum, u.StudentNum

		//redis set
		k23 := fmt.Sprintf(userGroup, u.School, u.Department, u.Grade, u.Class)

		script := redis.NewScript(23, lua)
		_, err := script.Do(c, k1, k2, k3, k4, k5, k6, k7, k8, k9, k10, k11, k12,
			k13, k14, k15, k16, k17, k18, k19, k20, k21, k22, k23)
		c.Close()
		if err != nil {
			log.Println("Error create user:", err)
		}
	}()
	return u.Id
}

func createUserAvatar(uid int, avatarUrl string) error {
	c := conn.Pool.Get()
	defer c.Close()
	user := "user:" + strconv.Itoa(uid)
	_, err := c.Do("HSET", user, UAvatar, avatarUrl)
	return err
}

func readPassword(uid int) (pwd string, err error) {
	c := conn.Pool.Get()
	defer c.Close()
	user := "user:" + strconv.Itoa(uid)
	pwd, err = redis.String(c.Do("HGET", user, UPassword))
	return
}

func updatePassword(id int, pwd string) error {
	c := conn.Pool.Get()
	defer c.Close()
	user := "user:" + strconv.Itoa(id)
	_, err := c.Do("HSET", user, UPassword, pwd)
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
	lua := `
	local tid = KEYS[2]
	redis.call("HMSET", "todo:"..tid,
					KEYS[1], KEYS[2], KEYS[3], KEYS[4], KEYS[5], KEYS[6], KEYS[7], KEYS[8],
					KEYS[9], KEYS[10], KEYS[11], KEYS[12], KEYS[13], KEYS[14])
	redis.call("RPUSH", KEYS[15], tid)
	redis.call("SADD", KEYS[16], tid)
	`
	//to-do model
	td.Id, _ = redis.Int(c.Do("INCR", "autoIncrTodo"))
	k1, k2 := TId, td.Id
	k3, k4 := TDesc, td.Desc
	k5, k6 := TStartTime, td.StartTime
	k7, k8 := TDeadline, td.Deadline
	k9, k10 := TDone, td.Done
	k11, k12 := TOwnerId, td.OwnerId
	k13, k14 := TMissionId, td.MissionId
	//redis list
	k15 := fmt.Sprintf(userTdList, td.OwnerId)
	k16 := fmt.Sprintf(userTdNotDoneSet, td.OwnerId)

	script := redis.NewScript(16, lua)
	_, err := script.Do(c, k1, k2, k3, k4, k5, k6, k7, k8, k9, k10, k11, k12, k13, k14, k15, k16)
	return err
}

func updateTodoStatus(uid, tid int) error {
	c := conn.Pool.Get()
	defer c.Close()
	done := fmt.Sprintf(userTdDoneSet, uid)
	notDone := fmt.Sprintf(userTdNotDoneSet, uid)
	_, err := c.Do("SMOVE", notDone, done, tid)
	return err
}

//redis actions of model mission
func createMission(m *Mission) error {
	c := conn.Pool.Get()
	defer c.Close()
	lua := `
	local mid = KEYS[2]
	redis.call("HMSET", "mission:"..mid,
					   KEYS[1], KEYS[2], KEYS[3], KEYS[4], KEYS[5], KEYS[6],
					   KEYS[7], KEYS[8],KEYS[9], KEYS[10])
	redis.call("SADD", KEYS[11], mid)
	redis.call("SADD", KEYS[12], mid)
	`
	//mission models
	m.Id, _ = redis.Int(c.Do("INCR", "autoIncrMission"))
	k1, k2 := MId, m.Id
	k3, k4 := MStartTime, m.StartTime
	//k5, k6 := MissionDeadline, m.Deadline
	k7, k8 := MDesc, m.Desc
	k9, k10 := MPubId, m.PublisherId
	//redis set
	k11 := fmt.Sprintf(userMsJoinedSet, m.PublisherId)
	k12 := fmt.Sprintf(userMsPublishedSet, m.PublisherId)

	script := redis.NewScript(10, lua)
	_, err := script.Do(c, k1, k2, k3, k4, k7, k8, k9, k10, k11, k12)
	return err
}

func readMissionRcv(mid int) (reply []*User, err error) {
	c := conn.Pool.Get()
	defer c.Close()
	rcvSet := fmt.Sprintf(missionRcvSet, mid)
	lua := `
	local data = redis.call("SMEMBERS", KEYS[1])
	local ret = {}
  	for idx=1,#data do
  		ret[idx] = redis.call("HGETALL","user:"..data[idx])
  	end
  	return ret
   	`
	script := redis.NewScript(0, lua)
	users, err := redis.Values(script.Do(c, rcvSet))
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

package model

import (
	"fmt"
	"github.com/evolsnow/samaritan/base"
	"github.com/evolsnow/samaritan/conn"
	"github.com/garyburd/redigo/redis"
	"log"
	"strconv"
	"time"
)

//user redis key name
const (
	UId       = "id"
	UPid      = "pid"
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
	TId  = "id"
	TPid = "pid"

	TStartTime = "startTime"
	TDeadline  = "deadline"
	TDesc      = "desc"
	TOwnerId   = "ownerId"
	TDone      = "done"
	TProjectId = "projectId"
)

//project redis key name
const (
	PId          = "id"
	PPid         = "pid"
	PName        = "name"
	PCreateTime  = "createTime"
	PDesc        = "desc"
	PPublisherId = "publisherId"
)

//other useful index set key name
const (
	userGroup        = "%s:%s:%d:%s"          //just for further analysis-> school:department:grade:class
	userTdList       = "user:%d:todoList"     //user's all to-do, redis-type:List
	userTdNotDoneSet = "user:%d:todoStatus:0" //to-do status, redis-type:Set
	userTdDoneSet    = "user:%d:todoStatus:1"

	userPjJoinedSet    = "user:%d:projects:participate" //user's all projects redis-type:Set
	userPjPublishedSet = "user:%d:projects:publish"
	userPjColorList    = "user:%d:project:%d:color" //user defined project color redis-type:List

	projectMembersSet = "project:%d:members" //project's receivers redis-type:Set
	ClientId          = "clientId:"          //index for userId, clientId:John return john's userId

)

//redis actions of model User
func createUser(u *User) {
	c := conn.Pool.Get()
	u.Id, _ = redis.Int(c.Do("INCR", "autoIncrUser"))
	u.Pid = base.HashedUserId(u.Id)
	//return user's id asap
	go func() {
		lua := `
			local uid = KEYS[2]
			redis.call("HMSET", "user:"..uid,
					KEYS[1], KEYS[2], KEYS[3], KEYS[4], KEYS[5], KEYS[6], KEYS[7], KEYS[8],
					KEYS[9], KEYS[10], KEYS[11], KEYS[12], KEYS[13], KEYS[14], KEYS[15], KEYS[16],
					KEYS[17], KEYS[18], KEYS[19], KEYS[20], KEYS[21], KEYS[22],
					KEYS[23], KEYS[24])
			redis.call("SADD", KEYS[25], uid)
			redis.call("SET", KEYS[26], uid)
			`
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
		k23, k24 := UPid, u.Pid
		//redis set
		//todo not now (::0:)?
		k25 := fmt.Sprintf(userGroup, u.School, u.Department, u.Grade, u.Class)
		k26 := ClientId + u.Name
		script := redis.NewScript(26, lua)
		_, err := script.Do(c, k1, k2, k3, k4, k5, k6, k7, k8, k9, k10, k11, k12,
			k13, k14, k15, k16, k17, k18, k19, k20, k21, k22, k23, k24, k25, k26)
		c.Close()
		if err != nil {
			log.Println("Error create user:", err)
		}
	}()
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
func createTodo(td *Todo) {
	c := conn.Pool.Get()
	td.Id, _ = redis.Int(c.Do("INCR", "autoIncrTodo"))
	td.Pid = base.HashedTodoId(td.Id)
	go func() {
		lua := `
			local tid = redis.call("INCR", "autoIncrTodo")
			redis.call("HMSET", "todo:"..tid,
					KEYS[1], tid, KEYS[3], KEYS[4], KEYS[5], KEYS[6], KEYS[7], KEYS[8],
					KEYS[9], KEYS[10], KEYS[11], KEYS[12], KEYS[13], KEYS[14],
					KEYS[15], KEYS[16])
			redis.call("RPUSH", KEYS[17], tid)
			redis.call("SADD", KEYS[18], tid)
			`
		//to-do model
		k1, k2 := TId, td.Id
		k3, k4 := TDesc, td.Desc
		k5, k6 := TStartTime, td.StartTime
		k7, k8 := TDeadline, td.Deadline
		k9, k10 := TDone, td.Done
		k11, k12 := TOwnerId, td.OwnerId
		k13, k14 := TProjectId, td.ProjectId
		k15, k16 := TPid, td.Pid
		//redis list
		k17 := fmt.Sprintf(userTdList, td.OwnerId)
		k18 := fmt.Sprintf(userTdNotDoneSet, td.OwnerId)

		script := redis.NewScript(18, lua)
		_, err := script.Do(c, k1, k2, k3, k4, k5, k6, k7, k8, k9, k10, k11, k12, k13, k14, k15, k16, k17, k18)
		if err != nil {
			log.Println("Error create todo:", err)
		}
		c.Close()
	}()
}

func updateTodoStatus(uid, tid int) error {
	c := conn.Pool.Get()
	defer c.Close()
	done := fmt.Sprintf(userTdDoneSet, uid)
	notDone := fmt.Sprintf(userTdNotDoneSet, uid)
	_, err := c.Do("SMOVE", notDone, done, tid)
	return err
}

//redis actions of model project
func createProject(p *Project) {
	c := conn.Pool.Get()
	p.Id, _ = redis.Int(c.Do("INCR", "autoIncrProject"))
	p.Pid = base.HashedProjectId(p.Id)
	go func() {
		lua := `
			local pid = KEYS[2]
			redis.call("HMSET", "project:"..pid,
					   KEYS[1], pid, KEYS[3], KEYS[4], KEYS[5], KEYS[6],
					   KEYS[7], KEYS[8], KEYS[9], KEYS[10], KEYS[11], KEYS[12])
			redis.call("SADD", KEYS[13], pid)
			redis.call("SADD", KEYS[14], pid)
			`
		//project models
		k1, k2 := PId, p.Id
		k3, k4 := PCreateTime, time.Now().Unix()
		k5, k6 := PDesc, p.Desc
		k7, k8 := PPublisherId, p.PublisherId
		k9, k10 := PName, p.Name
		k11, k12 := PPid, p.Pid
		//redis set
		k13 := fmt.Sprintf(userPjJoinedSet, p.PublisherId)
		k14 := fmt.Sprintf(userPjPublishedSet, p.PublisherId)

		script := redis.NewScript(14, lua)
		_, err := script.Do(c, k1, k2, k3, k4, k5, k6, k7, k8, k9, k10, k11, k12, k13, k14)
		if err != nil {
			log.Println("Error create project:", err)
		}
		c.Close()
	}()
}

func readProjectMembers(pid int) (reply []*User, err error) {
	c := conn.Pool.Get()
	defer c.Close()
	rcvSet := fmt.Sprintf(projectMembersSet, pid)
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

func readProject(pid int) (p *Project, err error) {
	c := conn.Pool.Get()
	defer c.Close()
	project := "project:" + strconv.Itoa(pid)
	ret, err := redis.Values(c.Do("HGETALL", project))
	if err != nil {
		return
	}
	err = redis.ScanStruct(ret, p)
	return
}

func updateProjectMember(pid, uid, action int) (err error) {
	c := conn.Pool.Get()
	defer c.Close()
	memSet := fmt.Sprintf(projectMembersSet, pid)
	if action > 0 {
		_, err = c.Do("SADD", memSet, uid)

	} else {
		_, err = c.Do("SREM", memSet, uid)
	}
	return
}

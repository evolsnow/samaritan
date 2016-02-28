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
	UEmail    = "email"
	UAvatar   = "avatar"
	USchool   = "school"
	UDep      = "depart"
	UGrade    = "grade"
	UClass    = "class"
	UStuNum   = "stuNum"
)

//to-do thing redis key name
const (
	TId           = "id"
	TPid          = "pid"
	TStartTime    = "startTime"
	TTaskTime     = "taskTime"
	TPlace        = "place"
	TRepeat       = "repeat"
	TRepeatPeriod = "repeatPeriod"
	TDesc         = "desc"
	TOwnerId      = "ownerId"
	TDone         = "done"
	TFinishTime   = "finishTime"
	TProjectId    = "projectId"
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
	ClientId          = "clientId:"          //index for userId, ClientId:john's token return john's userId
	GroupId           = "groupId:"           //index for groupId, GroupId:a return group a's Id

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
					KEYS[23], KEYS[24], KEYS[25], KEYS[26])
			redis.call("SADD", KEYS[27], uid)
			redis.call("SET", KEYS[28], uid)
			`
		ka := []interface{}{
			//user model
			UId, u.Id,
			UPid, u.Pid,
			UAlias, u.Alias,
			UName, u.Name,
			UPhone, u.Phone,
			UPassword, u.Password,
			UEmail, u.Email,
			UAvatar, u.Avatar,
			USchool, u.School,
			UDep, u.Department,
			UGrade, u.Grade,
			UClass, u.Class,
			UStuNum, u.StudentNum,
			//redis set
			//todo not now (::0:)?
			fmt.Sprintf(userGroup, u.School, u.Department, u.Grade, u.Class),
			ClientId + u.Pid,
		}
		script := redis.NewScript(len(ka), lua)
		_, err := script.Do(c, ka...)
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
			local tid = KEYS[2]
			redis.call("HMSET", "todo:"..tid,
					KEYS[1], tid, KEYS[3], KEYS[4], KEYS[5], KEYS[6], KEYS[7], KEYS[8],
					KEYS[9], KEYS[10], KEYS[11], KEYS[12], KEYS[13], KEYS[14],
					KEYS[15], KEYS[16], KEYS[17], KEYS[18], KEYS[19], KEYS[20],
					KEYS[21], KEYS[22], KEYS[23], KEYS[24])
			redis.call("RPUSH", KEYS[25], tid)
			redis.call("SADD", KEYS[26], tid)
			`
		ka := []interface{}{
			//to-do model
			TId, td.Id,
			TDesc, td.Desc,
			TStartTime, td.StartTime,
			TTaskTime, td.TaskTime,
			TPlace, td.Place,
			TRepeat, td.Repeat,
			TRepeatPeriod, td.RepeatPeriod,
			TDone, td.Done,
			TFinishTime, td.FinishTime,
			TOwnerId, td.OwnerId,
			TProjectId, td.ProjectId,
			TPid, td.Pid,
			//redis list
			fmt.Sprintf(userTdList, td.OwnerId),
			fmt.Sprintf(userTdNotDoneSet, td.OwnerId),
		}
		script := redis.NewScript(len(ka), lua)
		_, err := script.Do(c, ka...)
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
			redis.call("SET", KEYS[15], pid)
			`
		ka := []interface{}{
			//project models
			PId, p.Id,
			PCreateTime, time.Now().Unix(),
			PDesc, p.Desc,
			PPublisherId, p.PublisherId,
			PName, p.Name,
			PPid, p.Pid,
			//redis set
			fmt.Sprintf(userPjJoinedSet, p.PublisherId),
			fmt.Sprintf(userPjPublishedSet, p.PublisherId),
			GroupId + p.Name,
		}
		script := redis.NewScript(len(ka), lua)
		_, err := script.Do(c, ka...)
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

func readMemIdsWithName(name string) (ids []int, err error) {
	c := conn.Pool.Get()
	defer c.Close()
	idx := GroupId + name
	// warning: one hard code here
	lua := `
	local pid = redis.call("GET", KEYS[1])
	return redis.call("SMEMBERS", "project:"..pid..":members")
	`
	script := redis.NewScript(1, lua)
	ids, err = redis.Ints(script.Do(c, idx))
	return
}

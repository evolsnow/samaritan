package model

import (
	"fmt"
	"github.com/evolsnow/samaritan/common/base"
	"github.com/evolsnow/samaritan/common/dbms"
	"github.com/evolsnow/samaritan/common/log"
	"github.com/garyburd/redigo/redis"
	"strconv"
	"time"
)

//user redis key name
const (
	UId         = "id"
	UPid        = "pid"
	USamId      = "samId"
	UCreateTime = "createTime"
	UAlias      = "alias"
	UName       = "name"
	UPhone      = "phone"
	UPassword   = "password"
	UEmail      = "email"
	UAvatar     = "avatar"
	USchool     = "school"
	UDep        = "depart"
	UGrade      = "grade"
	UClass      = "class"
	UStuNum     = "stuNum"
)

//to-do thing redis key name
const (
	TId         = "id"
	TPid        = "pid"
	TCreateTime = "createTime"
	TStartTime  = "startTime"
	TPlace      = "place"
	TRepeat     = "repeat"
	TAllDay     = "allDay"
	TRepeatMode = "repeatMode"
	TDesc       = "desc"
	TRemark     = "remark"
	TOwnerId    = "ownerId"
	TDone       = "done"
	TFinishTime = "finishTime"
	TMissionId  = "missionId"
)

//mission redis key name
const (
	MId            = "id"
	MPid           = "pid"
	MName          = "name"
	MCreateTime    = "createTime"
	MDesc          = "desc"
	MPublisherId   = "publisherId"
	MCompletionNum = "completionNum"
	MCompletedTime = "completedTime"
	//comments
	CId         = "id"
	CPid        = "pid"
	CWhen       = "when"
	CCriticPid  = "criticPid"
	CCriticName = "criticName"
)

//project redis key name
const (
	PId         = "id"
	PPid        = "pid"
	PName       = "name"
	PCreateTime = "createTime"
	PDesc       = "desc"
	PCreatorId  = "creatorId"
	PPrivate    = "private"
)

const (
	ChId        = "id"
	ChConvId    = "convId"
	ChType      = "type"
	ChTarget    = "target"
	ChMsg       = "msg"
	ChGroupName = "groupName"
	ChFrom      = "from"
)

//other useful index set key name
const (
	//user
	userGroup          = "%s:%s:%d:%s"          //just for further analysis-> school:department:grade:class
	userTdList         = "user:%d:todoList"     //user's all to-do, redis-type:List
	userTdNotDoneSet   = "user:%d:todoStatus:0" //to-do status, redis-type:Set
	userTdDoneSet      = "user:%d:todoStatus:1"
	userPjJoinedSet    = "user:%d:projects:participate" //user's all projects redis-type:Set
	userPjCreatedSet   = "user:%d:projects:create"
	userMsAcceptedSet  = "user:%d:missions:accept" //user's all missions redis-type:Set
	userMsPublishedSet = "user:%d:missions:publish"
	userPjColorList    = "user:%d:project:%d:color" //user defined project color redis-type:List

	//to-do
	todoPictureList = "todo:%d:pictures" //to-do's pictures redis-type:List

	//mission
	missionReceiversSet = "mission:%d:receivers" //mission's receivers redis-type:Set
	missionCommentsList = "mission:%d:comments"  //mission's comments redis-type:List

	//project
	projectMembersSet = "project:%d:members" //project's members redis-type:Set

	//chat
	deviceToken    = "deviceToken:%d"     //ios device token
	offlineMsgList = "user:%d:offlineMsg" //redis type:list

)

//redis actions of model User
func createUser(u *User) {
	c := dbms.Pool.Get()
	u.Id, _ = redis.Int(c.Do("INCR", "autoIncrUser"))
	u.Pid = base.HashedUserId(u.Id)
	//index
	go dbms.CreateUserIndex(u.Id, u.Pid)
	//return user's id asap
	go func() {
		lua := `
			local uid = KEYS[2]
			redis.call("HMSET", "user:"..uid,
					KEYS[1], KEYS[2], KEYS[3], KEYS[4], KEYS[5], KEYS[6], KEYS[7], KEYS[8],
					KEYS[9], KEYS[10], KEYS[11], KEYS[12], KEYS[13], KEYS[14], KEYS[15], KEYS[16],
					KEYS[17], KEYS[18], KEYS[19], KEYS[20], KEYS[21], KEYS[22],
					KEYS[23], KEYS[24], KEYS[25], KEYS[26], KEYS[27], KEYS[28],
					KEYS[29], KEYS[30])
			redis.call("SADD", KEYS[31], uid)
			`
		ka := []interface{}{
			//user model
			UId, u.Id,
			UPid, u.Pid,
			USamId, u.SamId,
			UCreateTime, time.Now().Unix(),
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
		}
		script := redis.NewScript(len(ka), lua)
		_, err := script.Do(c, ka...)
		c.Close()
		if err != nil {
			log.Error("Error create user:", err)
		}
	}()
}

func readUser(id int) (*User, error) {
	c := dbms.Pool.Get()
	defer c.Close()
	user := "user:" + strconv.Itoa(id)
	ret, err := redis.Values(c.Do("HGETALL", user))
	if err != nil {
		return nil, err
	}
	u := new(User)
	err = redis.ScanStruct(ret, u)
	return u, err
}

func createUserAvatar(uid int, avatarUrl string) error {
	c := dbms.Pool.Get()
	defer c.Close()
	user := "user:" + strconv.Itoa(uid)
	_, err := c.Do("HSET", user, UAvatar, avatarUrl)
	return err
}

func readPassword(uid int) (pwd string, err error) {
	c := dbms.Pool.Get()
	defer c.Close()
	user := "user:" + strconv.Itoa(uid)
	pwd, err = redis.String(c.Do("HGET", user, UPassword))
	return
}

func updateUser(uid int, kvMap map[string]interface{}) error {
	c := dbms.Pool.Get()
	defer c.Close()
	for k, v := range kvMap {
		c.Send("HSET", "user:"+strconv.Itoa(uid), k, v)
	}
	return c.Flush()
}

//redis actions of model to-do
func createTodo(td *Todo) {
	c := dbms.Pool.Get()
	td.Id, _ = redis.Int(c.Do("INCR", "autoIncrTodo"))
	td.Pid = base.HashedTodoId(td.Id)
	go dbms.CreateTodoIndex(td.Id, td.Pid)
	go func() {
		lua := `
			local tid = KEYS[2]
			redis.call("HMSET", "todo:"..tid,
					KEYS[1], tid, KEYS[3], KEYS[4], KEYS[5], KEYS[6], KEYS[7], KEYS[8],
					KEYS[9], KEYS[10], KEYS[11], KEYS[12], KEYS[13], KEYS[14],
					KEYS[15], KEYS[16], KEYS[17], KEYS[18], KEYS[19], KEYS[20],
					KEYS[21], KEYS[22], KEYS[23], KEYS[24], KEYS[25], KEYS[26],
					KEYS[27], KEYS[28])
			redis.call("RPUSH", KEYS[29], tid)
			redis.call("SADD", KEYS[30], tid)
			`
		ka := []interface{}{
			//to-do model
			TId, td.Id,
			TDesc, td.Desc,
			TRemark, td.Remark,
			TCreateTime, time.Now().Unix(),
			TStartTime, td.StartTime,
			TPlace, td.Place,
			TRepeat, td.Repeat,
			TAllDay, td.AllDay,
			TRepeatMode, td.RepeatMode,
			TDone, td.Done,
			TFinishTime, td.FinishTime,
			TOwnerId, td.OwnerId,
			TMissionId, td.MissionId,
			TPid, td.Pid,
			//redis list
			fmt.Sprintf(userTdList, td.OwnerId),
			fmt.Sprintf(userTdNotDoneSet, td.OwnerId),
		}
		script := redis.NewScript(len(ka), lua)
		_, err := script.Do(c, ka...)
		if err != nil {
			log.Error("Error create todo:", err)
		}
		c.Close()
	}()
}

func readOwner(tid int) (*User, error) {
	c := dbms.Pool.Get()
	defer c.Close()
	uid, err := redis.Int(c.Do("HGET", "todo:"+strconv.Itoa(tid), TOwnerId))
	if err != nil {
		return nil, err
	}
	return readUser(uid)
}

func readBelongedMission(tid int) (*Mission, error) {
	c := dbms.Pool.Get()
	defer c.Close()
	mid, err := redis.Int(c.Do("HGET", "todo:"+strconv.Itoa(tid), TMissionId))
	if err != nil {
		return nil, err
	}
	return readMission(mid)
}

func updateTodoStatus(uid, tid int) error {
	c := dbms.Pool.Get()
	defer c.Close()
	done := fmt.Sprintf(userTdDoneSet, uid)
	notDone := fmt.Sprintf(userTdNotDoneSet, uid)
	_, err := c.Do("SMOVE", notDone, done, tid)
	return err
}

func updateTodo(tid int, kvMap map[string]interface{}) error {
	c := dbms.Pool.Get()
	defer c.Close()
	for k, v := range kvMap {
		c.Send("HSET", "todo:"+strconv.Itoa(tid), k, v)
	}
	return c.Flush()
}

func deleteTodo(tid int) error {
	c := dbms.Pool.Get()
	defer c.Close()
	_, err := redis.Bool(c.Do("DEL", "todo:"+strconv.Itoa(tid)))
	return err
}

//redis actions of model mission
func createMission(m *Mission) {
	c := dbms.Pool.Get()
	m.Id, _ = redis.Int(c.Do("INCR", "autoIncrComment"))
	m.Pid = base.HashedMissionId(m.Id)
	go dbms.CreateMissionIndex(m.Id, m.Pid)
	go func() {
		lua := `
			local mid = KEYS[2]
			redis.call("HMSET", "mission:"..mid,
					   KEYS[1], mid, KEYS[3], KEYS[4], KEYS[5], KEYS[6],
					   KEYS[7], KEYS[8], KEYS[9], KEYS[10], KEYS[11], KEYS[12],
					   KEYS[13], KEYS[14], KEYS[15], KEYS[16])
			redis.call("SADD", KEYS[17], mid)
			redis.call("SADD", KEYS[18], mid)
			`
		ka := []interface{}{
			//mission models
			MId, m.Id,
			MPid, m.Pid,
			MName, m.Name,
			MCreateTime, time.Now().Unix(),
			MDesc, m.Desc,
			MPublisherId, m.PublisherId,
			MCompletionNum, m.CompletionNum,
			MCompletedTime, m.CompletedTime,
			//redis set
			fmt.Sprintf(userMsPublishedSet, m.PublisherId),
			fmt.Sprintf(userMsAcceptedSet, m.PublisherId),
		}
		script := redis.NewScript(len(ka), lua)
		_, err := script.Do(c, ka...)
		if err != nil {
			log.Error("Error create mission:", err)
		}
		c.Close()
	}()
}

func createMissionComment(cm *Comment) {
	c := dbms.Pool.Get()
	cm.Id, _ = redis.Int(c.Do("INCR", "autoIncrComment"))
	cm.Pid = base.HashedCommentId(cm.Id)
	go func() {
		mid := dbms.ReadMissionId(cm.MissionPid)
		lua := `
			local cmid = KEYS[2]
			redis.call("HMSET", "comment:"..cmid,
					   KEYS[1], cmid, KEYS[3], KEYS[4], KEYS[5], KEYS[6],
					   KEYS[7], KEYS[8], KEYS[9], KEYS[10])
			redis.call("RPUSH", KEYS[11], cmid)
			`
		ka := []interface{}{
			//comment models
			CId, cm.Id,
			CPid, cm.Pid,
			CWhen, cm.When,
			CCriticPid, cm.CriticPid,
			CCriticName, cm.CriticName,

			//redis list
			fmt.Sprintf(missionCommentsList, mid),
		}
		script := redis.NewScript(len(ka), lua)
		_, err := script.Do(c, ka...)
		if err != nil {
			log.Error("Error create comment:", err)
		}
		c.Close()
	}()
}

func readMission(mid int) (*Mission, error) {
	c := dbms.Pool.Get()
	defer c.Close()
	mission := "mission:" + strconv.Itoa(mid)
	ret, err := redis.Values(c.Do("HGETALL", mission))
	if err != nil {
		return nil, err
	}
	m := new(Mission)
	err = redis.ScanStruct(ret, m)
	return m, err
}

func readMissionComments(mid int) (cms []*Comment, err error) {
	c := dbms.Pool.Get()
	defer c.Close()
	key := fmt.Sprintf(missionCommentsList, mid)
	lua := `
		local data = redis.call("LRANGE", KEYS[1], 0, -1)
		local ret = {}
  		for idx=1,#data do
  			ret[idx] = redis.call("HGETALL","comment:"..data[idx])
  		end
  		return ret
	`
	script := redis.NewScript(1, lua)
	rets, err := redis.Values(script.Do(c, key))
	cms = make([]*Comment, len(rets))
	for _, v := range rets {
		cmt := new(Comment)
		err = redis.ScanStruct(v.([]interface{}), cmt)
		cms = append(cms, cmt)
	}
	return
}

func readMissionReceiversId(mid int) (ids []int, err error) {
	c := dbms.Pool.Get()
	defer c.Close()
	key := fmt.Sprintf(missionReceiversSet, mid)
	ids, err = redis.Ints(c.Do("SMEMBERS", key))
	return
}

func updateMission(mid int, kvMap map[string]interface{}) error {
	c := dbms.Pool.Get()
	defer c.Close()
	for k, v := range kvMap {
		c.Send("HSET", "mission:"+strconv.Itoa(mid), k, v)
	}
	return c.Flush()
}

//redis actions of model project
func createProject(p *Project) {
	c := dbms.Pool.Get()
	p.Id, _ = redis.Int(c.Do("INCR", "autoIncrProject"))
	p.Pid = base.HashedProjectId(p.Id)
	go dbms.CreateProjectIndex(p.Id, p.Pid)
	go func() {
		lua := `
			local pid = KEYS[2]
			redis.call("HMSET", "project:"..pid,
					   KEYS[1], pid, KEYS[3], KEYS[4], KEYS[5], KEYS[6],
					   KEYS[7], KEYS[8], KEYS[9], KEYS[10], KEYS[11], KEYS[12]
					   KEYS[13], KYES[14])
			redis.call("SADD", KEYS[15], pid)
			redis.call("SADD", KEYS[16], pid)
			redis.call("SET", KEYS[17], pid)
			`
		ka := []interface{}{
			//project models
			PId, p.Id,
			PPid, p.Pid,
			PCreateTime, time.Now().Unix(),
			PDesc, p.Desc,
			PCreatorId, p.CreatorId,
			PPrivate, p.Private,
			PName, p.Name,
			//redis set
			fmt.Sprintf(userPjJoinedSet, p.CreatorId),
			fmt.Sprintf(userPjCreatedSet, p.CreatorId),
		}
		script := redis.NewScript(len(ka), lua)
		_, err := script.Do(c, ka...)
		if err != nil {
			log.Error("Error create project:", err)
		}
		c.Close()
	}()
}

func readProjectMembers(pid int) (reply []*User, err error) {
	c := dbms.Pool.Get()
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
	c := dbms.Pool.Get()
	defer c.Close()
	project := "project:" + strconv.Itoa(pid)
	ret, err := redis.Values(c.Do("HGETALL", project))
	if err != nil {
		return
	}
	err = redis.ScanStruct(ret, p)
	return
}

func readProjectMembersId(pid int) (ids []int, err error) {
	c := dbms.Pool.Get()
	defer c.Close()
	key := fmt.Sprintf(projectMembersSet, pid)
	ids, err = redis.Ints(c.Do("SMEMBERS", key))
	return
}

func updateProjectMember(pid, uid, action int) (err error) {
	c := dbms.Pool.Get()
	defer c.Close()
	memSet := fmt.Sprintf(projectMembersSet, pid)
	if action > 0 {
		_, err = c.Do("SADD", memSet, uid)

	} else {
		_, err = c.Do("SREM", memSet, uid)
	}
	return
}

func updateProject(pid int, kvMap map[string]interface{}) error {
	c := dbms.Pool.Get()
	defer c.Close()
	for k, v := range kvMap {
		c.Send("HSET", "project:"+strconv.Itoa(pid), k, v)
	}
	return c.Flush()
}

//chat
func createChat(ct *Chat) int {
	c := dbms.Pool.Get()
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
		ChConvId, ct.ConversationId,
		ChType, ct.Type,
		ChTarget, ct.Target,
		ChMsg, ct.Msg,

		ChGroupName, ct.GroupName,
		ChFrom, ct.From,
	}
	script := redis.NewScript(len(ka), lua)
	id, err := redis.Int(script.Do(c, ka...))
	if err != nil {
		log.Error("Error create offline conversation:", err)
	}
	return id
}

func readChatMembers(ct *Chat) (ids []int) {
	//read mission members
	m := new(Mission)
	m.Pid = ct.GroupName
	ids = m.GetReceiversId()
	return
}

func readDeviceToken(uid int) (token string, err error) {
	c := dbms.Pool.Get()
	defer c.Close()
	key := fmt.Sprintf(deviceToken, uid)
	token, err = redis.String(c.Do("GET", key))
	return
}

func updateOfflineMsg(uid, convId int) {
	c := dbms.Pool.Get()
	defer c.Close()
	key := fmt.Sprintf(offlineMsgList, uid)
	c.Do("RPUSH", key, convId)
}

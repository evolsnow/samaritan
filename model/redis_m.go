package model

import (
	"fmt"
	"github.com/evolsnow/samaritan/base"
	"github.com/evolsnow/samaritan/base/log"
	"github.com/evolsnow/samaritan/conn"
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
	UPassword   = "passwd"
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

	//additional
	UserId    = "userId:"    //index for userId, UserId:john's pid return john's userId
	MissionId = "missionId:" //index for Mission Id, missionId:a's pid return mission a's Id
	ProjectId = "project:"   //project:a's name return a's id

)

//visible func
func ReadUserId(uPid string) (uid int) {
	key := UserId + uPid
	uid, _ = redis.Int(conn.Get(key))
	return
}

func ReadMissionId(mPid string) (mid int) {
	key := MissionId + mPid
	mid, _ = redis.Int(conn.Get(key))
	return
}

func ReadProjectId(pPid string) (pid int) {
	key := ProjectId + pPid
	pid, _ = redis.Int(conn.Get(key))
	return
}

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
					KEYS[23], KEYS[24], KEYS[25], KEYS[26], KEYS[27], KEYS[28],
					KEYS[29], KEYS[30])
			redis.call("SADD", KEYS[31], uid)
			redis.call("SET", KEYS[32], uid)
			`
		ka := []interface{}{
			//user model
			UId, u.Id,
			UPid, u.Pid,
			USamId, u.SamId,
			UCreateTime, u.createTime,
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
			UserId + u.Pid,
		}
		script := redis.NewScript(len(ka), lua)
		_, err := script.Do(c, ka...)
		c.Close()
		if err != nil {
			log.Error("Error create user:", err)
		}
	}()
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
					KEYS[21], KEYS[22], KEYS[23], KEYS[24], KEYS[25], KEYS[26])
			redis.call("RPUSH", KEYS[27], tid)
			redis.call("SADD", KEYS[28], tid)
			`
		ka := []interface{}{
			//to-do model
			TId, td.Id,
			TDesc, td.Desc,
			TRemark, td.Remark,
			TCreateTime, td.createTime,
			TStartTime, td.StartTime,
			TPlace, td.Place,
			TRepeat, td.Repeat,
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

func updateTodoStatus(uid, tid int) error {
	c := conn.Pool.Get()
	defer c.Close()
	done := fmt.Sprintf(userTdDoneSet, uid)
	notDone := fmt.Sprintf(userTdNotDoneSet, uid)
	_, err := c.Do("SMOVE", notDone, done, tid)
	return err
}

//redis actions of model mission
func createMission(m *Mission) {
	c := conn.Pool.Get()
	m.Id, _ = redis.Int(c.Do("INCR", "autoIncrComment"))
	m.Pid = base.HashedMissionId(m.Id)
	go func() {
		lua := `
			local mid = KEYS[2]
			redis.call("HMSET", "mission:"..mid,
					   KEYS[1], mid, KEYS[3], KEYS[4], KEYS[5], KEYS[6],
					   KEYS[7], KEYS[8], KEYS[9], KEYS[10], KEYS[11], KEYS[12],
					   KEYS[13], KEYS[14], KEYS[15], KEYS[16])
			redis.call("SADD", KEYS[17], mid)
			redis.call("SADD", KEYS[18], mid)
			redis.call("SET", KEYS[19], mid)
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
			MissionId + m.Pid,
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
	c := conn.Pool.Get()
	cm.Id, _ = redis.Int(c.Do("INCR", "autoIncrComment"))
	cm.Pid = base.HashedCommentId(cm.Id)
	go func() {
		mid := ReadMissionId(cm.missionPid)
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

func readMissionComments(mid int) (cms []*Comment, err error) {
	c := conn.Pool.Get()
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
	c := conn.Pool.Get()
	defer c.Close()
	key := fmt.Sprintf(missionReceiversSet, mid)
	ids, err = redis.Ints(c.Do("SMEMBERS", key))
	return
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
			MissionId + p.Name,
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

func readProjectMembersId(pid int) (ids []int, err error) {
	c := conn.Pool.Get()
	defer c.Close()
	key := fmt.Sprintf(projectMembersSet, pid)
	ids, err = redis.Ints(c.Do("SMEMBERS", key))
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

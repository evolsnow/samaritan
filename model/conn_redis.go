package model

import (
	"encoding/json"
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
	MDeadline      = "deadline"
	MProjectId     = "projectId"
	//comments
	CId         = "id"
	CPid        = "pid"
	CCreateTime = "createTime"
	CCriticPid  = "criticPid"
	CCriticName = "criticName"
)

//project redis key name
const (
	PId            = "id"
	PPid           = "pid"
	PName          = "name"
	PCreateTime    = "createTime"
	PDesc          = "desc"
	PBackgroundPic = "backgroundPic"
	PCreatorId     = "creatorId"
	PPrivate       = "private"
)

const (
	ChId        = "id"
	ChPid       = "pid"
	ChType      = "type"
	ChTarget    = "target"
	ChMsg       = "msg"
	ChGroupName = "groupName"
	ChFrom      = "from"
	ChTimeStamp = "timestamp"
	ChInfo      = "info"
	ChDealt     = "dealt"
)

//other useful index set key name
const (
	//user
	userGroup          = "%s:%s:%d:%s"          //just for further analysis-> school:department:grade:class
	userTdList         = "user:%d:todoList"     //user's all to-do, redis-type:List
	userTdNotDoneSet   = "user:%d:todoStatus:0" //to-do status, redis-type:Set
	userTdDoneSet      = "user:%d:todoStatus:1"
	userPjJoinedSet    = "user:%d:projects:join" //user's all projects redis-type:Set
	userPjCreatedSet   = "user:%d:projects:create"
	userMsAcceptedSet  = "user:%d:missions:accept" //user's all missions redis-type:Set
	userMsPublishedSet = "user:%d:missions:publish"
	userMsCompletedSet = "user:%d:missions:complete"
	userCsDealtSet     = "user:%d:chats:complete"

	userPjColorList = "user:%d:project:%d:color" //user defined project color redis-type:List

	//to-do
	todoPictureList = "todo:%d:pictures" //to-do's pictures redis-type:List

	//mission
	missionReceiversSet     = "mission:%d:receivers"      //mission's receivers redis-type:Set
	missionCommentsList     = "mission:%d:comments"       //mission's comments redis-type:List
	missionCompletedUserSet = "mission:%d:completedUsers" //mission's completed user redis-type:Set
	missionPictureList      = "mission:%d:pictures"       //mission's pictures redis-type:List

	//project
	projectMembersSet   = "project:%d:members" //project's members redis-type:Set
	projectMissionsList = "project:%d:missions"
	//chat
	deviceToken    = "deviceToken:%d"     //ios device token
	offlineMsgList = "user:%d:offlineMsg" //redis type:list

)

//redis actions of model User

//create user in redis
func createUser(u *User) {
	c := dbms.Pool.Get()
	defer c.Close()
	u.Id, _ = redis.Int(c.Do("INCR", "autoIncrUser"))
	u.Pid = base.HashedUserId(u.Id)
	//index
	go dbms.CreateUserIndex(u.Id, u.Pid)
	//return user's id asap
	//go func() {
	lua := `
			local uid = KEYS[2]
			redis.call("HMSET", "user:"..uid,
					KEYS[1], KEYS[2], KEYS[3], KEYS[4], KEYS[5], KEYS[6], KEYS[7], KEYS[8],
					KEYS[9], KEYS[10], KEYS[11], KEYS[12], KEYS[13], KEYS[14], KEYS[15], KEYS[16],
					KEYS[17], KEYS[18], KEYS[19], KEYS[20], KEYS[21], KEYS[22],
					KEYS[23], KEYS[24], KEYS[25], KEYS[26], KEYS[27], KEYS[28],
					KEYS[29], KEYS[30])
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
	}
	script := redis.NewScript(len(ka), lua)
	_, err := script.Do(c, ka...)
	if err != nil {
		log.Error("Error create user:", err)
	}
	//redis set
	//todo not now (::0:)?
	c.Do("SADD", fmt.Sprintf(userGroup, u.School, u.Department, u.Grade, u.Class), u.Id)
	//}()
}

//return model User
func readUserWithId(id int) (*User, error) {
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

//save user avatar path
func createUserAvatar(uid int, avatarUrl string) error {
	c := dbms.Pool.Get()
	defer c.Close()
	user := "user:" + strconv.Itoa(uid)
	_, err := c.Do("HSET", user, UAvatar, avatarUrl)
	return err
}

//read user created projects
func readCreatedProjects(uid int) ([]Project, error) {
	key := fmt.Sprintf(userPjCreatedSet, uid)
	return readProjects(key)
}

//read user joined projects
func readJoinedProjects(uid int) ([]Project, error) {
	key := fmt.Sprintf(userPjJoinedSet, uid)
	return readProjects(key)
}

//read projects with the given key
func readProjects(key string) ([]Project, error) {
	c := dbms.Pool.Get()
	defer c.Close()
	lua := `
		local data = redis.call("SMEMBERS", KEYS[1])
		local ret = {}
  		for idx=1,#data do
  			local info = redis.call("HGETALL","project:"..data[idx])
			if (info == false) then
				redis.call("SREM", KEYS[1], data[idx])
			else
				ret[idx] = info
  			end
  		end
  		return ret
	`
	script := redis.NewScript(1, lua)
	results, err := redis.Values(script.Do(c, key))
	ps := make([]Project, len(results))
	for i, _ := range results {
		p := new(Project)
		err = redis.ScanStruct(results[i].([]interface{}), p)
		log.DebugJson(*p)
		ps[i] = *p
	}
	return ps, err
}

//read user's all completed missions' id
func readCompletedMissionsId(uid int) ([]int, error) {
	c := dbms.Pool.Get()
	defer c.Close()
	key := fmt.Sprintf(userMsCompletedSet, uid)
	ids, err := redis.Ints(c.Do("SMEMBERS", key))
	return ids, err

}

//read user's accepted missions's id
func readAcceptedMissionsId(uid int) ([]int, error) {
	c := dbms.Pool.Get()
	defer c.Close()
	ids, err := redis.Ints(c.Do("SMEMBERS", fmt.Sprintf(userMsAcceptedSet, uid)))
	return ids, err
}

//read user's dealt chat's id
func readDealtChatsId(uid int) ([]int, error) {
	c := dbms.Pool.Get()
	defer c.Close()
	ids, err := redis.Ints(c.Do("SMEMBERS", fmt.Sprintf(userCsDealtSet, uid)))
	return ids, err
}

//read user's password
func readPassword(uid int) (pwd string, err error) {
	c := dbms.Pool.Get()
	defer c.Close()
	user := "user:" + strconv.Itoa(uid)
	pwd, err = redis.String(c.Do("HGET", user, UPassword))
	return
}

//read user's name
func readName(uid int) (name string, err error) {
	c := dbms.Pool.Get()
	defer c.Close()
	user := "user:" + strconv.Itoa(uid)
	name, err = redis.String(c.Do("HGET", user, UName))
	return
}

//read offline msg
func readOfflineMsg(uid int) (reply []*Chat, err error) {
	c := dbms.Pool.Get()
	defer c.Close()
	key := fmt.Sprintf(offlineMsgList, uid)
	lua := `
	local data = redis.call("LRANGE", KEYS[1], "0", "-1")
	local ret = {}
  	for idx=1,#data do
  	  	local info = redis.call("HGETALL","chat:"..data[idx])
		if (info == false) then
			redis.call("LREM", KEYS[1], 1, data[idx])
		else
			ret[idx] = info
		end
  	end
  	return ret
   	`
	script := redis.NewScript(1, lua)
	chs, err := redis.Values(script.Do(c, key))
	reply = make([]*Chat, len(chs))
	for i, v := range chs {
		ch := new(Chat)
		err = redis.ScanStruct(v.([]interface{}), ch)
		json.Unmarshal([]byte(ch.SerializedInfo), &ch.ExtraInfo)
		reply[i] = ch
	}
	return reply, err
}

//add to user's completed mission
func updateCompletedMission(uid, mid int) error {
	c := dbms.Pool.Get()
	defer c.Close()
	c.Send("SADD", fmt.Sprintf(userMsCompletedSet, uid), mid)
	c.Send("SADD", fmt.Sprintf(missionCompletedUserSet, mid), uid)
	return c.Flush()
}

//add to user's uncompleted mission
func updateUnCompletedMission(uid, mid int) error {
	c := dbms.Pool.Get()
	defer c.Close()
	c.Send("SREM", fmt.Sprintf(userMsCompletedSet, uid), mid)
	c.Send("SREM", fmt.Sprintf(missionCompletedUserSet, mid), uid)
	return c.Flush()
}

//add to user's accepted mission
func updateAcceptedMission(uid, mid int) error {
	c := dbms.Pool.Get()
	defer c.Close()
	_, err := c.Do("SADD", fmt.Sprintf(userMsAcceptedSet, uid), mid)
	return err
}

//add to user's joined project
func updateJoinedProject(uid, pid int) error {
	c := dbms.Pool.Get()
	defer c.Close()
	_, err := c.Do("SADD", fmt.Sprintf(userPjJoinedSet, uid), pid)
	return err
}

//add to user's completed mission
func updateDealtChat(uid, cid int, dealt bool) (err error) {
	c := dbms.Pool.Get()
	defer c.Close()
	if dealt {
		_, err = c.Do("SADD", fmt.Sprintf(userCsDealtSet, uid), cid)
	} else {
		_, err = c.Do("SREM", fmt.Sprintf(userCsDealtSet, uid), cid)
	}
	return
}

//update user with given value
func updateUser(uid int, kvMap map[string]interface{}) error {
	c := dbms.Pool.Get()
	defer c.Close()
	for k, v := range kvMap {
		c.Send("HSET", "user:"+strconv.Itoa(uid), k, v)
	}
	return c.Flush()
}

//redis actions of model to-do
//create to-do in redis
func createTodo(td *Todo) {
	c := dbms.Pool.Get()
	defer c.Close()
	td.Id, _ = redis.Int(c.Do("INCR", "autoIncrTodo"))
	td.Pid = base.HashedTodoId(td.Id)
	go dbms.CreateTodoIndex(td.Id, td.Pid)

	//go func() {
	lua := `
			local tid = KEYS[2]
			redis.call("HMSET", "todo:"..tid,
					KEYS[1], tid, KEYS[3], KEYS[4], KEYS[5], KEYS[6], KEYS[7], KEYS[8],
					KEYS[9], KEYS[10], KEYS[11], KEYS[12], KEYS[13], KEYS[14],
					KEYS[15], KEYS[16], KEYS[17], KEYS[18], KEYS[19], KEYS[20],
					KEYS[21], KEYS[22], KEYS[23], KEYS[24], KEYS[25], KEYS[26],
					KEYS[27], KEYS[28])

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
	}
	script := redis.NewScript(len(ka), lua)
	_, err := script.Do(c, ka...)
	if err != nil {
		log.Error("Error create todo:", err)
	}

	c.Send("RPUSH", fmt.Sprintf(userTdList, td.OwnerId), td.Id)
	c.Send("SADD", fmt.Sprintf(userTdNotDoneSet, td.OwnerId), td.Id)
	c.Flush()
	//}()
}

//return model To-do
func readTodoWithId(id int) (*Todo, error) {
	c := dbms.Pool.Get()
	defer c.Close()
	todo := "todo:" + strconv.Itoa(id)
	ret, err := redis.Values(c.Do("HGETALL", todo))
	if err != nil {
		return nil, err
	}
	t := new(Todo)
	err = redis.ScanStruct(ret, t)
	return t, err
}

//get to-do's owner
func readOwner(tid int) (*User, error) {
	c := dbms.Pool.Get()
	defer c.Close()
	uid, err := redis.Int(c.Do("HGET", "todo:"+strconv.Itoa(tid), TOwnerId))
	if err != nil {
		return nil, err
	}
	return readUserWithId(uid)
}

//get to-do's belonged mission
func readBelongedMission(tid int) (*Mission, error) {
	c := dbms.Pool.Get()
	defer c.Close()
	mid, err := redis.Int(c.Do("HGET", "todo:"+strconv.Itoa(tid), TMissionId))
	if err != nil {
		return nil, err
	}
	return readMissionWithId(mid)
}

//read to-do pics
func readTodoPics(tid int) (pics []string, err error) {
	c := dbms.Pool.Get()
	defer c.Close()
	key := fmt.Sprintf(todoPictureList, tid)
	pics, err = redis.Strings(c.Do("LRANGE", key, 0, -1))
	return
}

//set to-do status to done
func updateTodoStatus(uid, tid int) error {
	c := dbms.Pool.Get()
	defer c.Close()
	done := fmt.Sprintf(userTdDoneSet, uid)
	notDone := fmt.Sprintf(userTdNotDoneSet, uid)
	_, err := c.Do("SMOVE", notDone, done, tid)
	return err
}

//update to-do pics
func updateTodoPics(tid int, pics []string) error {
	c := dbms.Pool.Get()
	defer c.Close()
	key := fmt.Sprintf(todoPictureList, tid)
	c.Send("DEL", key)
	for _, v := range pics {
		c.Send("RPUSH", key, v)

	}
	return c.Flush()
}

//update to-do with given value
func updateTodo(tid int, kvMap map[string]interface{}) error {
	c := dbms.Pool.Get()
	defer c.Close()
	for k, v := range kvMap {
		c.Send("HSET", "todo:"+strconv.Itoa(tid), k, v)
	}
	if _, ok := kvMap[TDone]; ok {
		uid := kvMap[TOwnerId]
		go updateTodoStatus(uid.(int), tid)
	}
	return c.Flush()
}

//delete a to-do
func deleteTodo(tid int) error {
	c := dbms.Pool.Get()
	defer c.Close()
	c.Send("DEL", "todo:"+strconv.Itoa(tid))
	c.Send("DEL", fmt.Sprintf(todoPictureList), tid)
	//c.Send("LREM", fmt.Sprintf(userTdList, tid), 1, tid)
	//c.Send("SREM", fmt.Sprintf(userTdNotDoneSet, tid), tid)
	//c.Send("SREM", fmt.Sprintf(userTdDoneSet, tid), tid)
	return c.Flush()
}

//redis actions of model mission
//create a mission in redis
func createMission(m *Mission) {
	c := dbms.Pool.Get()
	defer c.Close()
	m.Id, _ = redis.Int(c.Do("INCR", "autoIncrComment"))
	m.Pid = base.HashedMissionId(m.Id)
	go dbms.CreateMissionIndex(m.Id, m.Pid)
	//go func() {
	lua := `
			local mid = KEYS[2]
			redis.call("HMSET", "mission:"..mid,
					   KEYS[1], mid, KEYS[3], KEYS[4], KEYS[5], KEYS[6],
					   KEYS[7], KEYS[8], KEYS[9], KEYS[10], KEYS[11], KEYS[12],
					   KEYS[13], KEYS[14], KEYS[15], KEYS[16], KEYS[17], KEYS[18],
					   KEYS[19], KEYS[20])

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
		MDeadline, m.Deadline,
		MProjectId, m.ProjectId,
	}
	script := redis.NewScript(len(ka), lua)
	_, err := script.Do(c, ka...)
	if err != nil {
		log.Error("Error create mission:", err)
	}

	c.Send("SADD", fmt.Sprintf(userMsPublishedSet, m.PublisherId), m.Id)
	c.Send("SADD", fmt.Sprintf(userMsAcceptedSet, m.PublisherId), m.Id)
	if m.ProjectId != 0 {
		c.Send("LPUSH", fmt.Sprintf(projectMissionsList, m.ProjectId), m.Id)
	}
	c.Flush()
	//}()
}

//create a mission comment
func createMissionComment(cm *Comment) {
	c := dbms.Pool.Get()
	defer c.Close()
	cm.Id, _ = redis.Int(c.Do("INCR", "autoIncrComment"))
	cm.Pid = base.HashedCommentId(cm.Id)
	//go func() {
	mid := dbms.ReadMissionId(cm.MissionPid)
	lua := `
			local cmid = KEYS[2]
			redis.call("HMSET", "comment:"..cmid,
					   KEYS[1], cmid, KEYS[3], KEYS[4], KEYS[5], KEYS[6],
					   KEYS[7], KEYS[8], KEYS[9], KEYS[10])
			`
	ka := []interface{}{
		//comment models
		CId, cm.Id,
		CPid, cm.Pid,
		CCreateTime, time.Now().Unix(),
		CCriticPid, cm.CriticPid,
		CCriticName, cm.CriticName,
	}
	script := redis.NewScript(len(ka), lua)
	_, err := script.Do(c, ka...)
	if err != nil {
		log.Error("Error create comment:", err)
	}
	c.Do("RPUSH", fmt.Sprintf(missionCommentsList, mid), cm.Id)
	//}()
}

//return model Mission
func readMissionWithId(mid int) (*Mission, error) {
	c := dbms.Pool.Get()
	defer c.Close()
	mission := "mission:" + strconv.Itoa(mid)
	ret, err := redis.Values(c.Do("HGETALL", mission))
	if err != nil {
		return nil, err
	}
	m := new(Mission)
	err = redis.ScanStruct(ret, m)
	//m.ReceiversId, _ = readMissionReceiversId(m.Id)
	return m, err
}

//get all mission's comments
func readMissionComments(mid int) (cms []Comment, err error) {
	c := dbms.Pool.Get()
	defer c.Close()
	key := fmt.Sprintf(missionCommentsList, mid)
	lua := `
		local data = redis.call("LRANGE", KEYS[1], 0, -1)
		local ret = {}
  		for idx=1,#data do
			local info = redis.call("HGETALL","comment:"..data[idx])
			if (info == false) then
				redis.call("LREM", KEYS[1], 1, data[idx])
			else
				ret[idx] = info
			end
  		end
  		return ret
	`
	script := redis.NewScript(1, lua)
	results, err := redis.Values(script.Do(c, key))
	cms = make([]Comment, len(results))
	for i, v := range results {
		cmt := new(Comment)
		err = redis.ScanStruct(v.([]interface{}), cmt)
		cms[i] = *cmt
	}
	return
}

//get all mission's receivers id
func readMissionReceiversId(mid int) (ids []int, err error) {
	c := dbms.Pool.Get()
	defer c.Close()
	key := fmt.Sprintf(missionReceiversSet, mid)
	ids, err = redis.Ints(c.Do("SMEMBERS", key))
	return
}

//get all mission's receivers id
func readMissionReceiversName(mid int) (names []string, err error) {
	c := dbms.Pool.Get()
	defer c.Close()
	key := fmt.Sprintf(missionReceiversSet, mid)
	lua := `
		local data = redis.call("SMEMBERS", KEYS[1])
		local ret = {}
  		for idx=1,#data do
  			local info = redis.call("HGET","user:"..data[idx], KEYS[2])
			if (info == false) then
				redis.call("SREM", KEYS[1], data[idx])
			else
				ret[idx] = info
			end
  		end
  		return ret
	`
	script := redis.NewScript(2, lua)
	names, err = redis.Strings(script.Do(c, key, UName))
	return
}

//get all mission receivers id who has completed it
func readMissionCompletedUsersId(mid int) (ids []int, err error) {
	c := dbms.Pool.Get()
	defer c.Close()
	key := fmt.Sprintf(missionCompletedUserSet, mid)
	ids, err = redis.Ints(c.Do("SMEMBERS", key))
	return
}

//get mission's pics
func readMissionPics(mid int) (pics []string, err error) {
	c := dbms.Pool.Get()
	defer c.Close()
	key := fmt.Sprintf(missionPictureList, mid)
	pics, err = redis.Strings(c.Do("LRANGE", key, 0, -1))
	return
}

//update mission pics
func updateMissionPics(mid int, pics []string) error {
	c := dbms.Pool.Get()
	defer c.Close()
	key := fmt.Sprintf(missionPictureList, mid)
	c.Send("DEL", key)
	for _, v := range pics {
		c.Send("RPUSH", key, v)
	}
	return c.Flush()
}

//update mission with given value
func updateMission(mid int, kvMap map[string]interface{}) error {
	c := dbms.Pool.Get()
	defer c.Close()
	for k, v := range kvMap {
		c.Send("HSET", "mission:"+strconv.Itoa(mid), k, v)
	}
	return c.Flush()
}

//add or remove mission receivers
func updateMissionReceiver(mid, uid, action int) (err error) {
	c := dbms.Pool.Get()
	defer c.Close()
	memSet := fmt.Sprintf(missionReceiversSet, mid)
	if action > 0 {
		_, err = c.Do("SADD", memSet, uid)
	} else {
		_, err = c.Do("SREM", memSet, uid)
	}
	return
}

//delete a mission
func deleteMission(mid int) error {
	c := dbms.Pool.Get()
	defer c.Close()
	c.Send("DEL", "mission:"+strconv.Itoa(mid))
	//c.Send("SREM",fmt.Sprintf(userMsAcceptedSet,mid),mid)
	//c.Send("SREM",fmt.Sprintf(userMsPublishedSet,mid),mid)
	//c.Send("SREM",fmt.Sprintf(userMsCompletedSet,mid),mid)

	//c.Send("DEL",fmt.Sprintf(missionReceiversSet,mid))
	//c.Send("DEL",fmt.Sprintf(missionCommentsList,mid))
	//c.Send("DEL",fmt.Sprintf(missionCompletedUserSet,mid))
	c.Send("DEL", fmt.Sprintf(missionPictureList, mid))
	return c.Flush()
}

//redis actions of model project
//create project in redis
func createProject(p *Project) {
	c := dbms.Pool.Get()
	defer c.Close()
	p.Id, _ = redis.Int(c.Do("INCR", "autoIncrProject"))
	p.Pid = base.HashedProjectId(p.Id)
	go dbms.CreateProjectIndex(p.Id, p.Pid)
	//go func() {
	lua := `
			local pid = KEYS[2]
			redis.call("HMSET", "project:"..pid,
					   KEYS[1], pid, KEYS[3], KEYS[4], KEYS[5], KEYS[6],
					   KEYS[7], KEYS[8], KEYS[9], KEYS[10], KEYS[11], KEYS[12],
					   KEYS[13], KEYS[14], KEYS[15], KEYS[16])
			`
	ka := []interface{}{
		//project models
		PId, p.Id,
		PPid, p.Pid,
		PCreateTime, time.Now().Unix(),
		PDesc, p.Desc,
		PBackgroundPic, p.BackgroundPic,
		PCreatorId, p.CreatorId,
		PPrivate, p.Private,
		PName, p.Name,
	}
	script := redis.NewScript(len(ka), lua)
	_, err := script.Do(c, ka...)
	if err != nil {
		log.Error("Error create project:", err)
	}
	c.Send("SADD", fmt.Sprintf(userPjJoinedSet, p.CreatorId), p.Id)
	c.Send("SADD", fmt.Sprintf(userPjJoinedSet, p.CreatorId), p.Id)
	c.Send("SADD", fmt.Sprintf(projectMembersSet, p.Id), p.CreatorId)
	c.Flush()
	//}()
}

//get project's members
func readProjectMembers(pid int) (reply []*User, err error) {
	c := dbms.Pool.Get()
	defer c.Close()
	rcvSet := fmt.Sprintf(projectMembersSet, pid)
	lua := `
	local data = redis.call("SMEMBERS", KEYS[1])
	local ret = {}
  	for idx=1,#data do
  	  	local info = redis.call("HGETALL","user:"..data[idx])
		if (info == false) then
			redis.call("SREM", KEYS[1], data[idx])
		else
			ret[idx] = info
		end
  	end
  	return ret
   	`
	script := redis.NewScript(1, lua)
	users, err := redis.Values(script.Do(c, rcvSet))
	reply = make([]*User, len(users))
	for i, v := range users {
		rcv := new(User)
		err = redis.ScanStruct(v.([]interface{}), rcv)
		reply[i] = rcv
	}
	return reply, err
}

//get project's missions
func readProjectMissions(pid int) (reply []*Mission, err error) {
	c := dbms.Pool.Get()
	defer c.Close()
	msList := fmt.Sprintf(projectMissionsList, pid)
	lua := `
	local data = redis.call("LRANGE", KEYS[1], "0", "-1")
	local ret = {}
  	for idx=1,#data do
  	  	local info = redis.call("HGETALL","mission:"..data[idx])
		if (info == false) then
			redis.call("SREM", KEYS[1], data[idx])
		else
			ret[idx] = info
		end
  	end
  	return ret
   	`
	script := redis.NewScript(1, lua)
	ms, err := redis.Values(script.Do(c, msList))
	reply = make([]*Mission, len(ms))
	for i, v := range ms {
		m := new(Mission)
		err = redis.ScanStruct(v.([]interface{}), m)
		reply[i] = m
	}
	return reply, err
}

//return model Project
func readProjectWithId(pid int) (*Project, error) {
	c := dbms.Pool.Get()
	defer c.Close()
	project := "project:" + strconv.Itoa(pid)
	ret, err := redis.Values(c.Do("HGETALL", project))
	if err != nil {
		return nil, err
	}
	p := new(Project)
	err = redis.ScanStruct(ret, p)
	return p, err
}

//get project's creator
func readCreator(pid int) (*User, error) {
	c := dbms.Pool.Get()
	defer c.Close()
	uid, err := redis.Int(c.Do("HGET", "project:"+strconv.Itoa(pid), PCreatorId))
	if err != nil {
		return nil, err
	}
	return readUserWithId(uid)
}

//get projects's members' id
func readProjectMembersId(pid int) (ids []int, err error) {
	c := dbms.Pool.Get()
	defer c.Close()
	key := fmt.Sprintf(projectMembersSet, pid)
	ids, err = redis.Ints(c.Do("SMEMBERS", key))
	return
}

//get projects's members' name
func readProjectMembersName(pid int) (names []string, err error) {
	c := dbms.Pool.Get()
	defer c.Close()
	pSet := fmt.Sprintf(projectMembersSet, pid)
	lua := `
	local data = redis.call("SMEMBERS", KEYS[1])
	local ret = {}
  	for idx=1,#data do
  	  	local info = redis.call("HGET", "user:"..data[idx], KEYS[2])
		if (info == false) then
			redis.call("SREM", KEYS[1], data[idx])
		else
			ret[idx] = info
		end
  	end
  	return ret
   	`
	script := redis.NewScript(2, lua)
	names, err = redis.Strings(script.Do(c, pSet, UName))
	return
}

//add or remove project member
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

//update project with given value
func updateProject(pid int, kvMap map[string]interface{}) error {
	c := dbms.Pool.Get()
	defer c.Close()
	for k, v := range kvMap {
		c.Send("HSET", "project:"+strconv.Itoa(pid), k, v)
	}
	return c.Flush()
}

//delete a project
func deleteProject(pid int) error {
	c := dbms.Pool.Get()
	defer c.Close()
	c.Send("DEL", "project:"+strconv.Itoa(pid))
	c.Send("DEL", fmt.Sprintf(projectMembersSet, pid))
	c.Send("DEL", fmt.Sprintf(projectMissionsList, pid))
	return c.Flush()
}

//delete project from user project set
func deleteFromUserProjectSet(pid, uid int) error {
	c := dbms.Pool.Get()
	defer c.Close()
	c.Send("SREM", fmt.Sprintf(userPjJoinedSet, uid), pid)
	c.Send("SREM", fmt.Sprintf(userPjCreatedSet, uid), pid)
	return c.Flush()
}

//chat
//create chat in redis
func createChat(ct *Chat) int {
	c := dbms.Pool.Get()
	defer c.Close()
	ct.Id, _ = redis.Int(c.Do("INCR", "autoIncrChat"))
	ct.Pid = base.HashedChatId(ct.Id)
	go dbms.CreateChatIndex(ct.Id, ct.Pid)
	//todo expire the msg
	lua := `
	local cid = KEYS[2]
	redis.call("HMSET", "chat:"..cid,
					KEYS[1], cid, KEYS[3], KEYS[4], KEYS[5], KEYS[6],
					KEYS[7], KEYS[8], KEYS[9], KEYS[10], KEYS[11], KEYS[12],
					KEYS[13], KEYS[14], KEYS[15], KEYS[16], KEYS[17], KEYS[18],
					KEYS[19], KEYS[20])
	return cid
	`
	ka := []interface{}{
		ChId, ct.Id,
		ChPid, ct.Pid,
		ChType, ct.Type,
		ChTarget, ct.Target,
		ChMsg, ct.Msg,
		ChInfo, ct.SerializedInfo,
		ChTimeStamp, ct.Timestamp,
		ChGroupName, ct.GroupName,
		ChFrom, ct.From,
		ChDealt, ct.Dealt,
	}
	script := redis.NewScript(len(ka), lua)
	id, err := redis.Int(script.Do(c, ka...))
	if err != nil {
		log.Error("Error create offline conversation:", err)
	}
	return id
}

//get all chat members
func readChatMembers(ct *Chat) (ids []int) {
	//read mission members todo:mission or project
	m := new(Mission)
	m.Pid = ct.GroupName
	ids = m.GetReceiversId()
	return
}

//create offline msg when failed to push
func createOfflineMsg(uid, convId int) {
	c := dbms.Pool.Get()
	defer c.Close()
	key := fmt.Sprintf(offlineMsgList, uid)
	c.Do("LPUSH", key, convId)
}

func readUserMsgsId(uid int) (ids []int, err error) {
	c := dbms.Pool.Get()
	defer c.Close()
	ids, err = redis.Ints(c.Do("LRANGE", fmt.Sprintf(offlineMsgList, uid), "0", "-1"))
	return
}

//return model Chat
func readChatWithId(cid int) (*Chat, error) {
	c := dbms.Pool.Get()
	defer c.Close()
	chat := "chat:" + strconv.Itoa(cid)
	ret, err := redis.Values(c.Do("HGETALL", chat))
	if err != nil {
		return nil, err
	}
	ch := new(Chat)
	err = redis.ScanStruct(ret, ch)
	json.Unmarshal([]byte(ch.SerializedInfo), ch.ExtraInfo)
	//m.ReceiversId, _ = readMissionReceiversId(m.Id)
	return ch, err
}

package dbms

import (
	"github.com/garyburd/redigo/redis"
	"time"
)

const (
	samIdSet = "allSamId"

	SearchPhoneIndex = "index:search:phone"
	SearchMailIndex  = "index:search:mail"
	SearchSamIndex   = "index:search:sam"
)

const (
	//index
	UserIndex    = "index:user"
	TodoIndex    = "index:todo"
	MissionIndex = "index:mission"
	ProjectIndex = "index:project"
)

// All extra redis actions

func Get(key string) (string, error) {
	c := Pool.Get()
	defer c.Close()
	value, err := redis.String(c.Do("GET", key))
	return value, err
}

func CacheGet(key string) string {
	c := CachePool.Get()
	defer c.Close()
	value, _ := redis.String(c.Do("GET", key))
	return value
}

func CacheGetSet(key, newValue string) string {
	c := CachePool.Get()
	defer c.Close()
	value, _ := redis.String(c.Do("GETSET", key, newValue))
	return value
}

func CacheSet(key string, value interface{}, px time.Duration) {
	c := CachePool.Get()
	defer c.Close()
	c.Do("SET", key, value, "PX", int64(px/time.Millisecond))
}

func CacheDelete(key string) {
	c := CachePool.Get()
	defer c.Close()
	c.Do("DEL", key)
}

func CreateSearchIndex(uid int, info, searchType string) {
	c := Pool.Get()
	defer c.Close()
	switch searchType {
	case "phone":
		c.Do("HSET", SearchPhoneIndex, info, uid)
	case "mail":
		c.Do("HSET", SearchMailIndex, info, uid)
	case "samId":
		c.Do("HSET", SearchSamIndex, info, uid)
	}
}

func ReadUserIdWithIndex(info, loginType string) (uid int) {
	c := Pool.Get()
	defer c.Close()
	switch loginType {
	case "phone":
		uid, _ = redis.Int(c.Do("HGET", SearchPhoneIndex, info))
	case "mail":
		uid, _ = redis.Int(c.Do("HGET", SearchMailIndex, info))
	case "samId":
		uid, _ = redis.Int(c.Do("HGET", SearchSamIndex, info))
	}
	return
}

func ReadIfSamIdExist(sid string) (exist bool) {
	c := Pool.Get()
	defer c.Close()
	exist, _ = redis.Bool(c.Do("SISMEMBER", samIdSet, sid))
	return
}

func UpdateSamIdSet(sid string) {
	c := Pool.Get()
	defer c.Close()
	c.Do("SADD", samIdSet, sid)
}

func DeleteSamId(sid string) {
	c := Pool.Get()
	defer c.Close()
	c.Do("SREM", samIdSet, sid)
}

//create index
func CreateUserIndex(uid int, uPid string) {
	c := Pool.Get()
	defer c.Close()
	c.Do("HSET", UserIndex, uPid, uid)
}

func CreateTodoIndex(tid int, tPid string) {
	c := Pool.Get()
	defer c.Close()
	c.Do("HSET", TodoIndex, tPid, tid)
}

func CreateMissionIndex(mid int, mPid string) {
	c := Pool.Get()
	defer c.Close()
	c.Do("HSET", MissionIndex, mPid, mid)
}

func CreateProjectIndex(pid int, pPid string) {
	c := Pool.Get()
	defer c.Close()
	c.Do("HSET", ProjectIndex, pPid, pid)
}

//get real id from public id
func ReadUserId(uPid string) (uid int) {
	c := Pool.Get()
	defer c.Close()
	uid, _ = redis.Int(c.Do("HGET", UserIndex, uPid))
	return
}

func ReadMissionId(mPid string) (mid int) {
	c := Pool.Get()
	defer c.Close()
	mid, _ = redis.Int(c.Do("HGET", MissionIndex, mPid))
	return
}

func ReadProjectId(pPid string) (pid int) {
	c := Pool.Get()
	defer c.Close()
	pid, _ = redis.Int(c.Do("HGET", ProjectIndex, pPid))
	return
}

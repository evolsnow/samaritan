/*
Package dbms provides common redis action
*/
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
	DeviceIndex  = "index:device"
)

// All extra redis actions

// Get reads a key from redis db
func Get(key string) (string, error) {
	c := Pool.Get()
	defer c.Close()
	value, err := redis.String(c.Do("GET", key))
	return value, err
}

// CacheGet reads a key from redis cache db
func CacheGet(key string) string {
	c := CachePool.Get()
	defer c.Close()
	value, _ := redis.String(c.Do("GET", key))
	return value
}

// CacheGetSet gets a key from the cache db and sets with a new value
func CacheGetSet(key, newValue string) string {
	c := CachePool.Get()
	defer c.Close()
	value, _ := redis.String(c.Do("GETSET", key, newValue))
	return value
}

// CacheSet sets value with expire
func CacheSet(key string, value interface{}, px time.Duration) {
	c := CachePool.Get()
	defer c.Close()
	c.Do("SET", key, value, "PX", int64(px/time.Millisecond))
}

// CacheDelete delete the key from cache db
func CacheDelete(key string) {
	c := CachePool.Get()
	defer c.Close()
	c.Do("DEL", key)
}

// CreateSearchIndex creates index from user id, phone/mail/sam id
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

// ReadUserIdWithIndex reads user id from redis db with index type
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

// ReadIfSamIdExist tests if the given sam id is in redis db set
func ReadIfSamIdExist(sid string) (exist bool) {
	c := Pool.Get()
	defer c.Close()
	exist, _ = redis.Bool(c.Do("SISMEMBER", samIdSet, sid))
	return
}

// UpdateSamIdSet adds a sam id to the set
func UpdateSamIdSet(sid string) {
	c := Pool.Get()
	defer c.Close()
	c.Do("SADD", samIdSet, sid)
}

// DeleteSamIdSet deletes a sam id from the set
func DeleteSamId(sid string) {
	c := Pool.Get()
	defer c.Close()
	c.Do("SREM", samIdSet, sid)
}

// CreateUserIndex creates 'user public id==>> user real id' index
func CreateUserIndex(uid int, uPid string) {
	c := Pool.Get()
	defer c.Close()
	c.Do("HSET", UserIndex, uPid, uid)
}

// CreateTodoIndex creates 'to-do public id==>> to-do real id' index
func CreateTodoIndex(tid int, tPid string) {
	c := Pool.Get()
	defer c.Close()
	c.Do("HSET", TodoIndex, tPid, tid)
}

// CreateMissionIndex creates 'mission public id==>> mission real id' index
func CreateMissionIndex(mid int, mPid string) {
	c := Pool.Get()
	defer c.Close()
	c.Do("HSET", MissionIndex, mPid, mid)
}

// CreateProjectIndex creates 'project public id==>> project real id' index
func CreateProjectIndex(pid int, pPid string) {
	c := Pool.Get()
	defer c.Close()
	c.Do("HSET", ProjectIndex, pPid, pid)
}

// CreateDeviceIndex creates 'user real id==>> user device token' index
func CreateDeviceIndex(uid int, dt string) {
	c := Pool.Get()
	defer c.Close()
	c.Do("HSET", DeviceIndex, uid, dt)
}

// ReadUserId gets user real id with public id
func ReadUserId(uPid string) (uid int) {
	c := Pool.Get()
	defer c.Close()
	uid, _ = redis.Int(c.Do("HGET", UserIndex, uPid))
	return
}

// ReadTodoId gets to-do real id with public id
func ReadTodoId(tPid string) (tid int) {
	c := Pool.Get()
	defer c.Close()
	tid, _ = redis.Int(c.Do("HGET", TodoIndex, tPid))
	return
}

// ReadMissionId gets mission real id with public id
func ReadMissionId(mPid string) (mid int) {
	c := Pool.Get()
	defer c.Close()
	mid, _ = redis.Int(c.Do("HGET", MissionIndex, mPid))
	return
}

// ReadProjectId gets project real id with public id
func ReadProjectId(pPid string) (pid int) {
	c := Pool.Get()
	defer c.Close()
	pid, _ = redis.Int(c.Do("HGET", ProjectIndex, pPid))
	return
}

// ReadDeviceToken gets user device token with real id
func ReadDeviceToken(uid int) (dt string) {
	c := Pool.Get()
	defer c.Close()
	dt, _ = redis.String(c.Do("HGET", DeviceIndex, uid))
	return
}

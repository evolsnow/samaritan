package dbms

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"time"
)

const (
	userToken = "user:%d:token"
	samIdSet  = "allSamId"
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

func CreateToken(uid int, token string) {
	c := Pool.Get()
	defer c.Close()
	key := fmt.Sprintf(userToken, uid)
	c.Do("SET", key, token)
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

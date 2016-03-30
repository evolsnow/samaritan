package dbms

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"time"
)

const (
	userToken = "user:%d:token"
	samIdSet  = "allSamId"

	LoginPhoneIndex = "login:phone:%s"
	LoginMailIndex  = "login:mail:%s"
	LoginSamIndex   = "login:sam:%s"
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

func CreateLoginIndex(uid int, info, loginType string) {
	c := Pool.Get()
	defer c.Close()
	switch loginType {
	case "phone":
		c.Do("SET", fmt.Sprintf(LoginPhoneIndex, info), uid)
	case "mail":
		c.Do("SET", fmt.Sprintf(LoginMailIndex, info), uid)
	case "samId":
		c.Do("SET", fmt.Sprintf(LoginSamIndex, info), uid)
	}
}

func ReadLoginUid(info, loginType string) (uid int) {
	c := Pool.Get()
	defer c.Close()
	switch loginType {
	case "phone":
		uid, _ = redis.Int(c.Do("GET", fmt.Sprintf(LoginPhoneIndex, info)))
	case "mail":
		uid, _ = redis.Int(c.Do("GET", fmt.Sprintf(LoginMailIndex, info)))
	case "samId":
		uid, _ = redis.Int(c.Do("GET", fmt.Sprintf(LoginSamIndex, info)))
	}
	return
}

func ReadToken(uid int) (token string) {
	c := Pool.Get()
	defer c.Close()
	key := fmt.Sprintf(userToken, uid)
	token, _ = redis.String(c.Do("GET", key))
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

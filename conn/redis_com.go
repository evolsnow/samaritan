package conn

import (
	"github.com/garyburd/redigo/redis"
	"time"
)

//All extra redis actions
//
//func GetSignKey(id string) (key, lastVisit string) {
//	c := Pool.Get()
//	defer c.Close()
//	user := "vsuser:" + id
//	ret, err := redis.Strings(c.Do("HMGET", user, "appKey", "lastVisit"))
//	if err != nil {
//		log.Println("no user %s", user)
//		return
//	}
//	return ret[0], ret[1]
//}
//
//func UpdateSign(id, lv string) {
//	c := Pool.Get()
//	defer c.Close()
//	user := "vsuser:" + id
//	c.Do("HSET", user, "lastVisit", lv)
//}

func Get(key string) string {
	c := Pool.Get()
	defer c.Close()
	value, _ := redis.String(c.Do("GET", key))
	return value
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

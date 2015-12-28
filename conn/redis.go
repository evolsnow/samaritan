package conn

import (
	"github.com/garyburd/redigo/redis"
	"log"
)

func GetSignKey(id string) (key, lastVisit string) {
	c := Pool.Get()
	defer c.Close()
	user := "vsuser:" + id
	ret, err := redis.Strings(c.Do("HMGET", user, "appKey", "lastVisit"))
	if err != nil {
		log.Println("no user %s", user)
		return
	}
	return ret[0], ret[1]
}

func UpdateSign(id, lv string) {
	c := Pool.Get()
	defer c.Close()
	user := "vsuser:" + id
	c.Do("HSET", user, "lastVisit", lv)
}

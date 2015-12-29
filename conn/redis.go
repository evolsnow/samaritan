package conn

import (
	"github.com/garyburd/redigo/redis"
	"log"
)

//All redis actions

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

func GetPassword(id string) string {
	c := Pool.Get()
	defer c.Close()
	user := "user:" + id
	pwd, _ := redis.String(c.Do("HGET", user, "passwd"))
	return pwd
}

func SetPassword(id, pwd string) {
	c := Pool.Get()
	defer c.Close()
	user := "user:" + id
	_, err := c.Do("HSET", user, "passwd", pwd)
	if err != nil {
		log.Println("Failed to save password for user:%s", id)
	}
}

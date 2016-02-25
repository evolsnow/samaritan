package model

import (
	"github.com/evolsnow/samaritan/conn"
	"github.com/garyburd/redigo/redis"
	"log"
	"strconv"
	"time"
)

//All model's redis actions

func getPassword(id int) string {
	c := conn.Pool.Get()
	defer c.Close()
	user := "user:" + strconv.Itoa(id)
	pwd, _ := redis.String(c.Do("HGET", user, "passwd"))
	return pwd
}

func updatePassword(id int, pwd string) {
	c := conn.Pool.Get()
	defer c.Close()
	user := "user:" + strconv.Itoa(id)
	_, err := c.Do("HSET", user, "passwd", pwd)
	if err != nil {
		log.Println("Failed to update password for user:%s", id)
	}
}

func readUser(id int) (reply []interface{}, err error) {
	c := conn.Pool.Get()
	defer c.Close()
	user := "user:" + strconv.Itoa(id)
	reply, err = redis.Values(c.Do("HGETALL", user))
	return
}

func readMissionRcv(mid int, uids []int) ([]interface{}, error) {
	c := conn.Pool.Get()
	defer c.Close()
	replys := make([]interface{}, len(uids))
	receiverSet := "mission:receivers" + strconv.Itoa(mid)
	//todo lua smembers

	return replys, nil
}

func readMission(mid int) (reply []interface{}, err error) {
	c := conn.Pool.Get()
	defer c.Close()
	mission := "mission:" + strconv.Itoa(mid)
	reply, err = redis.Values(c.Do("HGETALL", mission))
	return
}

func updateMissionRcv(mid, uid int) (err error) {
	c := conn.Pool.Get()
	defer c.Close()
	receiverSet := "mission:receivers" + strconv.Itoa(mid)
	_, err = c.Do("SADD", receiverSet, uid)
	return
}

func createMission(m *Mission) {

}

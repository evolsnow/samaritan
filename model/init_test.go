package model

import (
	"github.com/evolsnow/samaritan/common/caches"
	"github.com/evolsnow/samaritan/common/dbms"
	"time"
)

var cache *caches.SimpleCache

func init() {
	dbms.Pool = dbms.NewPool("127.0.0.1:6379", "", "2")
	dbms.CachePool = dbms.NewPool("127.0.0.1:6379", "", "9")
	c := dbms.Pool.Get()
	defer c.Close()
	c.Do("FLUSHDB")
	cc := dbms.CachePool.Get()
	defer cc.Close()
	cc.Do("FLUSHDB")
	beforeTest()
}

func beforeTest() {
	cache = caches.NewCache()
	u := &User{
		SamId:    "evol",
		Alias:    "evol",
		Name:     "张三",
		Phone:    "13212341234",
		Password: "oldpwd",
		Email:    "gsc1215225@gmail.com",
	}
	u.Save()
	dbms.CreateSearchIndex(u.Id, "gsc1215225@gmail.com", "mail")
	cache.Set("gsc1215225@gmail.com:code", "123456", time.Minute*5)

	t := &Todo{
		OwnerId:   u.Id,
		StartTime: time.Now().Unix(),
		Desc:      "desc here",
	}
	t.Save()
	cache.Set("delete_test_todo_pid", t.Pid, time.Minute*5)
}

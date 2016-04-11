package handler

import (
	"github.com/evolsnow/samaritan/common/caches"
	"github.com/evolsnow/samaritan/common/dbms"
	"github.com/evolsnow/samaritan/model"
	"time"
)

//var cache *caches.SimpleCache

func init() {
	dbms.Pool = dbms.NewPool("127.0.0.1:6379", "", "1")
	dbms.CachePool = dbms.NewPool("127.0.0.1:6379", "", "8")
	c := dbms.Pool.Get()
	c.Do("FLUSHDB")
	c.Close()

	cc := dbms.CachePool.Get()
	cc.Do("FLUSHDB")
	cc.Close()

	beforeTest()
}

func beforeTest() {
	cache = caches.NewCache()
	u := &model.User{
		SamId:    "evol",
		Alias:    "evol",
		Name:     "张三",
		Phone:    "13212341234",
		Password: "oldpwd",
		Email:    "gsc1215225@gmail.com",
	}
	u.Save()
	go u.CreateAvatar()
	dbms.CreateSearchIndex(u.Id, "gsc1215225@gmail.com", "mail")
	cache.Set("gsc1215225@gmail.com:code", "123456", time.Minute*5)

	t1 := &model.Todo{
		OwnerId:   u.Id,
		StartTime: time.Now().Unix(),
		Desc:      "todo 1 desc here",
	}
	t1.Save()
	//dbms.CreateTodoIndex(t.Id, t.Pid)
	cache.Set("delete_test_todo_pid", t1.Pid, time.Minute*5)

	t2 := &model.Todo{
		OwnerId:   u.Id,
		StartTime: time.Now().Unix(),
		Desc:      "todo 2 desc here",
	}
	t2.Save()
	//dbms.CreateTodoIndex(t2.Id, t2.Pid)
	cache.Set("put_test_todo_pid", t2.Pid, time.Minute*5)

	m := &model.Mission{
		Name: "test mission",
	}
	m.Save()
	m.AddReceiver(u.Id)
	m.AddReceiver(10)
	cache.Set("put_test_mission_pid", m.Pid, time.Minute*5)

	p := &model.Project{
		CreatorId: u.Id,
		Name:      "pj name",
		Desc:      "pj desc",
	}
	p.Save()
	cache.Set("delete_test_project_pid", p.Pid, time.Minute*5)
}

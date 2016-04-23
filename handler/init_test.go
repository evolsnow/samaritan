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
	u2 := &model.User{
		SamId:    "snow",
		Alias:    "snow",
		Name:     "王二",
		Phone:    "13212341234",
		Password: "oldpwd",
		Email:    "snow@gmail.com",
	}
	u2.Save()
	dbms.CreateSearchIndex(u2.Id, "snow@gmail.com", "mail")

	cache.Set("get_test_user_pid", u2.Pid, time.Minute*5)
	t1 := &model.Todo{
		OwnerId:   u.Id,
		StartTime: time.Now().Unix(),
		Desc:      "todo 1 desc here",
	}
	t1.Save()
	//dbms.CreateTodoIndex(t.Id, t.Pid)

	t2 := &model.Todo{
		OwnerId:   u.Id,
		StartTime: time.Now().Unix(),
		Desc:      "todo 2 desc here",
	}
	t2.Save()
	//dbms.CreateTodoIndex(t2.Id, t2.Pid)
	cache.Set("put_test_todo_pid", t2.Pid, time.Minute*5)

	t3 := &model.Todo{
		OwnerId:   u.Id,
		StartTime: time.Now().Unix(),
		Desc:      "todo 3 desc here",
	}
	t3.Save()
	cache.Set("delete_test_todo_pid", t3.Pid, time.Minute*5)

	m := &model.Mission{
		Name:        "test mission",
		Desc:        "test mission desc",
		PublisherId: u.Id,
		Deadline:    147258369,
	}
	m.Save()
	m.AddReceiver(u.Id)
	m.AddReceiver(u2.Id)
	cache.Set("put_test_mission_pid", m.Pid, time.Minute*5)
	cache.Set("post_test_mission_pid", m.Pid, time.Minute*5)

	m2 := &model.Mission{
		Name:        "test mission",
		Desc:        "test mission desc",
		PublisherId: u.Id,
		Deadline:    147258369,
	}
	m2.Save()
	cache.Set("delete_test_mission_pid", m2.Pid, time.Minute*5)

	p := &model.Project{
		CreatorId: u.Id,
		Name:      "pj name",
		Desc:      "pj desc",
	}
	p.Save()
	cache.Set("delete_test_project_pid", p.Pid, time.Minute*5)

	p2 := &model.Project{
		CreatorId: u.Id,
		Name:      "pj2 name",
		Desc:      "pj2 desc",
	}
	p2.Save()
	cache.Set("put_test_project_pid", p2.Pid, time.Minute*5)
	cache.Set("post_test_project_pid", p2.Pid, time.Minute*5)
}

package dbms

import (
	"github.com/garyburd/redigo/redis"
	"testing"
)

func init() {
	Pool = NewPool("127.0.0.1:6379", "", "1")
	CachePool = NewPool("127.0.0.1:6379", "", "8")
	c := Pool.Get()
	defer c.Close()
	c.Do("FLUSHDB")
	cc := CachePool.Get()
	defer cc.Close()
	cc.Do("FLUSHDB")
}

func TestRedis(t *testing.T) {
	_, err := Get("foo")
	if err != redis.ErrNil {
		t.Error("redis test failed:", err)
		t.FailNow()
	}
}

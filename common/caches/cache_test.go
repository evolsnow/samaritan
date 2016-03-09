package caches

import (
	"github.com/evolsnow/samaritan/common/dbms"
	"testing"
	"time"
)

func init() {
	dbms.CachePool = dbms.NewPool("127.0.0.1:6379", "", 10)
	Cache = NewCache()
}

func TestCacheSetGet(t *testing.T) {
	Cache.Set("foo", "bar", time.Second*10)
	if Cache.Get("foo") != "bar" {
		t.Error("cache set get error")
	}
}

func TestCacheDelete(t *testing.T) {
	Cache.Set("foo", "bar", time.Second*10)
	Cache.Delete("foo")
	if Cache.Get("foo") == "bar" {
		t.Error("cache delete error")
	}
}

func TestCacheExpiry(t *testing.T) {
	Cache.Set("foo", "bar", time.Second*1/2)
	time.Sleep(time.Second)
	if Cache.Get("foo") == "bar" {
		t.Error("cache expiry error")
	}
}

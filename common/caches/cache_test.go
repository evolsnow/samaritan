package caches

import (
	"github.com/evolsnow/samaritan/common/dbms"
	"testing"
	"time"
)

var cache *SimpleCache

func init() {
	dbms.CachePool = dbms.NewPool("127.0.0.1:6379", "", "1")
	cache = NewCache()
}

func TestCacheSetGet(t *testing.T) {
	cache.Set("foo", "bar", time.Second*10)
	if cache.Get("foo") != "bar" {
		t.Error("cache set get error")
	}
}

func TestCacheGetSet(t *testing.T) {
	cache.Set("foo", "bar", time.Second*10)
	if cache.GetSet("foo", "no_bar") != "bar" {
		t.Error("cache getset get error")
	}
	if cache.Get("foo") != "no_bar" {
		t.Error("cache getset set error")
	}
}

func TestCacheDelete(t *testing.T) {
	cache.Set("foo", "bar", time.Second*10)
	cache.Delete("foo")
	if cache.Get("foo") == "bar" {
		t.Error("cache delete error")
	}
}

func TestCacheExpiry(t *testing.T) {
	cache.Set("foo", "bar", time.Second*1/2)
	time.Sleep(time.Second)
	if cache.Get("foo") == "bar" {
		t.Error("cache expiry error")
	}
}

package caches

import (
	"fmt"
	"testing"
)

func init() {
	LRUCache = NewLRUCache(100)
	lru = LRUCache
}

var lru *LCache

func TestLruAddGet(t *testing.T) {
	lru.Add("foo", "bar")
	ret, _ := lru.Get("foo")
	if ret != "bar" {
		t.Error("lru add get error")
	}
}

func TestLruRemove(t *testing.T) {
	lru.Add("foo", "bar")
	lru.Remove("foo")
	ret, _ := lru.Get("foo")
	if ret != nil {
		t.Error("lru delete error")
	}
}

func TestLruLen(t *testing.T) {
	for i := 0; i < 10; i++ {
		key := fmt.Sprintf("key:%d", i)
		lru.Add(key, i)
	}
	if lru.Len() != 10 {
		t.Error("lru len error")
	}
}

func TestLruRotation(t *testing.T) {
	for i := 0; i < 101; i++ {
		key := fmt.Sprintf("key:%d", i)
		lru.Add(key, i)
	}
	ret, _ := lru.Get("key:0")
	if ret != nil {
		t.Error("lru rotate error")
	}
}

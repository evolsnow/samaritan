package caches

import (
	"github.com/evolsnow/samaritan/common/dbms"
	"time"
)

// Different from LRU cache, this normal cache is saved in redis
// and shouldn't be deleted by the system automatically

//var Cache *SimpleCache

type SimpleCache struct {
	cache map[string]interface{}
}

func NewCache() *SimpleCache {
	return &SimpleCache{
		cache: make(map[string]interface{}),
	}
}

// Add adds a value to the cache.
func (c *SimpleCache) Set(key string, value interface{}, px time.Duration) {
	dbms.CacheSet(key, value, px)
}

// Get looks up a key's value from the cache.
func (c *SimpleCache) Get(key string) (value string) {
	return dbms.CacheGet(key)
}

// Delete deletes a key immediately
func (c *SimpleCache) Delete(key string) {
	dbms.CacheDelete(key)
}

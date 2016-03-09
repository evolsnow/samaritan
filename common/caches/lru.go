/*
Copyright 2013 Google Inc.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
     http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Implements an LRU cache.
package caches

import (
	"container/list"
	"github.com/evolsnow/samaritan/common/dbms"
)

var LRUCache *LCache

// Cache is an LRU cache. It is not safe for concurrent access.
// But it is not necessary for project samaritan.
type LCache struct {
	// MaxEntries is the maximum number of cache entries before
	// an item is evicted. Zero means no limit.
	MaxEntries int

	// OnEvicted optionally specifics a callback function to be
	// executed when an entry is purged from the cache.
	//OnEvicted  func(key string, value interface{})

	ll       *list.List
	lruCache map[string]*list.Element
}

// Project samaritan just need string type key

type entry struct {
	key   string
	value interface{}
}

// New creates a new Cache.
// If maxEntries is zero, the cache has no limit and it's assumed
// that eviction is done by the caller.
func NewLRUCache(maxEntries int) *LCache {
	return &LCache{
		MaxEntries: maxEntries,
		ll:         list.New(),
		lruCache:   make(map[string]*list.Element),
	}
}

// Add adds a value to the cache.
func (c *LCache) Add(key string, value interface{}) {
	if c.lruCache == nil {
		c.lruCache = make(map[string]*list.Element)
		c.ll = list.New()
	}
	if ee, ok := c.lruCache[key]; ok {
		c.ll.MoveToFront(ee)
		ee.Value.(*entry).value = value
		return
	}
	ele := c.ll.PushFront(&entry{key, value})
	c.lruCache[key] = ele
	if c.MaxEntries != 0 && c.ll.Len() > c.MaxEntries {
		c.RemoveOldest()
	}
}

// Get looks up a key's value from the cache.
// If failed, load from redis.
func (c *LCache) GetOrRedis(key string) (value interface{}) {
	//	if c.lruCache == nil {
	//		return
	//	}
	if ele, hit := c.lruCache[key]; hit {
		c.ll.MoveToFront(ele)
		return ele.Value.(*entry).value
	} else {
		// look up this key from redis
		ele, _ := dbms.Get(key)
		if ele != "" {
			go c.Add(key, ele)
			return ele
		}
	}
	return
}

// Get looks up a key's value from the cache.
func (c *LCache) Get(key string) (value interface{}, ok bool) {
	if c.lruCache == nil {
		return
	}
	if ele, hit := c.lruCache[key]; hit {
		c.ll.MoveToFront(ele)
		return ele.Value.(*entry).value, true
	}
	return
}

// Remove removes the provided key from the cache.
func (c *LCache) Remove(key string) {
	if c.lruCache == nil {
		return
	}
	if ele, hit := c.lruCache[key]; hit {
		c.removeElement(ele)
	}
}

// RemoveOldest removes the oldest item from the cache.
func (c *LCache) RemoveOldest() {
	if c.lruCache == nil {
		return
	}
	ele := c.ll.Back()
	if ele != nil {
		c.removeElement(ele)
	}
}

func (c *LCache) removeElement(e *list.Element) {
	c.ll.Remove(e)
	kv := e.Value.(*entry)
	delete(c.lruCache, kv.key)
	//	if c.OnEvicted != nil {
	//		c.OnEvicted(kv.key, kv.value)
	//	}
}

// Len returns the number of items in the cache.
func (c *LCache) Len() int {
	if c.lruCache == nil {
		return 0
	}
	return c.ll.Len()
}

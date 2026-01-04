package cache

import (
	"sync"
	"time"

	kurohelpercore "github.com/kuro-helper/core/v2"
)

// cache struct
type Cache struct {
	Value    any
	ExpireAt time.Time
}

type CacheStore struct {
	data       map[string]*Cache
	expireTime time.Duration
	mu         sync.RWMutex
}

var (
	// 一般查詢快取
	SearchCache = NewCacheStore(10 * time.Minute)
	// 個人資料快取
	UserInfoCache = NewCacheStore(1 * time.Minute)
)

// make new cache store
func NewCacheStore(expireTime time.Duration) *CacheStore {
	return &CacheStore{
		data:       make(map[string]*Cache),
		expireTime: expireTime,
	}
}

// set cache
func (c *CacheStore) Set(key string, value any) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data[key] = &Cache{
		Value:    value,
		ExpireAt: time.Now().Add(c.expireTime),
	}
}

// get cache
func (c *CacheStore) Get(key string) (any, error) {
	c.mu.RLock()
	item, ok := c.data[key]
	c.mu.RUnlock()

	// 不存在或已過期
	if !ok || time.Now().After(item.ExpireAt) {
		c.mu.Lock()
		delete(c.data, key)
		c.mu.Unlock()
		return nil, kurohelpercore.ErrCacheLost
	}

	return item.Value, nil
}

// clean all cache
func (c *CacheStore) Clean() (count int) {
	c.mu.Lock()
	defer c.mu.Unlock()

	for k := range c.data {
		delete(c.data, k)
		count++
	}

	return
}

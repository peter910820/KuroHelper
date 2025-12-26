package cache

import (
	"sync"
	"time"

	kurohelpercore "github.com/peter910820/kurohelper-core"
)

// cache struct
type Cache struct {
	Value    any
	ExpireAt time.Time
}

type CacheStore struct {
	data map[string]*Cache
	mu   sync.RWMutex
}

var (
	UserInfoCache = NewCacheStore()
)

func NewCacheStore() *CacheStore {
	return &CacheStore{
		data: make(map[string]*Cache),
	}
}

func (c *CacheStore) Set(key string, value any) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data[key] = &Cache{
		Value:    value,
		ExpireAt: time.Now().Add(10 * time.Minute),
	}
}

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

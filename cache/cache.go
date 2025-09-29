package cache

import (
	"sync"
	"time"

	kurohelpererrors "kurohelper/errors"
)

// cache struct
type Cache struct {
	Value    any
	ExpireAt time.Time
}

var (
	commonCache   = make(map[string]*Cache)
	commonCacheMu sync.RWMutex
)

func Set(key string, value any) {
	commonCacheMu.Lock()
	defer commonCacheMu.Unlock()
	commonCache[key] = &Cache{
		Value:    value,
		ExpireAt: time.Now().Add(10 * time.Minute),
	}
}

func Get(key string) (any, error) {
	commonCacheMu.RLock()
	item, ok := commonCache[key]
	commonCacheMu.RUnlock()
	// 檢查過期或不存在
	if !ok || time.Now().After(item.ExpireAt) {
		commonCacheMu.Lock()
		delete(commonCache, key)
		commonCacheMu.Unlock()
		return nil, kurohelpererrors.ErrCacheLost
	}
	return item.Value, nil
}

func Clean() {
	commonCacheMu.Lock()
	defer commonCacheMu.Unlock()
	for k := range commonCache {
		delete(commonCache, k)
	}
}
